package main

import (
	"context"
	"flag"
	"log"
	"time"

	// "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/huydevct/todo-grpc/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

func main() {
	address := flag.String("server", "", ":9090")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	insert_at := timestamppb.New(t)
	pfx := t.Format(time.RFC3339Nano)

	// call create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			InsertAt:    insert_at,
			UpdateAt:    insert_at,
		},
	}

	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// call read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)

	// call update
	req3 := v1.UpdateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id:          res2.ToDo.Id,
			Title:       res2.ToDo.Title,
			Description: res2.ToDo.Description + " + updated",
			UpdateAt:    res2.ToDo.UpdateAt,
		},
	}

	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("Update failed: %v",err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)

	// call readAll
	req5 := v1.DeleteRequest{
		Api: apiVersion,
		Id: id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)


}
