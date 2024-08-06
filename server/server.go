package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "github.com/Aravinthvvs/gRPC/proto/train/train"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTicketServiceServer
	mu          sync.Mutex
	receipts    map[string]*pb.ReceiptResponse
	userSeats   map[string]string
	sectionA    map[string]string
	sectionB    map[string]string
	seatCounter int
}

func newServer() *server {
	return &server{
		receipts:    make(map[string]*pb.ReceiptResponse),
		userSeats:   make(map[string]string),
		sectionA:    make(map[string]string),
		sectionB:    make(map[string]string),
		seatCounter: 0, // Initialize the seat counter
	}
}

func (s *server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Basic validation
	if req.User == nil {
		return nil, fmt.Errorf("user information is required")
	}
	if req.From == "" || req.To == "" {
		return nil, fmt.Errorf("from and to fields are required")
	}

	// Generate a receipt ID
	receiptID := fmt.Sprintf("rec-%d", len(s.receipts)+1)

	// Generate a unique seat number
	s.seatCounter++
	seat := fmt.Sprintf("Seat-%d", s.seatCounter)

	// Store receipt and user seat allocation
	s.receipts[receiptID] = &pb.ReceiptResponse{
		From:      req.From,
		To:        req.To,
		User:      req.User,
		PricePaid: 20,
		Seat:      seat,
	}
	s.userSeats[req.User.Email] = seat

	// Define sections
	sections := []string{"SectionA", "SectionB"}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Randomly select a section
	section := sections[rand.Intn(len(sections))]

	if section == "SectionA" {
		s.sectionA[req.User.Email] = seat
	} else {
		s.sectionB[req.User.Email] = seat
	}

	return &pb.PurchaseResponse{ReceiptId: receiptID}, nil
}

func (s *server) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.receipts[req.ReceiptId]
	if !exists {
		return nil, fmt.Errorf("receipt not found")
	}

	return receipt, nil
}

func (s *server) ViewUsersBySection(ctx context.Context, req *pb.ViewUsersRequest) (*pb.ViewUsersResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userSeats map[string]string
	if req.Section == "SectionA" {
		userSeats = s.sectionA
	} else if req.Section == "SectionB" {
		userSeats = s.sectionB
	} else {
		return nil, fmt.Errorf("invalid section")
	}

	var userSeatList []*pb.UserSeat
	for email, seat := range userSeats {
		userSeatList = append(userSeatList, &pb.UserSeat{
			User: &pb.User{
				Email: email,
			},
			Seat: seat,
		})
	}

	return &pb.ViewUsersResponse{UserSeats: userSeatList}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.userSeats[req.Email]
	if !exists {
		return &pb.RemoveUserResponse{Success: false}, nil
	}

	// Remove from sections
	delete(s.sectionA, req.Email)
	delete(s.sectionB, req.Email)

	delete(s.userSeats, req.Email)

	// Remove receipt
	for receiptID, receipt := range s.receipts {
		if receipt.User.Email == req.Email {
			delete(s.receipts, receiptID)
			break
		}
	}

	return &pb.RemoveUserResponse{Success: true}, nil
}

func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.userSeats[req.Email]
	if !exists {
		return &pb.ModifySeatResponse{Success: false}, nil
	}

	// Update seat in the corresponding section
	section := "SectionA"
	if _, ok := s.sectionB[req.Email]; ok {
		section = "SectionB"
	}

	if section == "SectionA" {
		s.sectionA[req.Email] = req.NewSeat
	} else {
		s.sectionB[req.Email] = req.NewSeat
	}

	// Update receipt with new seat
	for _, receipt := range s.receipts {
		if receipt.User.Email == req.Email {
			receipt.Seat = req.NewSeat
			break
		}
	}

	s.userSeats[req.Email] = req.NewSeat

	return &pb.ModifySeatResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTicketServiceServer(s, newServer())
	log.Println("Starting server on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
