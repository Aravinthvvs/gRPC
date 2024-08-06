package main

import (
	"context"
	"testing"

	pb "github.com/Aravinthvvs/gRPC/proto/train/train"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockTicketServiceClient is a mock implementation of TicketServiceClient
type MockTicketServiceClient struct {
	mock.Mock
}

func (m *MockTicketServiceClient) PurchaseTicket(ctx context.Context, in *pb.PurchaseRequest, opts ...grpc.CallOption) (*pb.PurchaseResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.PurchaseResponse), args.Error(1)
}

func (m *MockTicketServiceClient) GetReceipt(ctx context.Context, in *pb.ReceiptRequest, opts ...grpc.CallOption) (*pb.ReceiptResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ReceiptResponse), args.Error(1)
}

func (m *MockTicketServiceClient) ViewUsersBySection(ctx context.Context, in *pb.ViewUsersRequest, opts ...grpc.CallOption) (*pb.ViewUsersResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ViewUsersResponse), args.Error(1)
}

func (m *MockTicketServiceClient) RemoveUser(ctx context.Context, in *pb.RemoveUserRequest, opts ...grpc.CallOption) (*pb.RemoveUserResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.RemoveUserResponse), args.Error(1)
}

func (m *MockTicketServiceClient) ModifySeat(ctx context.Context, in *pb.ModifySeatRequest, opts ...grpc.CallOption) (*pb.ModifySeatResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ModifySeatResponse), args.Error(1)
}

// TestPurchaseTicketClient tests the PurchaseTicket function
func TestPurchaseTicketClient(t *testing.T) {
	// Create a new instance of MockTicketServiceClient
	mockClient := new(MockTicketServiceClient)

	// Define the expected response and error
	expectedResponse := &pb.PurchaseResponse{ReceiptId: "rec-123"}
	mockClient.On("PurchaseTicket", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Create a mock gRPC connection (just for the sake of having a valid connection instance)
	conn := &grpc.ClientConn{}
	client := NewClient(conn) // Use the NewClient function
	client.client = mockClient

	// Define the user and request parameters
	user := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	from := "London"
	to := "France"

	// Call the PurchaseTicket method
	resp, err := client.PurchaseTicket(context.Background(), from, to, user)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert that the response ReceiptId is as expected
	assert.Equal(t, expectedResponse.ReceiptId, resp.ReceiptId)

	// Verify that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestGetReceipt tests the GetReceipt method of the client
func TestGetReceipt(t *testing.T) {
	// Create a new instance of the mock client
	mockClient := new(MockTicketServiceClient)

	// Define the expected response and behavior for the mock
	expectedReceipt := &pb.ReceiptResponse{
		From:      "London",
		To:        "France",
		User:      &pb.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		PricePaid: 20,
		Seat:      "Seat-1",
	}

	mockClient.On("GetReceipt", mock.Anything, &pb.ReceiptRequest{ReceiptId: "rec-1"}).
		Return(expectedReceipt, nil)

	// Create a new client using the mock
	client := &Client{client: mockClient}

	// Call the method under test
	resp, err := client.GetReceipt(context.Background(), "rec-1")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedReceipt, resp)

	// Verify that the mock expectations were met
	mockClient.AssertExpectations(t)
}

// TestRemoveUser tests the RemoveUser method of the client
func TestRemoveUser(t *testing.T) {
	mockClient := new(MockTicketServiceClient)
	expectedResponse := &pb.RemoveUserResponse{Success: true}

	mockClient.On("RemoveUser", mock.Anything, &pb.RemoveUserRequest{Email: "alice.smith@example.com"}).
		Return(expectedResponse, nil)

	client := &Client{client: mockClient}
	resp, err := client.RemoveUser(context.Background(), "alice.smith@example.com")

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	mockClient.AssertExpectations(t)
}

// TestModifySeat tests the ModifySeat method of the client
func TestModifySeat(t *testing.T) {
	mockClient := new(MockTicketServiceClient)
	expectedResponse := &pb.ModifySeatResponse{Success: true}

	mockClient.On("ModifySeat", mock.Anything, &pb.ModifySeatRequest{Email: "alice.smith@example.com", NewSeat: "Seat-42"}).
		Return(expectedResponse, nil)

	client := &Client{client: mockClient}
	resp, err := client.ModifySeat(context.Background(), "alice.smith@example.com", "Seat-42")

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	mockClient.AssertExpectations(t)
}
