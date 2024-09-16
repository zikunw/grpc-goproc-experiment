package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/zikunw/grpc-goproc-experiment/message"
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
	for _, stream := range reciever_streams {
		go run(stream, timeChan, buffer, numTuples/numReciever)
	}

	for range reciever_streams {
		fmt.Println(<-timeChan)
	}
}

func run(stream *BatchStream, timeChan chan SETime, buffer *Buffer, numTuples int) {
	count := 0
	start_time := time.Now()
	//=======Critical Section=====
	for count+len(buffer.Content) < numTuples {
		count += int(buffer.Size)
		stream.Put(buffer)
	}
	//============================
	end_time := time.Now()
	timeChan <- SETime{
		start: start_time,
		end:   end_time,
	}
	stream.Close()
}
