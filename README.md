# grpc-go-course

This repository has couple of examples to understand gRPC framework using go language

**Concepts covered in this course are:**
- How to build 
  -   Unary gRPC API
  -   Server Streaming gRPC API
  -   Client Streaming gRPC API
  -   Bi-Directional Streaming gRPC API
  -   Unary gRPC API with Dead Line
  -   Error Handling for gRPC error codes
  -   SSL security for the gRPC API


**Generate the go files from a .proto file use the below command**
  
  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative **<proto file path>**
  
  Run below command to generate go files from book.proto file (from **book-search-service** example)
    
    book-search-service(master)]$  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative bookpb/book.proto
  
  OR Simply run the below script
    
    book-search-service(master)]$ sh generatepb.sh
  
  
**Run the grpc client and server using below commands**

 Run the server
 
    book-search-service/server(master)]$ go run server.go
 
 Run the client
 
    book-search-service/client(master)]$ go run client.go
 
 
