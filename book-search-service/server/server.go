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
	"io"
	"log"
	"net"
	"time"

	"github.com/vpulimamidi/grpc-go-course/book-search-service/bookpb"
	"google.golang.org/grpc"
)

const (
	port = ":8989"
)

type server struct {
	bookpb.UnimplementedBookSearchAPIServer
}

func (s *server) GetBook(ctx context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	fmt.Printf("\nRequest information : %v", req)
	book := getBookByTitle(req.GetTitle())
	if book != nil {
		response := &bookpb.Book{
			Title:    book.title,
			Subject:  book.subject,
			Audience: book.audience,
			Author:   book.author,
			Price:    book.price,
		}
		return &bookpb.GetBookResponse{
			Book: response,
		}, nil
	}
	return nil, errors.New("Book is not found for a given Title: " + req.GetTitle())
}

func (s *server) GetAllBooks(req *bookpb.GetAllBooksRequest, stream bookpb.BookSearchAPI_GetAllBooksServer) error {
	fmt.Printf("\n Request: %v", req)
	books := getAllTheBookByTitle(req.GetTitle())
	if books != nil {
		for i := 0; i < len(books); i++ {
			book := books[i]
			result := &bookpb.Book{
				Title:    book.title,
				Subject:  book.subject,
				Audience: book.audience,
				Author:   book.author,
				Price:    book.price,
			}
			stream.Send(&bookpb.GetAllBooksResponse{
				Book: result,
			})
			// Sleep for a second
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

func (s *server) GetBooksForGivenTitles(stream bookpb.BookSearchAPI_GetBooksForGivenTitlesServer) error {
	fmt.Println("GetBooksForGivenTitles")
	index := 0
	books := make([]*bookpb.Book, 100)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			finalBooks := make([]*bookpb.Book, index)
			copy(finalBooks, books)
			// we have finished reading the client stream
			return stream.SendAndClose(&bookpb.GetBooksForGivenTitlesResponse{
				Book: finalBooks,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		book := getBookByTitle(req.GetTitle())
		bk := &bookpb.Book{
			Title:    book.title,
			Subject:  book.subject,
			Audience: book.audience,
			Author:   book.author,
			Price:    book.price,
		}
		books[index] = bk
		index++
	}
}

func (s *server) GetEachBook(stream bookpb.BookSearchAPI_GetEachBookServer) error {
	fmt.Println("Invoking Bi Directional GetBooks method with a streaming request")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		book := getBookByTitle(req.GetTitle())
		if book != nil {
			response := &bookpb.Book{
				Title:    book.title,
				Subject:  book.subject,
				Audience: book.audience,
				Author:   book.author,
				Price:    book.price,
			}
			sendError := stream.Send(&bookpb.GetEachBookResponse{
				Book: response,
			})
			if sendError != nil {
				log.Fatalf("Error while sending data to the client: %v", sendError)
				return sendError
			}
		}
	}
}

func main() {
	log.Printf("Book search server is running.......")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bookpb.RegisterBookSearchAPIServer(s, &server{})
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
	books := make([]book, 4)
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
	books[2] = book{
		title:    "Java",
		subject:  "Head First Java",
		audience: "Software Engineers",
		author:   "Kathy Sierra",
		price:    999.0,
	}
	books[3] = book{
		title:    "Java",
		subject:  "Effective Java",
		audience: "Software Engineers",
		author:   "Joshua Bloch",
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

// Find the book by title and return it
func getAllTheBookByTitle(title string) []book {
	books := getBooks()
	tempBooks := make([]book, len(books))
	index := 0
	for _, book := range books {
		if book.title == title {
			tempBooks[index] = book
			index++
		}
	}
	matchedBooks := make([]book, index)
	copy(matchedBooks, tempBooks)
	return matchedBooks
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
