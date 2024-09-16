package main

import (
	"github.com/zikunw/grpc-goproc-experiment/message"
)

type Buffer struct {
	Content []message.KV
	Size    uint
}
