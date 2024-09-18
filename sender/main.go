package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/zikunw/grpc-goproc-experiment/message"
	"google.golang.org/grpc"
)

// Commandline input
var numReciever int
var gomaxprocs int
var numTuples int
var batchSize int

// Reciever list (index=#reciever-1)
var RECIEVER_ADDRS = []string{
	"localhost:10000",
	"localhost:10001",
	"localhost:10002",
	"localhost:10003",
	"localhost:10004",
	"localhost:10005",
	"localhost:10006",
	"localhost:10007",
	"localhost:10008",
	"localhost:10009",
	"localhost:10010",
	"localhost:10011",
	"localhost:10012",
}

var f *os.File

func init() {
	var err error
	f, err = os.OpenFile("./result.txt", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}
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
	var wg sync.WaitGroup
	for i, stream := range reciever_streams {
		wg.Add(1)
		go run(stream, &wg, buffer, numTuples/numReciever, recievers[i])
	}

	wg.Wait()
}

func run(stream *BatchStream, wg *sync.WaitGroup, buffer *Buffer, numTuples int, addr string) {
	count := 0
	fmt.Println("Downstream started running")
	start := time.Now()
	for count+len(buffer.Content) < numTuples {
		count += int(buffer.Size)
		err := stream.Put(buffer)
		if err != nil {
			panic(err)
		}
	}
	duration := time.Since(start)
	stream.Close()
	shutdown(addr)

	_, err := f.WriteString(fmt.Sprintf("%d,%d,%s,%d\n", numReciever, gomaxprocs, addr, duration))
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func shutdown(address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := message.NewWorkerClient(conn)
	client.Shutdown(context.Background(), &message.Empty{})
}
