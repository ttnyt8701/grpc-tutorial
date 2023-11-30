package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-playground/pb"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// 証明書 ssl
	certFile := "/Users/s23300/Library/Application Support/mkcert/rootCA.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	// server 接続
	conn, err := grpc.Dial("localhost:5003", grpc.WithTransportCredentials(creds)) // 第二引数は通信が暗号化されずに非推奨。SSLにするべきだが学習なのでよし
	if err != nil {
		log.Fatalf("Fataled to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn) //FileServiceClientを取得できる
	//callListFiles(client)

	callDownload(client)

	//callUpload(client)

	//callUploadAndNotifyProgress(client)
}

func callListFiles(client pb.FileServiceClient) {
	// gRPCリクエストのためのメタデータを新規に作成
	md := metadata.New(map[string]string{"authorization": "Bearer test-token"})
	// 新しいgRPCコンテキストを作成。このコンテキストにはメタデータ（md）が添付され、これにより認証情報がリクエストに含まれる
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//　gRPCクライアントがサーバーに対してListFilesというメソッドを呼び出している
	res, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	//res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFiles())
}

func callDownload(client pb.FileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // X秒後timeout deadline
	defer cancel()

	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("Error code: %v, ErrorMessage: %v", resErr.Code(), resErr.Message())
				} else if resErr.Code() == codes.DeadlineExceeded {
					log.Fatalln("deadline Exceeded")
				} else {
					log.Fatalln("unknown grpc error")
				}
			} else {
				log.Fatalln(err)
			}
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
	//defer file.Close()

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

func callUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/s23300/learn/go-grpc-tutorial/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// request
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatal(sendErr)
			}
			time.Sleep(1 * time.Second)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("recevide message: %v", res.GetMsg())
		}
	}()
	<-ch

}
