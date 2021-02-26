package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vpulimamidi/grpc-go-course/compute-service/computepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	targetHost = "localhost:9988"
)

func main() {
	fmt.Printf("Client for Compute service\n")
	tlsEnabled := false
	opts := grpc.WithInsecure()
	if tlsEnabled {
		// Certificate Authority Trust certificate
		certFile := "../../ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	// Create a client connection using target host
	clientConnection, err := grpc.Dial(targetHost, opts)
	if err != nil {
		log.Fatalf("Could not connect to server %v", err)
	}
	// This code will make sure to close the connection at the end of this main function
	defer clientConnection.Close()
	calculatorClient := computepb.NewCalculatorAPIClient(clientConnection)
	req1 := &computepb.DivideRequest{
		Dividend: 20,
		Divisor:  10,
	}
	req2 := &computepb.DivideRequest{
		Dividend: 20,
		Divisor:  0,
	}
	// success case
	doDevide(calculatorClient, req1)
	// failure case
	doDevide(calculatorClient, req2)

	//success
	doSum(calculatorClient, time.Second*5)
	// failure
	doSum(calculatorClient, time.Second*1)
}

func doDevide(calculatorClient computepb.CalculatorAPIClient, req *computepb.DivideRequest) {
	fmt.Println("\nRequest: ", req)
	res, err := calculatorClient.Divide(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println("Response Code: ", respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Invalid divisor is sent!")
				return
			}
		} else {
			log.Fatalf("Error calling Divide method: %v", err)
			return
		}
	}
	fmt.Println("Response: ", res)
}

func doSum(c computepb.CalculatorAPIClient, timeout time.Duration) {
	fmt.Println("\nUnary RPC call with dead line...")
	req := &computepb.SumRequest{
		Number1: 10,
		Number2: 20,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	fmt.Println("Request: ", req)
	fmt.Printf("\nServer should be respond within %f seconds!", timeout.Seconds())
	defer cancel()
	res, err := c.Sum(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Error Code: ", statusErr.Code())
				fmt.Println("Error Message: ", statusErr.Message())
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling Sum RPC: %v", err)
		}
		return
	}
	log.Printf("\nResponse: %v", res.Result)
}
