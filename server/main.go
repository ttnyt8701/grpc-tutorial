package main

import (
	"bytes"
	"grpc-playground/pb"
	"io"
	"time"
)

/*
Unary server の実装
req/resは1:1
*/

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type server struct {
	// FileServiceServerインターフェースを満したデフォルト実装. オーバーライドして具体的な実装をする必要がある。
	pb.UnimplementedFileServiceServer
}

// ListFiles
func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := "/Users/s23300/learn/go-grpc-tutorial/storage"

	paths, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Files: filenames,
	}

	return res, nil
}

// server stream
func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := "/Users/s23300/learn/go-grpc-tutorial/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	buf := make([]byte, 16)

	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res) // クライアントにデータを送る
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// client stream
func (s *server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer // アップロードされたデータを格納するBuffer
	for {
		req, err := stream.Recv() // クライアントから複数のリクエストを取得
		if err == io.EOF {        // 終了したとき
			res := &pb.UploadResponse{Size: int32(buf.Len())} // bufのサイズ
			return stream.SendAndClose(res)                   // serverからのレスポンスを返す
		}
		if err != nil {
			return err
		}

		data := req.GetData() // リクエストからのデータを格納
		log.Printf("received data(bytes): %v", req.GetData())
		log.Printf("received data(bytes): %v", string(req.GetData()))
		buf.Write(data) // 出力
	}
}

// middleware
// logging
func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("request data: %v", req)

		resp, err = handler(ctx, req) // クライアントからのコールされたメソッド
		if err != nil {
			return nil, err
		}
		log.Printf("response data: %v", resp)

		return resp, nil
	}
}

func main() {
	// gRPCサーバーのインスタンス作成
	s := grpc.NewServer(grpc.UnaryInterceptor(myLogging())) // interceptor追加

	// gRPCサーバーにサービス登録
	pb.RegisterFileServiceServer(s, &server{})

	// server設定
	lis, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// server起動
	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
