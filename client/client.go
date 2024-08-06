package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Aravinthvvs/gRPC/proto/train/train"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTicketServiceClient(conn)

	// Example of Purchase Ticket
	users := []pb.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
		{
			FirstName: "Aravinth",
			LastName:  "Veeramuthu",
			Email:     "aveeramu@example.com",
		},
	}
	for _, user := range users {
		userPtr := &user
		purchaseResp, err := c.PurchaseTicket(context.Background(), &pb.PurchaseRequest{
			From: "London",
			To:   "France",
			User: userPtr,
		})
		if err != nil {
			log.Fatalf("could not purchase ticket: %v", err)
		}
		fmt.Printf("Purchase Response: %s\n", purchaseResp.ReceiptId)

		// Example of Get Receipt
		receiptResp, err := c.GetReceipt(context.Background(), &pb.ReceiptRequest{
			ReceiptId: purchaseResp.ReceiptId,
		})
		if err != nil {
			log.Fatalf("could not get receipt: %v", err)
		}
		fmt.Printf("Receipt: %+v\n", receiptResp)
	}

	// Example of Remove User
	resp, err := c.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: "aveeramu@example.com"})
	if err != nil {
		log.Fatalf("could not remove user: %v", err)
	}
	if resp.Success {
		fmt.Println("User removed successfully.")
	} else {
		fmt.Println("Failed to remove user.")
	}

	// Example of Modify Seat
	modresp, err := c.ModifySeat(context.Background(), &pb.ModifySeatRequest{Email: "john.doe@example.com", NewSeat: "Seat-1"})
	if err != nil {
		log.Fatalf("could not modify seat: %v", err)
	}
	if modresp.Success {
		fmt.Println("Seat modified successfully.")
	} else {
		fmt.Println("Failed to modify seat.")
	}

	// Example of View Users By Sections
	seatsResp, err := c.ViewUsersBySection(context.Background(), &pb.ViewUsersRequest{Section: "SectionA"})

	if err != nil {
		log.Fatalf("Could not get the seats deatils for the Section : %v", "SectionA")
	}
	fmt.Printf("Seats Deatils for the requested sectionA : %v \n", seatsResp.UserSeats)

	sectionBResp, err := c.ViewUsersBySection(context.Background(), &pb.ViewUsersRequest{Section: "SectionB"})

	if err != nil {
		log.Fatalf("Could not get the seats deatils for the Section : %v", "SectionB")
	}
	fmt.Printf("Seats Deatils for the requested sectionB : %v \n", sectionBResp.UserSeats)

}
