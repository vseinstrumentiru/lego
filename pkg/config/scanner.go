package config

import (
	"reflect"

	"emperror.dev/errors"
)

// Scanner gathers information from struct fields
type Scanner interface {
	Scan(f *Field)
}

// ScanFunc is func implementation of Scanner
type ScanFunc func(f *Field)

func (fn ScanFunc) Scan(f *Field) {
	fn(f)
}

// Scanners is slice implementation of Scanner
type Scanners []Scanner

func (s Scanners) Scan(f *Field) {
	for _, scanner := range s {
		scanner.Scan(f)
	}
}

// ScannerGroup group scanners in one
func ScannerGroup(scanners ...Scanner) Scanner {
	return Scanners(scanners)
}

func getDefaultScanner() Scanner {
	return ScannerGroup(
		ScanFunc(typeScanner),
		ScanFunc(mapStructureTagScanner),
		ScanFunc(defaultsScanner),
		ScanFunc(flagTagScanner),
	)
}

func scan(path string, v reflect.Value, scanner Scanner) ([]*Field, error) {
	t := v.Elem().Type()
	var fields []*Field
	for i := 0; i < t.NumField(); i++ {
		field, err := newField(path, t.Field(i), v.Elem().Field(i))

		if errors.Is(err, ErrNotValidValue) {
			continue
		} else if err != nil {
			return nil, err
		}

		scanner.Scan(field)

		if field.isIgnored {
			continue
		}

		field.settings = append(field.settings, bindEnv)
		fields = append(fields, field)

		if !field.hasFields {
			continue
		}

		key := field.key()
		if field.isSquashed {
			key = path
		}

		subItems, err := scan(key, field.value, scanner)
		if err != nil {
			return nil, err
		}

		fields = append(fields, subItems...)
	}

	return fields, nil
}
