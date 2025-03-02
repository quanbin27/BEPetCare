package main

import (
	"github.com/quanbin27/commons/config"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	dsn := config.Envs.RecordsDSN
	store, err := NewMongoStore(dsn)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	log.Println("Connecting to database ...", dsn)
	grpcAddr := config.Envs.RecordsGrpcAddr
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("failed to listen")
	}
	defer l.Close()
	service := NewPetRecordService(store)
	NewGrpcHandler(grpcServer, service)
	log.Println("Pet Record Service Listening on", grpcAddr)
	grpcServer.Serve(l)
}
