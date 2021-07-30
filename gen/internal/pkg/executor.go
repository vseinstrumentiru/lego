package pkg

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Executor interface {
	Exec(fn func() error, comment string) error
	CopyStub(from, to string, args ...interface{}) error
	Mkdir(path string) (exist bool, err error)
	Exist(path string) (bool, error)
	Result() []string
	log(msg string)
}

func NewExecutor(workDir string, verbose bool, dryRun bool) Executor {
	if dryRun {
		return &dryRunner{}
	}

	return &runner{path: workDir, verbose: verbose}
}

func NewStubExecutor(stubs embed.FS, workDir string, verbose bool, dryRun bool) Executor {
	if dryRun {
		return &dryRunner{stubs: stubs}
	}

	return &runner{path: workDir, verbose: verbose, stubs: stubs}
}

type runner struct {
	path    string
	stubs   embed.FS
	res     []string
	verbose bool
}

func (r *runner) log(msg string) {
	r.res = append(r.res, msg)
	if r.verbose {
		//nolint:forbidigo
		fmt.Println(msg)
	}
}

func (r *runner) Exist(path string) (bool, error) {
	fullPath, err := filepath.Abs(r.path + "/" + path)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		r.log("file " + path + " already exist")

		return true, nil
	}

	return false, nil
}

func (r *runner) Exec(fn func() error, comment string) error {
	if err := fn(); err != nil {
		return err
	}

	r.log(comment)

	return nil
}

func (r *runner) Mkdir(path string) (exist bool, err error) {
	fullPath, err := filepath.Abs(r.path + "/" + path)
	if err != nil {
		return false, err
	}

	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		err = r.Exec(func() error {
			return os.MkdirAll(fullPath, os.ModePerm)
		}, fmt.Sprintf("creating directory %s", path))

		if err != nil {
			return false, err
		}
	} else {
		return true, nil
	}

	return false, nil
}

func (r *runner) CopyStub(from string, to string, args ...interface{}) error {
	fullPath, err := filepath.Abs(r.path + "/" + to)
	if err != nil {
		return err
	}

	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		r.log("file " + to + " already exist")

		return nil
	}

	data, err := r.stubs.ReadFile(from)
	if err != nil {
		return err
	}

	if strings.HasSuffix(from, ".tmpl") {
		tpl, err := template.New(to).Parse(string(data))
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		var vars interface{}
		if len(args) == 1 {
			vars = args[0]
		}
		if err = tpl.Execute(buf, vars); err != nil {
			return err
		}

		data = buf.Bytes()
	}

	return r.Exec(func() error {
		return os.WriteFile(fullPath, data, os.ModePerm)
	}, "creating "+to)
}

func (r *runner) Result() []string {
	return r.res
}

type dryRunner struct {
	res   []string
	stubs embed.FS
}

func (r *dryRunner) log(msg string) {
	r.res = append(r.res, msg)
}

func (r *dryRunner) Exec(_ func() error, comment string) error {
	r.log(comment)

	return nil
}

func (r *dryRunner) CopyStub(from, to string, args ...interface{}) error {
	data, err := r.stubs.ReadFile(from)
	if err != nil {
		return err
	}

	if strings.HasSuffix(from, ".tmpl") {
		tpl, err := template.New(to).Parse(string(data))
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		var vars interface{}
		if len(args) == 1 {
			vars = args[0]
		}
		if err = tpl.Execute(buf, vars); err != nil {
			return err
		}
	}

	r.log("creating " + to)

	return nil
}

func (r *dryRunner) Mkdir(path string) (exist bool, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return false, err
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		_ = r.Exec(func() error {
			return nil
		}, fmt.Sprintf("creating directory %s", path))
	} else {
		return true, nil
	}

	return false, nil
}

func (r *dryRunner) Exist(path string) (bool, error) {
	return false, nil
}

func (r *dryRunner) Result() []string {
	return r.res
}
