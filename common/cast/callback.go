package cast

import (
	"net"
	"reflect"
	"time"

	"emperror.dev/errors"
)

func On(c Castable, key string, callback interface{}) (err error) {
	val := reflect.ValueOf(callback)

	c = NoDefaults(c)

	if val.Kind() != reflect.Func {
		return errors.New("callback must be a function")
	}

	if val.Type().NumIn() != 1 {
		return errors.New("callback must accept only one argument")
	}

	if val.Type().NumIn() != 1 {
		return errors.New("callback must accept only one argument")
	}

	if val.Type().In(0).Kind() == reflect.Ptr {
		return errors.New("callback argument must be non-pointer value")
	}

	t := reflect.New(val.Type().In(0)).Elem().Interface()
	var arg interface{}

	switch t.(type) {
	case string:
		arg, err = c.GetString(key)
	case []string:
		arg, err = c.GetStringSlice(key)
	case bool:
		arg, err = c.GetBool(key)
	case []bool:
		arg, err = c.GetBoolSlice(key)
	case []byte:
		arg, err = c.GetBytesHex(key)
	case time.Duration:
		arg, err = c.GetDuration(key)
	case []time.Duration:
		arg, err = c.GetDurationSlice(key)
	case float32:
		arg, err = c.GetFloat32(key)
	case []float32:
		arg, err = c.GetFloat32Slice(key)
	case float64:
		arg, err = c.GetFloat64(key)
	case []float64:
		arg, err = c.GetFloat64Slice(key)
	case int:
		arg, err = c.GetInt(key)
	case int8:
		arg, err = c.GetInt8(key)
	case int16:
		arg, err = c.GetInt16(key)
	case int32:
		arg, err = c.GetInt32(key)
	case int64:
		arg, err = c.GetInt64(key)
	case []int:
		arg, err = c.GetIntSlice(key)
	case []int32:
		arg, err = c.GetInt32Slice(key)
	case []int64:
		arg, err = c.GetInt64Slice(key)
	case net.IP:
		arg, err = c.GetIP(key)
	case []net.IP:
		arg, err = c.GetIPSlice(key)
	case net.IPMask:
		arg, err = c.GetIPv4Mask(key)
	case net.IPNet:
		arg, err = c.GetIPNet(key)
	default:
		return errors.NewWithDetails("type not found", "key", key)
	}

	if err != nil {
		return err
	}

	val.Call([]reflect.Value{reflect.ValueOf(arg)})

	return nil
}

func OnCheck(c Castable, key string, callback interface{}) (ok bool) {
	return On(c, key, callback) == nil
}

func NewCallback(set Castable) Callback {
	return &callback{set}
}

type Callback interface {
	On(key string, callback interface{}) (err error)
}

func NewCallbackCheck(set Castable) CallbackCheck {
	return &callbackCheck{set}
}

type CallbackCheck interface {
	On(key string, callback interface{}) (ok bool)
}

type callback struct {
	Castable
}

func (c *callback) On(key string, cb interface{}) error {
	return On(c, key, cb)
}

type callbackCheck struct {
	Castable
}

func (c *callbackCheck) On(key string, cb interface{}) bool {
	return OnCheck(c, key, cb)
}
