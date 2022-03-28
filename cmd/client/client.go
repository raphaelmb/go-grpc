package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/raphaelmb/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser - unary
	// AddUser(client)

	// AddUserVerbose - server stream
	// AddUserVerbose(client)

	// AddUsers - client stream
	// AddUsers(client)

	// AddUserStreamBoth - both client and server stream
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{Id: "0", Name: "Joao", Email: "j@j.com"}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{Id: "0", Name: "Joao", Email: "j@j.com"}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status:", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "r1",
			Name:  "Raphael 1",
			Email: "rpl@rpl.com",
		},
		{
			Id:    "r2",
			Name:  "Raphael 2",
			Email: "rpl2@rpl.com",
		},
		{
			Id:    "r3",
			Name:  "Raphael 3",
			Email: "rpl3@rpl.com",
		},
		{
			Id:    "r4",
			Name:  "Raphael 4",
			Email: "rpl4@rpl.com",
		},
		{
			Id:    "r5",
			Name:  "Raphael 5",
			Email: "rpl5@rpl.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "r1",
			Name:  "Raphael 1",
			Email: "rpl@rpl.com",
		},
		{
			Id:    "r2",
			Name:  "Raphael 2",
			Email: "rpl2@rpl.com",
		},
		{
			Id:    "r3",
			Name:  "Raphael 3",
			Email: "rpl3@rpl.com",
		},
		{
			Id:    "r4",
			Name:  "Raphael 4",
			Email: "rpl4@rpl.com",
		},
		{
			Id:    "r5",
			Name:  "Raphael 5",
			Email: "rpl5@rpl.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
