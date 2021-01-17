package provide

import "github.com/vseinstrumentiru/lego/v2/transport/grpc"

func Grpc() []interface{} {
	return []interface{}{
		grpc.Provide,
	}
}
