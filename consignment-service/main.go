package main

import (
	"sync"
	"context"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	pb "github.com/Nishant967/shipping-microservice-project/consignment-service/proto/consignment"
	"google.golang.org/grpc"

)

const (
	port = ":5000"
)

type inventory interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Inventory struct {
	consignments []*pb.Consignment
	mutex sync.RWMutex
}

// Create a new consignment
func (inv *Inventory) Create(consignment *pb.Consignment) (*pb.Consignment,error) {
	inv.mutex.Lock()
	updated := append(inv.consignments, consignment)
	inv.consignments = updated
	inv.mutex.Unlock()
	return consignment, nil
}

// GetAll consignments
func (inv *Inventory) GetAll() []*pb.Consignment {
	return inv.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.

type service struct {
	inven inventory
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	// Save our consignment
	consignment, err := s.inven.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// GetConsignments -
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.inven.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {

	inv := &Inventory{}

	// Set-up our gRPC server.
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(s, &service{inv})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}