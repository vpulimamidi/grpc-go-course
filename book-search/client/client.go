package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	// Create book search client using clientConnection
	bookSearchClient := pb.NewBookSearchClient(clientConnection)
	// search book using title
	bookSearchByTitle(bookSearchClient)
	// search book using author
	bookSearchByAuthor(bookSearchClient)

}

func bookSearchByTitle(bookSearchClient pb.BookSearchClient) {
	log.Println("Book search Request using Title")
	// Build a BookSearchByTitleRequest
	bookSearchByTitleReq := &pb.BookSearchByTitleRequest{
		Title: "Java",
	}
	bookSearchResponse, err := bookSearchClient.BookSearchByTitle(context.Background(), bookSearchByTitleReq)
	if err != nil {
		log.Printf("Error message: %s", err)
	}
	if bookSearchResponse != nil {
		log.Printf("\nRequest:  %v\n", bookSearchByTitleReq)
		printResponse(bookSearchResponse)
	}
}

func bookSearchByAuthor(bookSearchClient pb.BookSearchClient) {
	log.Println("Book search Request using Author")
	// Build a BookSearchByAuthorRequest
	bookSearchByAuthorReq := &pb.BookSearchByAuthorRequest{
		Author: "Eric Evans",
	}
	bookSearchResponse, err := bookSearchClient.BookSearchByAuthor(context.Background(), bookSearchByAuthorReq)
	if err != nil {
		log.Printf("Error message: %s", err)
	}
	if bookSearchResponse != nil {
		log.Printf("\nRequest:  %v\n", bookSearchByAuthorReq)
		printResponse(bookSearchResponse)
	}
}

func printResponse(bookSearchResponse *bookpb.BookSearchResponse) {
	if bookSearchResponse != nil {
		data, _ := json.MarshalIndent(bookSearchResponse.Book, "", "    ")
		// Print formatted JSON
		log.Println("\nResponse:\n", string(data))
	}
}
