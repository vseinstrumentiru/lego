package cast

import (
	"net"
	"time"
)

type Castable interface {
	GetString(key string) (string, error)
	GetStringSlice(key string) ([]string, error)
	GetBool(key string) (bool, error)
	GetBoolSlice(key string) ([]bool, error)
	GetBytesHex(key string) ([]byte, error)
	GetDuration(key string) (time.Duration, error)
	GetDurationSlice(key string) ([]time.Duration, error)
	GetFloat32(key string) (float32, error)
	GetFloat32Slice(key string) ([]float32, error)
	GetFloat64(key string) (float64, error)
	GetFloat64Slice(key string) ([]float64, error)
	GetInt(key string) (int, error)
	GetIntSlice(key string) ([]int, error)
	GetInt8(key string) (int8, error)
	GetInt16(key string) (int16, error)
	GetInt32(key string) (int32, error)
	GetInt32Slice(key string) ([]int32, error)
	GetInt64(key string) (int64, error)
	GetInt64Slice(key string) ([]int64, error)
	GetIP(key string) (net.IP, error)
	GetIPSlice(key string) ([]net.IP, error)
	GetIPv4Mask(key string) (net.IPMask, error)
	GetIPNet(key string) (net.IPNet, error)
}
