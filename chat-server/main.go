package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	pb "gRPC/chat"

	"google.golang.org/grpc"
)

type ChatServer struct {
	pb.UnimplementedChatServer
}

var (
	reader *bufio.Reader
)

// 채팅을 위한 함수 구현
func (cs *ChatServer) StartStreaming(stream pb.Chat_StartStreamingServer) error {
	for {
		// 메세시 수신
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("receive Message From Client is %s\n", message.GetContent())

		// 서버의 응답 생성

		input, _ := reader.ReadString('\n')
		response := &pb.ChatMessage{
			Content: "Server Response Message is " + input,
		}
		if err = stream.Send(response); err != nil {
			return err
		}

		time.Sleep(1000 * time.Millisecond)

	}
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	// gRPC 리스너 생성
	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	// gRPC 서버 생성
	gRPCServer := grpc.NewServer()

	// chatService 등록
	pb.RegisterChatServer(gRPCServer, &ChatServer{})

	// 서버 시작
	fmt.Printf("chat Server Start From %s \n", ":50051")
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("gRPCServer err to serve : %v", err)
	}

}
