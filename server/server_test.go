package main

import (
	"context"
	"testing"

	pb "github.com/Aravinthvvs/gRPC/proto/train/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestServer() *server {
	return &server{
		receipts:    make(map[string]*pb.ReceiptResponse),
		userSeats:   make(map[string]string),
		sectionA:    make(map[string]string),
		sectionB:    make(map[string]string),
		seatCounter: 0, // Initialize the seat counter
	}

}
func TestPurchaseTicket(t *testing.T) {
	server := newTestServer()
	tests := []struct {
		name            string
		req             *pb.PurchaseRequest
		expectedResp    *pb.PurchaseResponse
		expectedErrCode codes.Code
	}{
		{
			name: "successful purchase",
			req: &pb.PurchaseRequest{
				From: "London",
				To:   "France",
				User: &pb.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				},
			},
			expectedResp:    &pb.PurchaseResponse{ReceiptId: "rec-1"},
			expectedErrCode: codes.OK,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.PurchaseTicket(context.Background(), tt.req)
			if err != nil {
				st, _ := status.FromError(err)
				assert.Equal(t, tt.expectedErrCode, st.Code())
			}
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

// Test GetReceipt
func TestGetReceipt(t *testing.T) {
	s := newTestServer()
	// First, simulate a ticket purchase to generate a receipt
	req := &pb.PurchaseRequest{
		From: "London",
		To:   "France",
		User: &pb.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane.doe@example.com",
		},
	}
	purchaseResp, err := s.PurchaseTicket(context.Background(), req)
	require.NoError(t, err)

	receiptId := purchaseResp.ReceiptId
	getResp, err := s.GetReceipt(context.Background(), &pb.ReceiptRequest{ReceiptId: receiptId})
	require.NoError(t, err)
	assert.Equal(t, req.From, getResp.From)
	assert.Equal(t, req.To, getResp.To)
	assert.Equal(t, req.User, getResp.User)
}

type testInput struct {
	name            string
	initialUser     *pb.User
	removeEmail     string
	expectedSuccess bool
}

func TestRemoveUser(t *testing.T) {
	s := newTestServer()

	tests := []testInput{
		{
			name: "Successful removal",
			initialUser: &pb.User{
				FirstName: "Bob",
				LastName:  "Brown",
				Email:     "bob.brown@example.com",
			},
			removeEmail:     "bob.brown@example.com",
			expectedSuccess: true,
		},
		{
			name:            "Remove non-existent user",
			initialUser:     nil, // No initial user
			removeEmail:     "nonexistent.user@example.com",
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.initialUser != nil {
				// Simulate a ticket purchase for the initial user
				req := &pb.PurchaseRequest{
					From: "London",
					To:   "France",
					User: tt.initialUser,
				}
				_, err := s.PurchaseTicket(context.Background(), req)
				require.NoError(t, err)
			}

			// Test removal
			removeResp, err := s.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: tt.removeEmail})
			require.NoError(t, err)
			assert.Equal(t, tt.expectedSuccess, removeResp.Success)
		})
	}
}
