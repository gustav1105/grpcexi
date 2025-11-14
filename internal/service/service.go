package service

import "google.golang.org/grpc"

type GrpcService interface {
	Register(grpcServer *grpc.Server)
}
