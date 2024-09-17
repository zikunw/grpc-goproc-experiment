package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/zikunw/grpc-goproc-experiment/message"
	"google.golang.org/grpc"
)

// Commandline input
var numReciever int
var gomaxprocs int
var numTuples int
var batchSize int

// SETime
type SETime = struct {
	start time.Time
	end   time.Time
}

// Reciever list (index=#reciever-1)
var RECIEVER_ADDRS = []string{
	"localhost:10000",
	"localhost:10001",
	"localhost:10002",
	"localhost:10003",
	"localhost:10004",
	"localhost:10005",
}

func main() {
	flag.IntVar(&numReciever, "num", 1, "Num of recievers (default=1)")
	flag.IntVar(&gomaxprocs, "proc", 1, "Num for GOMAXPROCS config (default=1)")
	flag.IntVar(&numTuples, "t", 1, "Num of tuples send (default=1)")
	flag.IntVar(&batchSize, "batch", 1, "Batch Size (default=1)")
	flag.Parse()

	runtime.GOMAXPROCS(gomaxprocs)

	// Create downstreams
	recievers := RECIEVER_ADDRS[:numReciever]
	reciever_streams := []*BatchStream{}
	for _, reciever := range recievers {
		reciever_streams = append(reciever_streams, NewBatchStream(reciever))
	}

	// Populate buffer
	buffer := &Buffer{}
	for i := 0; i < numTuples; i++ {
		buffer.Content = append(buffer.Content, message.KV{Key: int64(i), Value: int64(i)})
		buffer.Size += 1
	}

	// Start experiment
	timeChan := make(chan SETime, 10)
	for i, stream := range reciever_streams {
		go run(stream, timeChan, buffer, numTuples/numReciever, recievers[i])
	}

	f, err := os.Create(path.Join("./result.txt"))
	if err != nil {
		panic(err)
	}

	for range reciever_streams {
		result := <-timeChan
		_, err = f.WriteString(fmt.Sprintf("%d\n", result.end.Sub(result.start)))
		if err != nil {
			panic(err)
		}
	}
	f.Close()
}

func run(stream *BatchStream, timeChan chan SETime, buffer *Buffer, numTuples int, addr string) {
	count := 0
	fmt.Println("Downstream started running")
	start_time := time.Now()
	//=======Critical Section=====
	for count+len(buffer.Content) < numTuples {
		count += int(buffer.Size)
		err := stream.Put(buffer)
		if err != nil {
			panic(err)
		}
	}
	//============================
	end_time := time.Now()
	timeChan <- SETime{
		start: start_time,
		end:   end_time,
	}
	stream.Close()
	shutdown(addr)
}

func shutdown(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := message.NewWorkerClient(conn)
	client.Shutdown(context.Background(), &message.Empty{})
}
