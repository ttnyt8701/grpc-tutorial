// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.0
// source: proto/file.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServiceClient interface {
	// Unary RPC 1:1
	ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesResponse, error)
	// server stream
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FileService_DownloadClient, error)
	// client stream
	Upload(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadClient, error)
	// bidirectional stream
	UploadAndNotifyProgress(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadAndNotifyProgressClient, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesResponse, error) {
	out := new(ListFilesResponse)
	err := c.cc.Invoke(ctx, "/FileService/ListFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FileService_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[0], "/FileService/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileService_DownloadClient interface {
	Recv() (*DownloadResponse, error)
	grpc.ClientStream
}

type fileServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *fileServiceDownloadClient) Recv() (*DownloadResponse, error) {
	m := new(DownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[1], "/FileService/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceUploadClient{stream}
	return x, nil
}

type FileService_UploadClient interface {
	Send(*UploadRequest) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type fileServiceUploadClient struct {
	grpc.ClientStream
}

func (x *fileServiceUploadClient) Send(m *UploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServiceUploadClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServiceClient) UploadAndNotifyProgress(ctx context.Context, opts ...grpc.CallOption) (FileService_UploadAndNotifyProgressClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[2], "/FileService/UploadAndNotifyProgress", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServiceUploadAndNotifyProgressClient{stream}
	return x, nil
}

type FileService_UploadAndNotifyProgressClient interface {
	Send(*UploadAndNotifyProgressRequest) error
	Recv() (*UploadAndNotifyProgressResponse, error)
	grpc.ClientStream
}

type fileServiceUploadAndNotifyProgressClient struct {
	grpc.ClientStream
}

func (x *fileServiceUploadAndNotifyProgressClient) Send(m *UploadAndNotifyProgressRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServiceUploadAndNotifyProgressClient) Recv() (*UploadAndNotifyProgressResponse, error) {
	m := new(UploadAndNotifyProgressResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations must embed UnimplementedFileServiceServer
// for forward compatibility
type FileServiceServer interface {
	// Unary RPC 1:1
	ListFiles(context.Context, *ListFilesRequest) (*ListFilesResponse, error)
	// server stream
	Download(*DownloadRequest, FileService_DownloadServer) error
	// client stream
	Upload(FileService_UploadServer) error
	// bidirectional stream
	UploadAndNotifyProgress(FileService_UploadAndNotifyProgressServer) error
	mustEmbedUnimplementedFileServiceServer()
}

// UnimplementedFileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileServiceServer struct {
}

func (UnimplementedFileServiceServer) ListFiles(context.Context, *ListFilesRequest) (*ListFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFiles not implemented")
}
func (UnimplementedFileServiceServer) Download(*DownloadRequest, FileService_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedFileServiceServer) Upload(FileService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFileServiceServer) UploadAndNotifyProgress(FileService_UploadAndNotifyProgressServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadAndNotifyProgress not implemented")
}
func (UnimplementedFileServiceServer) mustEmbedUnimplementedFileServiceServer() {}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_ListFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).ListFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/FileService/ListFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).ListFiles(ctx, req.(*ListFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).Download(m, &fileServiceDownloadServer{stream})
}

type FileService_DownloadServer interface {
	Send(*DownloadResponse) error
	grpc.ServerStream
}

type fileServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *fileServiceDownloadServer) Send(m *DownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FileService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).Upload(&fileServiceUploadServer{stream})
}

type FileService_UploadServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*UploadRequest, error)
	grpc.ServerStream
}

type fileServiceUploadServer struct {
	grpc.ServerStream
}

func (x *fileServiceUploadServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServiceUploadServer) Recv() (*UploadRequest, error) {
	m := new(UploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileService_UploadAndNotifyProgress_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).UploadAndNotifyProgress(&fileServiceUploadAndNotifyProgressServer{stream})
}

type FileService_UploadAndNotifyProgressServer interface {
	Send(*UploadAndNotifyProgressResponse) error
	Recv() (*UploadAndNotifyProgressRequest, error)
	grpc.ServerStream
}

type fileServiceUploadAndNotifyProgressServer struct {
	grpc.ServerStream
}

func (x *fileServiceUploadAndNotifyProgressServer) Send(m *UploadAndNotifyProgressResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServiceUploadAndNotifyProgressServer) Recv() (*UploadAndNotifyProgressRequest, error) {
	m := new(UploadAndNotifyProgressRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListFiles",
			Handler:    _FileService_ListFiles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _FileService_Download_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Upload",
			Handler:       _FileService_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadAndNotifyProgress",
			Handler:       _FileService_UploadAndNotifyProgress_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/file.proto",
}
