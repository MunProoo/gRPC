package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "gRPC/chat"

	"google.golang.org/grpc"
)

var (
	reader *bufio.Reader
)

func main() {
	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect %v", err)
	}
	defer conn.Close()

	// gRPC 클라이언트 생성
	client := pb.NewChatClient(conn)
	// pb.RegisterChatServer(conn, &ChatServer{})

	// 스트리밍을 위한 스트림 생성
	stream, err := client.StartStreaming(context.Background())
	if err != nil {
		log.Fatalf("Failed to start streaming : %v", err)
	}

	// 서버로 메시지 전송
	go func() {
		for {

			reader = bufio.NewReader(os.Stdin)
			content, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Failed to message read : %v", err)
			}

			request := pb.ChatMessage{
				Content: content,
			}

			// 메시지 전송
			if err = stream.Send(&request); err != nil {
				log.Fatalf("Failed to send message : %v", err)
			}

			// messages := []string{"Hello", "How are you?", "I'm fine, thank you!"}

			// for _, msg := range messages {
			// 	request := &pb.ChatMessage{
			// 		Content: msg,
			// 		// 다른 필드를 추가하려면 여기에 추가하세요.
			// 	}

			// 	if err := stream.Send(request); err != nil {
			// 		log.Fatalf("Error sending message: %v", err)
			// 	}
			// 	time.Sleep(1 * time.Second)
			// }

			time.Sleep(100 * time.Microsecond)
		}
	}()

	// 서버로부터 메시지 스트림 응답 수신
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// 서버가 스트림을 닫으면 종료
			break
		}
		if err != nil {
			log.Fatalf("Error receiving response: %v", err)
		}

		fmt.Printf("Received response from server: %s\n", response.GetContent())
	}
}
