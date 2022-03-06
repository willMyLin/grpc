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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "go_basic/grpc_test/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) ServerStreamSayHello(in *pb.HelloRequest, stream pb.Greeter_ServerStreamSayHelloServer) error {
	for i := 0; i < 3; i++ {
		err := stream.Send(&pb.HelloReply{Message: "reply" + in.GetName()})

		if err != nil {
			log.Println("stream send err:", err.Error())
		}
	}
	return nil //说明发送完成
}

func (s *server) ClientStreamSayHello(stream pb.Greeter_ClientStreamSayHelloServer) error {
	var i int
	for {
		i++
		req, err := stream.Recv()
		//接受完成
		if err == io.EOF {
			stream.SendAndClose(&pb.HelloReply{
				Message: fmt.Sprintf("message:%s,times:%d", req.GetName(), i),
			})
			break
		}
		if err != nil {
			log.Fatalln("recv err:", err.Error())
		}
	}
	return nil
}

func (s *server) BothStreamSayHello(stream pb.Greeter_BothStreamSayHelloServer) error {

	var msg = make(chan string)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		//发送数据
		defer wg.Done()
		for v := range msg {
			data := fmt.Sprintf("data from server %s", v)
			err := stream.Send(&pb.HelloReply{Message: data})
			if err != nil {
				log.Fatalln("recv err:", err.Error())
				continue
			}
		}

	}()

	wg.Add(1)
	go func() {
		//接收数据
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("recv err:", err.Error())
			}

			msg <- req.GetName()
			fmt.Println("recv name:", req.GetName())
		}
		close(msg)
	}()
	wg.Wait()
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
