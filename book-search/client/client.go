package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"bookpb.module/bookpb"
	pb "bookpb.module/bookpb"
	"google.golang.org/grpc"
)

const (
	targetHost = "localhost:8989"
)

func main() {
	fmt.Printf("---This is a book search client---\n")
	// Create a client connection using target host
	clientConnection, err := grpc.Dial(targetHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to server %v", err)
	}
	// This code will make to close the connection at the end of this main function
	defer clientConnection.Close()
	// Create book search client using clientConnection
	bookSearchClient := pb.NewBookSearchAPIClient(clientConnection)
	// Example for Unary API call logic

	doUnaryAPICall(bookSearchClient)
	// Example for Server streaming API call logic
	doServerStreamingAPICall(bookSearchClient)
	// Example for Client streaming API call logic
	doClientStreamingAPICall(bookSearchClient)
	// Example for BI-Directional API call logic
	doBiDiStreamingAPICall(bookSearchClient)
}

func doUnaryAPICall(bookSearchClient pb.BookSearchAPIClient) {
	fmt.Println("\n\n----------- Unary API Call Example - Start -----------")
	// Build a BookSearchByTitleRequest
	request := &pb.GetBookRequest{
		Title: "Java",
	}
	bookSearchResponse, err := bookSearchClient.GetBook(context.Background(), request)
	if err != nil {
		log.Printf("Error message: %s", err)
	}
	if bookSearchResponse != nil {
		log.Printf("Request:  %v\n", request)
		printBook(bookSearchResponse.Book)
	}
	fmt.Println("----------- Unary API Call Example - End -----------")
}

func doServerStreamingAPICall(c pb.BookSearchAPIClient) {
	fmt.Println("\n\n----------- Server Streaming RPC Call Example - Start -----------")
	req := &pb.GetAllBooksRequest{
		Title: "Java",
	}
	resStream, err := c.GetAllBooks(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetAllBooks RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		printBook(msg.GetBook())
		time.Sleep(1 * time.Second)
	}
	fmt.Println("----------- Server Streaming RPC Call Example - End -----------")
}

func doClientStreamingAPICall(c bookpb.BookSearchAPIClient) {
	fmt.Println("\n\n----------- Client Streaming RPC Call Example - Start -----------")
	requests := []*pb.GetBooksForGivenTitlesRequest{
		&bookpb.GetBooksForGivenTitlesRequest{
			Title: "Java",
		},
		&bookpb.GetBooksForGivenTitlesRequest{
			Title: "Domain Driven Design",
		},
	}
	stream, err := c.GetBooksForGivenTitles(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GetTheBooksForGivenTitles: %v", err)
	}
	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request..: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from GetBooksForGivenTitles: %v", err)
	}
	data, _ := json.MarshalIndent(res.Book, "", "    ")
	// Print formatted JSON
	log.Println("Response:\n", string(data))
	fmt.Println("----------- Client Streaming RPC Call Example - End -----------")
}

func doBiDiStreamingAPICall(c pb.BookSearchAPIClient) {
	fmt.Println("\n\n\n----------- Bi-Directional Streaming RPC Call Example - Start -----------")
	// Create a stream
	stream, err := c.GetEachBook(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	requests := []*pb.GetEachBookRequest{
		&pb.GetEachBookRequest{
			Title: "Domain Driven Design",
		},
		&pb.GetEachBookRequest{
			Title: "Java",
		},
	}
	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("\nSending Request..: %v\n", req)
			stream.Send(req)
			time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Println("Received Response:")
			printBook(res.GetBook())
			time.Sleep(2 * time.Second)
		}
		close(waitc)
	}()
	// block until everything is done
	<-waitc
	fmt.Println("----------- Bi-Directional Streaming RPC Call Example - End -----------")
}

func printBook(book *pb.Book) {
	data, _ := json.MarshalIndent(book, "", "    ")
	// Print formatted JSON
	log.Println("Response:\n", string(data))
}
