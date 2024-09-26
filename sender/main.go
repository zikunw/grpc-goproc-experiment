package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"

	"github.com/zikunw/grpc-goproc-experiment/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Commandline input
var numReciever int
var gomaxprocs int
var numTuples int
var batchSize int
var stackTrace bool

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
	_, err = f.WriteString("\n")
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.IntVar(&numReciever, "num", 1, "Num of recievers (default=1)")
	flag.IntVar(&gomaxprocs, "proc", 1, "Num for GOMAXPROCS config (default=1)")
	flag.IntVar(&numTuples, "t", 1, "Num of tuples send (default=1)")
	flag.IntVar(&batchSize, "batch", 1, "Batch Size (default=1)")
	flag.BoolVar(&stackTrace, "trace", false, "Use stack trace (default=false)")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	println("Running sender with", numReciever, gomaxprocs, numTuples, batchSize)

	runtime.GOMAXPROCS(gomaxprocs)

	if stackTrace {
		f, _ := os.Create("trace.out")
		defer f.Close()
		trace.Start(f)
		defer trace.Stop()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	_, task := trace.NewTask(context.Background(), "Preparation")

	// Populate buffer
	buffer := &Buffer{}
	for i := 0; i < batchSize; i++ {
		buffer.Content = append(buffer.Content, message.KV{Key: int64(i), Value: int64(i)})
		buffer.Size += 1
	}

	task.End()

	// Start experiment
	var wg sync.WaitGroup
	for _, addr := range RECIEVER_ADDRS[:numReciever] {
		wg.Add(1)
		go run(addr, &wg, buffer, numTuples/numReciever)
	}

	wg.Wait()
}

func run(addr string, wg *sync.WaitGroup, buffer *Buffer, numTuples int) {

	conn := getConn(addr)
	stream, err := conn.Input(context.Background())
	if err != nil {
		panic(err)
	}

	_, task := trace.NewTask(context.Background(), "run()")

	count := 0
	batchCache := &message.Batch{}
	var i uint
	for i = 0; i < buffer.Size; i++ { // For each key-value pair in the buffer
		if len(batchCache.Kvs) == cap(batchCache.Kvs) {
			batchCache.Kvs = append(batchCache.Kvs, message.KVFromVTPool())
		} else {
			batchCache.Kvs = batchCache.Kvs[:len(batchCache.Kvs)+1]
			if batchCache.Kvs[len(batchCache.Kvs)-1] == nil {
				batchCache.Kvs[len(batchCache.Kvs)-1] = message.KVFromVTPool()
			}
		}
		batchCache.Kvs[i].Key = (*buffer).Content[i].Key
		batchCache.Kvs[i].Value = (*buffer).Content[i].Value
	}

	fmt.Println("Downstream started running with", addr, numTuples)

	start := time.Now()
	fmt.Println(addr, count+int(buffer.Size), numTuples)

	for count+int(buffer.Size) < numTuples {
		var err error
		err = stream.Send(batchCache)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
		_, err = stream.Recv()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}

		count += int(buffer.Size)
	}
	fmt.Println("Downstream finished")
	duration := time.Since(start)
	stream.CloseSend()
	go shutdown(addr)

	task.End()

	_, err = f.WriteString(fmt.Sprintf("%d,%d,%s,%d\n", numReciever, gomaxprocs, addr, duration))
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

func getConn(addr string) message.WorkerClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return message.NewWorkerClient(conn)
}
