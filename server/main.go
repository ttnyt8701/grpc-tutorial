package main

import (
	"bytes"
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-playground/pb"
	"io"
	"log"
	"net"
	"os"
	"time"
)

/*
Unary server の実装
req/resは1:1
*/

//import (
//
//	"fmt"
//	"google.golang.org/grpc"
//	grpc_midllerware "github.com/grpc-ecosystem/go-grpc-middleware"
//	//grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/auth"
//	"log"
//	"net"
//	"os"
//)

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

	// ファイルの存在チェック
	//if _, err := os.Stat(path); os.IsNotExist(err) {
	//	return status.Error(codes.NotFound, "file was not found")
	//}

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

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("recieved data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("recieved %vbytes", size),
		}
		err = stream.Send(res)
	}
}

// interceptor
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

func authorize(ctx context.Context) (context.Context, error) {
	// クライアントから送信されたgRPCリクエストのコンテキスト（ctx）から「Bearer」トークンを抽出
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	// 認証
	if token != "test-token!!!" {
		//return nil, errors.New("bad token")
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}

	return ctx, nil
}

func main() {
	// 証明書
	creds, err := credentials.NewServerTLSFromFile(
		"ssl/localhost.pem",
		"ssl/localhost-key.pem")
	if err != nil {
		log.Fatalln(err)
	}

	// gRPCサーバーのインスタンス作成
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			myLogging(),
			grpc_auth.UnaryServerInterceptor(authorize),
		),
		),
	) // interceptor追加

	// gRPCサーバーにサービス登録
	pb.RegisterFileServiceServer(s, &server{})

	// server設定
	lis, err := net.Listen("tcp", "localhost:5003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// server起動
	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
