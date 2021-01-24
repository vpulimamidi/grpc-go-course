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
	"errors"
	"fmt"
	"log"
	"net"

	pb "bookpb.module/bookpb"

	"google.golang.org/grpc"
)

const (
	port = ":8989"
)

type server struct {
	pb.UnimplementedBookSearchServer
}

func (s *server) BookSearchByTitle(ctx context.Context, req *pb.BookSearchByTitleRequest) (*pb.BookSearchResponse, error) {
	fmt.Printf("\nRequest information : %v", req)
	book := getBookByTitle(req.GetTitle())
	if book != nil {
		response := &pb.Book{
			Title:    book.title,
			Subject:  book.subject,
			Audience: book.audience,
			Author:   book.author,
			Price:    book.price,
		}
		return &pb.BookSearchResponse{
			Book: response,
		}, nil
	}
	return nil, errors.New("Book is not found for a given Title: " + req.GetTitle())
}

func (s *server) BookSearchByAuthor(ctx context.Context, req *pb.BookSearchByAuthorRequest) (*pb.BookSearchResponse, error) {
	fmt.Printf("\nRequest information : %v", req)
	book := getBookByAuthor(req.GetAuthor())
	if book != nil {
		response := &pb.Book{
			Title:    book.title,
			Subject:  book.subject,
			Audience: book.audience,
			Author:   book.author,
			Price:    book.price,
		}
		return &pb.BookSearchResponse{
			Book: response,
		}, nil
	}
	return nil, errors.New("Book is not found for a given Author: " + req.GetAuthor())
}

func main() {
	log.Printf("Book search server is running.......")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBookSearchServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type book struct {
	title    string
	subject  string
	audience string
	author   string
	price    float32
}

// Return the list of books
func getBooks() []book {
	books := make([]book, 2)
	books[0] = book{
		title:    "Domain Driven Design",
		subject:  "Tackling complexity in the heart of Software",
		audience: "Software Engineers",
		author:   "Eric Evans",
		price:    999.0,
	}
	books[1] = book{
		title:    "Java",
		subject:  "Comprehensive guide to the entire Java laguage",
		audience: "Software Engineers",
		author:   "Herbert Schildt",
		price:    999.0,
	}
	return books
}

// Find the book by title and return it
func getBookByTitle(title string) *book {
	books := getBooks()
	for _, book := range books {
		if book.title == title {
			return &book
		}
	}
	return nil
}

// Find the book by Author
func getBookByAuthor(author string) *book {
	books := getBooks()
	for _, book := range books {
		if book.author == author {
			return &book
		}
	}
	return nil
}
