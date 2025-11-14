package main

import (
	"google.golang.org/grpc"
	"grpcexi/internal/repo"
	"grpcexi/internal/service"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	contactRepo := repo.NewSQLiteContactRepository("./contacts.db")
	services := []service.GrpcService{
		service.NewContactService(contactRepo),
	}

	for _, s := range services {
		s.Register(grpcServer)
	}

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
