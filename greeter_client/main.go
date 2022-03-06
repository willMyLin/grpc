/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	pb "go_basic/grpc_test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//SayHello(c, ctx)

	//serverStreamSayHello(c, ctx)
	//clientStreamSayHello(c, ctx)
	bothStreamSayHello(c, ctx)
}

func bothStreamSayHello(cli pb.GreeterClient, ctx context.Context) {

	stream, err := cli.BothStreamSayHello(ctx)
	if err != nil {
		log.Println("err:", err.Error())
	}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			rep, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("recv err:", err.Error())
				continue
			}
			fmt.Println("data from server:", rep.GetMessage())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			err = stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("client data :%s%d", "will", i)})
			if err != nil {
				log.Println("err:", err.Error())
			}
			time.Sleep(time.Second)
		}

		err = stream.CloseSend()
		if err != nil {
			log.Println("CloseSend err:", err.Error())
			return
		}
	}()

	wg.Wait()
}

func clientStreamSayHello(cli pb.GreeterClient, ctx context.Context) {

	stream, err := cli.ClientStreamSayHello(ctx)

	if err != nil {
		log.Fatalln("err:", err.Error())
	}

	for i := 0; i < 3; i++ {
		err := stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("name:%s-%d", "will", i)})
		if err != nil {
			log.Fatalln("err:", err.Error())
		}
	}

	rep, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("err:", err.Error())
	}
	fmt.Println("reply:", rep.GetMessage())
}

func serverStreamSayHello(cli pb.GreeterClient, ctx context.Context) {
	stream, err := cli.ServerStreamSayHello(ctx, &pb.HelloRequest{Name: "will"})

	if err != nil {
		log.Fatalln("err:", err.Error())
	}

	for {
		rep, err := stream.Recv()

		//全部发送完成
		if err == io.EOF {
			break
		}
		//发送错误
		if err != nil {
			log.Fatalln("recv err:", err.Error())
		}

		fmt.Println("recv：", rep.Message)
	}

}

func SayHello(cli pb.GreeterClient, ctx context.Context) {
	r, err := cli.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
