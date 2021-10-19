package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"

	"github.com/containerd/ttrpc"
	"github.com/kevpar/ttrpc-deadlock/svc"
)

type connType struct {
	listen func() (net.Listener, error)
	dial   func() (net.Conn, error)
}

func listen(typ string) (net.Listener, error) {
	return connTypes[typ].listen()
}

func dial(typ string) (net.Conn, error) {
	return connTypes[typ].dial()
}

var flagConnType = flag.String("conn-type", "", "")
var flagWorkers = flag.Int("workers", 4, "")

const (
	dataSizeKB = 200
)

var bigString string

func init() {
	var b [dataSizeKB * 1024]byte
	for i := range b {
		b[i] = 0x41
	}
	bigString = string(b[:])
}

type service struct{}

func (s *service) Foo(ctx context.Context, req *svc.FooRequest) (*svc.FooResponse, error) {
	return &svc.FooResponse{
		S: bigString,
	}, nil
}

func stackDumper(path string) {
	r := bufio.NewReader(os.Stdin)
	for {
		_, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		var (
			buf       []byte
			stackSize int
		)
		bufferLen := 16384
		for stackSize == len(buf) {
			buf = make([]byte, bufferLen)
			stackSize = runtime.Stack(buf, true)
			bufferLen *= 2
		}
		if err := ioutil.WriteFile(path, buf[:stackSize], 0664); err != nil {
			panic(err)
		}
		fmt.Printf("Wrote stacks to %s\n", path)
	}
}

func runServer() error {
	l, err := listen(*flagConnType)
	if err != nil {
		return err
	}
	s, err := ttrpc.NewServer()
	if err != nil {
		return err
	}
	service := &service{}
	svc.RegisterSvcService(s, service)
	go stackDumper("ServerStacks.txt")
	if err := s.Serve(context.Background(), l); err != nil {
		return err
	}
	return nil
}

func runClient() error {
	c, err := dial(*flagConnType)
	if err != nil {
		return err
	}
	client := svc.NewSvcClient(ttrpc.NewClient(c))
	go stackDumper("ClientStacks.txt")
	f := func(i int) {
		for {
			_, err := client.Foo(context.Background(), &svc.FooRequest{
				S: bigString,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(i)
		}
	}
	for i := 0; i < *flagWorkers; i++ {
		go f(i)
	}
	for {
	}
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "server":
		if err := runServer(); err != nil {
			panic(err)
		}
	case "client":
		if err := runClient(); err != nil {
			panic(err)
		}
	}
}
