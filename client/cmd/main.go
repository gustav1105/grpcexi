package main

import (
	"log"

	"grpcexi/client/internal/tui"
	pb "grpcexi/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"dns:///localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewContactServiceClient(conn)

	app := tui.NewApp(client)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
