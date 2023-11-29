package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-playground/pb"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// server 接続
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure()) // 第二引数は通信が暗号化されずに非推奨。SSLにするべきだが学習なのでよし
	if err != nil {
		log.Fatalf("Fataled to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn) //FileServiceClientを取得できる
	//callListFiles(client)

	//callDownload(client)

	callUpload(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFiles())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Response from Download(bytes: %v", res.GetData())
		log.Printf("Response from Download(bytes: %v", string(res.GetData()))

	}
}

func callUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/s23300/learn/go-grpc-tutorial/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}
