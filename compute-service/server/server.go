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

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vpulimamidi/grpc-go-course/compute-service/computepb"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const (
	port = ":9988"
)

type server struct {
	computepb.UnimplementedCalculatorAPIServer
}

//Divide ..
func (s *server) Divide(ctx context.Context, req *computepb.DivideRequest) (*computepb.DivideResponse, error) {
	fmt.Printf("\nRequest : %v", req)
	if req.GetDivisor() == 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received an invalid number: %v", req.GetDivisor()))
	}
	result := req.GetDividend() / req.GetDividend()
	return &computepb.DivideResponse{
		Result: float64(result),
	}, nil
}

func (s *server) Sum(ctx context.Context, req *computepb.SumRequest) (*computepb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	result := req.GetNumber1() + req.GetNumber2()
	res := &computepb.SumResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	log.Printf("Compute service is running.......")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	tlsEnabled := false
	if tlsEnabled {
		certFile := "../../ssl/server.crt"
		keyFile := "../../ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed loading certificates: %v", err)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}
	s := grpc.NewServer(opts...)
	computepb.RegisterCalculatorAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
