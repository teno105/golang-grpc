package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "golang-grpc/pkg/server/generated"

	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func newServer() *chatServer {
	return &chatServer{}
}

type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	streams []pb.ChatService_ChatServer
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	s.mu.Lock()
	s.streams = append(s.streams, stream)
	s.mu.Unlock()

	var err error
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		s.mu.Lock()
		for _, strm := range s.streams {
			strm.Send(&pb.ChatMsg{
				Sender:  in.Sender,
				Message: in.Message,
			})
		}
		s.mu.Unlock()
	}

	s.mu.Lock()
	for i, strm := range s.streams {
		if strm == stream {
			s.streams = append(s.streams[:i], s.streams[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
	return err
}
