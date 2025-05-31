package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Запрашиваем ID трека
	fmt.Print("Введите ID трека: ")
	trackID, _ := reader.ReadString('\n')
	trackID = strings.TrimSpace(trackID)

	// Запрашиваем путь к файлу с токеном
	fmt.Print("Введите путь к файлу с токеном: ")
	tokenPath, _ := reader.ReadString('\n')
	tokenPath = strings.TrimSpace(tokenPath)

	// Запрашиваем адрес сервера
	fmt.Print("Введите адрес сервера (например, localhost:8080): ")
	serverAddr, _ := reader.ReadString('\n')
	serverAddr = strings.TrimSpace(serverAddr)

	// Запрашиваем имя выходного файла
	fmt.Print("Введите имя выходного файла (например, output.mp3): ")
	outputFile, _ := reader.ReadString('\n')
	outputFile = strings.TrimSpace(outputFile)

	// Загрузка токена из файла
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("failed to read token: %v", err)
	}
	token := string(tokenBytes)

	// Установка соединения с сервером
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMusicalServiceClient(conn)

	// Добавление токена
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)

	// Создание запроса с введенным ID
	req := &pb.IDRequest{Id: trackID}

	stream, err := client.DownloadTrack(ctx, req)
	if err != nil {
		log.Fatalf("error calling DownloadTrack: %v", err)
	}

	// Открытие файла для записи
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer outFile.Close()

	fmt.Println("Receiving track...")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}

		chunk := resp.GetChunk()

		_, err = outFile.Write(chunk)
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
	}

	fmt.Printf("Track saved to %s\n", outputFile)
}