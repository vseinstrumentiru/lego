package cast

import (
	"net"
	"time"

	"emperror.dev/errors"
)

func NoDefaults(castable Castable) Castable {
	if _, ok := castable.(*noDefaults); ok {
		return castable
	}

	return &noDefaults{castable}
}

func ErrValueNotFound() error {
	return errors.New("value not found")
}

type noDefaults struct {
	i Castable
}

func (n *noDefaults) GetString(key string) (string, error) {
	val, err := n.i.GetString(key)
	if err != nil {
		return val, err
	}

	if val == "" {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetStringSlice(key string) ([]string, error) {
	val, err := n.i.GetStringSlice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetBool(key string) (bool, error) {
	val, err := n.i.GetBool(key)
	if err != nil {
		return val, err
	}

	if !val {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetBoolSlice(key string) ([]bool, error) {
	val, err := n.i.GetBoolSlice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetBytesHex(key string) ([]byte, error) {
	val, err := n.i.GetBytesHex(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetDuration(key string) (time.Duration, error) {
	val, err := n.i.GetDuration(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetDurationSlice(key string) ([]time.Duration, error) {
	val, err := n.i.GetDurationSlice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetFloat32(key string) (float32, error) {
	val, err := n.i.GetFloat32(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetFloat32Slice(key string) ([]float32, error) {
	val, err := n.i.GetFloat32Slice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetFloat64(key string) (float64, error) {
	val, err := n.i.GetFloat64(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetFloat64Slice(key string) ([]float64, error) {
	val, err := n.i.GetFloat64Slice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt(key string) (int, error) {
	val, err := n.i.GetInt(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetIntSlice(key string) ([]int, error) {
	val, err := n.i.GetIntSlice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt8(key string) (int8, error) {
	val, err := n.i.GetInt8(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt16(key string) (int16, error) {
	val, err := n.i.GetInt16(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt32(key string) (int32, error) {
	val, err := n.i.GetInt32(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt32Slice(key string) ([]int32, error) {
	val, err := n.i.GetInt32Slice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt64(key string) (int64, error) {
	val, err := n.i.GetInt64(key)
	if err != nil {
		return val, err
	}

	if val == 0 {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetInt64Slice(key string) ([]int64, error) {
	val, err := n.i.GetInt64Slice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetIP(key string) (net.IP, error) {
	val, err := n.i.GetIP(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetIPSlice(key string) ([]net.IP, error) {
	val, err := n.i.GetIPSlice(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetIPv4Mask(key string) (net.IPMask, error) {
	val, err := n.i.GetIPv4Mask(key)
	if err != nil {
		return val, err
	}

	if val == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}

func (n *noDefaults) GetIPNet(key string) (net.IPNet, error) {
	val, err := n.i.GetIPNet(key)
	if err != nil {
		return val, err
	}

	if val.IP == nil {
		return val, ErrValueNotFound()
	}

	return val, nil
}
