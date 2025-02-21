package main

import (
	common "github.com/quanbin27/commons"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", ":2000")
)

func main() {
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()
	NewGRPCHandler(grpcServer)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
