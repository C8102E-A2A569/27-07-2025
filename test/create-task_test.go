package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	proto "zip-archive/protos"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect: ", err)
	}
	defer conn.Close()

	client := proto.NewArchiveZipServiceClient(conn)

	resp, err := client.CreateTask(context.Background(), &proto.CreateTaskRequest{})
	if err != nil {
		log.Fatal("createTask failed: ", err)
	}

	log.Printf("task created: id=%s, status=%s", resp.TaskId, resp.Status)
	
}
