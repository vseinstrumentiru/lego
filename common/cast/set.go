package cast

import (
	"net"
	"time"

	"emperror.dev/errors"

	"github.com/vseinstrumentiru/lego/v2/common/set"
)

type CastableRWSet interface {
	Castable
	set.Set
}

type CastableSet interface {
	Castable
	set.ReadableSet
}

func NewCastableRWSet(s set.Set) CastableRWSet {
	return &castableRWSet{NewCastableSet(s), s}
}

func NewCastableSet(s set.ReadableSet) CastableSet {
	return &castableSet{s}
}

type castableRWSet struct {
	CastableSet
	set.CheckWritableSet
}

type castableSet struct {
	set.ReadableSet
}

func ErrWrongType() error {
	return errors.New("wrong value type")
}

func (w *castableSet) GetString(key string) (string, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return "", nil
	}

	val, ok := raw.(string)
	if !ok {
		return "", ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetStringSlice(key string) ([]string, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]string)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetBool(key string) (bool, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return false, nil
	}

	val, ok := raw.(bool)
	if !ok {
		return false, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetBoolSlice(key string) ([]bool, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]bool)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetBytesHex(key string) ([]byte, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]byte)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetDuration(key string) (time.Duration, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(time.Duration)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetDurationSlice(key string) ([]time.Duration, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]time.Duration)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetFloat32(key string) (float32, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(float32)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetFloat32Slice(key string) ([]float32, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]float32)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetFloat64(key string) (float64, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(float64)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetFloat64Slice(key string) ([]float64, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]float64)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt(key string) (int, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(int)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetIntSlice(key string) ([]int, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]int)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt8(key string) (int8, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(int8)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt16(key string) (int16, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(int16)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt32(key string) (int32, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(int32)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt32Slice(key string) ([]int32, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]int32)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt64(key string) (int64, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return 0, nil
	}

	val, ok := raw.(int64)
	if !ok {
		return 0, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetInt64Slice(key string) ([]int64, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]int64)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetIP(key string) (net.IP, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.(net.IP)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetIPSlice(key string) ([]net.IP, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.([]net.IP)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetIPv4Mask(key string) (net.IPMask, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return nil, nil
	}

	val, ok := raw.(net.IPMask)
	if !ok {
		return nil, ErrWrongType()
	}

	return val, nil
}

func (w *castableSet) GetIPNet(key string) (net.IPNet, error) {
	var raw interface{}
	if raw = w.Get(key); raw == nil {
		return net.IPNet{}, nil
	}

	val, ok := raw.(net.IPNet)
	if !ok {
		return net.IPNet{}, ErrWrongType()
	}

	return val, nil
}
