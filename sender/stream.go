package main

import (
	"context"
	"io"
	"sync"

	message "github.com/zikunw/grpc-goproc-experiment/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ========================
// Batch implementation
// ========================

// BatchStream is a stream between workers that sends batches of workloads.
// There is one BatchStream per downstream worker, meaning that each worker
// has multiple BatchStreams.
type BatchStream struct {
	mu sync.Mutex

	stream     message.Worker_InputClient
	batchCache *message.Batch

	addr string
}

// NewBatchStream creates a new batch stream to a worker located at addr.
func NewBatchStream(addr string) *BatchStream { // TODO[#36]: Should this take an utils.Address?
	if addr == "" {
		panic("BatchStream's address argument cannot be empty")
	}

	conn := getConn(addr)
	stream, err := conn.Input(context.Background())
	if err != nil {
		panic(err)
	}

	return &BatchStream{
		stream:     stream,
		batchCache: &message.Batch{},
		addr:       addr,
	}
}

// Put sends a batch of workloads down the stream.
// Put is thread-safe.
func (b *BatchStream) Put(kvs *Buffer) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var err error
	var i uint
	for i = 0; i < kvs.Size; i++ { // For each key-value pair in the buffer
		if len(b.batchCache.Kvs) == cap(b.batchCache.Kvs) {
			b.batchCache.Kvs = append(b.batchCache.Kvs, message.KVFromVTPool())
		} else {
			b.batchCache.Kvs = b.batchCache.Kvs[:len(b.batchCache.Kvs)+1]
			if b.batchCache.Kvs[len(b.batchCache.Kvs)-1] == nil {
				b.batchCache.Kvs[len(b.batchCache.Kvs)-1] = message.KVFromVTPool()
			}
		}

		// Move the key-value pair, instead of the workload object to avoid heap allocation.
		// This lets us reuse the workload objects.

		// We are repeatedly dereferencing the pointer to the workload object because `kv := (*kvs)[i]`
		// copies a lock, which isn't safe to do in general, but I think it might be safe here? - Ezra

		b.batchCache.Kvs[i].Key = (*kvs).Content[i].Key
		b.batchCache.Kvs[i].Value = (*kvs).Content[i].Value
	}

	// Send the batch down the stream.
	err = b.stream.Send(b.batchCache)
	if err != nil {
		if err == io.EOF {
			// TODO[#37]: handle EOF
			return nil
		}
		return err
	}

	// Wait to receive OK from the downstream worker.
	// This means that calls to Put are blocked until the downstream worker
	// has stored the batch in its input buffer.
	_, err = b.stream.Recv()
	if err != nil {
		if err == io.EOF {
			// TODO[#37]: handle EOF
			return nil
		}
		return err
	}

	return nil
}

// Close closes the stream.
func (b *BatchStream) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.stream.CloseSend()
}

func getConn(addr string) message.WorkerClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return message.NewWorkerClient(conn)
}
