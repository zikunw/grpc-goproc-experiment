// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.0
// source: message/rpc.proto

package message

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

// WorkerClient is the client API for Worker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkerClient interface {
	// APIs for workers to interact with each other
	Input(ctx context.Context, opts ...grpc.CallOption) (Worker_InputClient, error)
}

type workerClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkerClient(cc grpc.ClientConnInterface) WorkerClient {
	return &workerClient{cc}
}

func (c *workerClient) Input(ctx context.Context, opts ...grpc.CallOption) (Worker_InputClient, error) {
	stream, err := c.cc.NewStream(ctx, &Worker_ServiceDesc.Streams[0], "/message.Worker/Input", opts...)
	if err != nil {
		return nil, err
	}
	x := &workerInputClient{stream}
	return x, nil
}

type Worker_InputClient interface {
	Send(*Batch) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type workerInputClient struct {
	grpc.ClientStream
}

func (x *workerInputClient) Send(m *Batch) error {
	return x.ClientStream.SendMsg(m)
}

func (x *workerInputClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WorkerServer is the server API for Worker service.
// All implementations must embed UnimplementedWorkerServer
// for forward compatibility
type WorkerServer interface {
	// APIs for workers to interact with each other
	Input(Worker_InputServer) error
	mustEmbedUnimplementedWorkerServer()
}

// UnimplementedWorkerServer must be embedded to have forward compatible implementations.
type UnimplementedWorkerServer struct {
}

func (UnimplementedWorkerServer) Input(Worker_InputServer) error {
	return status.Errorf(codes.Unimplemented, "method Input not implemented")
}
func (UnimplementedWorkerServer) mustEmbedUnimplementedWorkerServer() {}

// UnsafeWorkerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkerServer will
// result in compilation errors.
type UnsafeWorkerServer interface {
	mustEmbedUnimplementedWorkerServer()
}

func RegisterWorkerServer(s grpc.ServiceRegistrar, srv WorkerServer) {
	s.RegisterService(&Worker_ServiceDesc, srv)
}

func _Worker_Input_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WorkerServer).Input(&workerInputServer{stream})
}

type Worker_InputServer interface {
	Send(*Response) error
	Recv() (*Batch, error)
	grpc.ServerStream
}

type workerInputServer struct {
	grpc.ServerStream
}

func (x *workerInputServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *workerInputServer) Recv() (*Batch, error) {
	m := new(Batch)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Worker_ServiceDesc is the grpc.ServiceDesc for Worker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Worker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.Worker",
	HandlerType: (*WorkerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Input",
			Handler:       _Worker_Input_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "message/rpc.proto",
}
