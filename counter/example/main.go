package main

import (
	"log"
	"net"
	"sync"

	pb "github.com/bakaoh/fresher/counter"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":56789"
)

type server struct{}

var m = map[string]int64{}
var lock = sync.Mutex{}

// SetBalance ...
func (s *server) SetBalance(ctx context.Context, in *pb.UserReq) (*pb.BalanceRes, error) {
	lock.Lock()
	defer lock.Unlock()
	m[in.UserId] = in.Balance
	return &pb.BalanceRes{Balance: in.Balance}, nil
}

// GetBalance ...
func (s *server) GetBalance(ctx context.Context, in *pb.UserReq) (*pb.BalanceRes, error) {
	lock.Lock()
	defer lock.Unlock()
	if b, ok := m[in.UserId]; ok {
		return &pb.BalanceRes{Balance: b}, nil
	}
	return &pb.BalanceRes{Balance: 0}, nil
}

// IncreaseBalance ...
func (s *server) IncreaseBalance(ctx context.Context, in *pb.UserReq) (*pb.BalanceRes, error) {
	lock.Lock()
	defer lock.Unlock()
	b, ok := m[in.UserId]
	if !ok {
		b = 0
	}
	m[in.UserId] = b + in.Balance
	return &pb.BalanceRes{Balance: m[in.UserId]}, nil
}

// DecreaseBalance ...
func (s *server) DecreaseBalance(ctx context.Context, in *pb.UserReq) (*pb.BalanceRes, error) {
	lock.Lock()
	defer lock.Unlock()
	b, ok := m[in.UserId]
	if !ok {
		b = 0
	}
	m[in.UserId] = b - in.Balance
	return &pb.BalanceRes{Balance: m[in.UserId]}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCounterServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
