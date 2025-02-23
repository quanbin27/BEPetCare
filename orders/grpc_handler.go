package main

import (
	"context"
	pb "github.com/quanbin27/commons/genproto"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGRPCHandler(grpcServer *grpc.Server) {
	handler := &grpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}
func (*grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Printf("CreateOrder called with req: %v", req)
	res := &pb.CreateOrderResponse{
		Status: "added",
	}
	return res, nil
}
