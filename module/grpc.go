package module

import "github.com/vseinstrumentiru/lego/v2/transport/grpc"

func GRPCServer() (interface{}, []interface{}) {
	return grpc.Provide, nil
}
