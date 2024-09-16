package main

import (
	"flag"
	"net"
	"runtime"
	"strconv"

	"github.com/zikunw/grpc-goproc-experiment/message"
	"google.golang.org/grpc"
)

var port int64
var gomaxprocs int

func main() {
	flag.Int64Var(&port, "port", 10000, "Port of the reciever (default=10000)")
	flag.IntVar(&gomaxprocs, "proc", 1, "Num for GOMAXPROCS config (default=1)")
	flag.Parse()

	runtime.GOMAXPROCS(gomaxprocs)

	s := grpc.NewServer()
	w := &Worker{}
	message.RegisterWorkerServer(s, w)
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(port, 10))
	if err != nil {
		panic(err)
	}
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type Worker struct {
	message.UnimplementedWorkerServer
}

func (w *Worker) Input(stream message.Worker_InputServer) error {
	for {
		_, err := stream.Recv()
		if err != nil {
			return err
		}
		stream.Send(&message.Response{})
	}
}
