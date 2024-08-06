package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Aravinthvvs/gRPC/proto/train/train"

	"google.golang.org/grpc"
)

func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [options]", os.Args[0])
	}

	command := os.Args[1]

	// Establish a connection to the server
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTicketServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	switch command {
	case "purchase":
		if len(os.Args) < 6 {
			log.Fatalf("Usage: %s purchase <from> <to> <first_name> <last_name> <email>", os.Args[0])
		}
		from := os.Args[2]
		to := os.Args[3]
		firstName := os.Args[4]
		lastName := os.Args[5]
		email := os.Args[6]

		user := &pb.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}
		resp, err := c.PurchaseTicket(ctx, &pb.PurchaseRequest{From: from, To: to, User: user})
		if err != nil {
			log.Fatalf("could not purchase ticket: %v", err)
		}
		fmt.Printf("Purchase Response: %s\n", resp.ReceiptId)

	case "get_receipt":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s get_receipt <receipt_id>", os.Args[0])
		}
		receiptId := os.Args[2]
		resp, err := c.GetReceipt(ctx, &pb.ReceiptRequest{ReceiptId: receiptId})
		if err != nil {
			log.Fatalf("could not get receipt: %v", err)
		}
		fmt.Printf("Receipt: %+v\n", resp)

	case "view_users":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s view_users <section>", os.Args[0])
		}
		section := os.Args[2]
		resp, err := c.ViewUsersBySection(ctx, &pb.ViewUsersRequest{Section: section})
		if err != nil {
			log.Fatalf("could not view users: %v", err)
		}
		fmt.Printf("Users in %s: %+v\n", section, resp.UserSeats)

	case "remove_user":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s remove_user <email>", os.Args[0])
		}
		email := os.Args[2]
		resp, err := c.RemoveUser(ctx, &pb.RemoveUserRequest{Email: email})
		if err != nil {
			log.Fatalf("could not remove user: %v", err)
		}
		if resp.Success {
			fmt.Println("User removed successfully.")
		} else {
			fmt.Println("Failed to remove user.")
		}

	case "modify_seat":
		if len(os.Args) < 4 {
			log.Fatalf("Usage: %s modify_seat <email> <new_seat>", os.Args[0])
		}
		email := os.Args[2]
		newSeat := os.Args[3]
		resp, err := c.ModifySeat(ctx, &pb.ModifySeatRequest{Email: email, NewSeat: newSeat})
		if err != nil {
			log.Fatalf("could not modify seat: %v", err)
		}
		if resp.Success {
			fmt.Println("Seat modified successfully.")
		} else {
			fmt.Println("Failed to modify seat.")
		}

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
