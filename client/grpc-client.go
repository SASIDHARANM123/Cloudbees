// client.go

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "ticket/pb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTicketServiceClient(conn)

	user := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	// Purchase Ticket
	ticket := &pb.TicketRequest{
		From: "London",
		To:   "France",
		User: user,
	}
	PurchaseTicket, err := client.PurchaseTicket(context.Background(), ticket)
	if err != nil {
		log.Fatalf("Failed to purchase ticket: %v", err)
	}
	fmt.Println("PurchaseTicket: ", PurchaseTicket)

	// Get Receipt
	GetReceipt, err := client.GetReceipt(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to get receipt: %v", err)
	}
	fmt.Println("GetReceipt: ", GetReceipt)

	// Get Seat Allocation
	reqSection := &pb.Section{Section: pb.SectionTypes_Class_A}
	seatAllocation, err := client.GetSeatAllocation(context.Background(), reqSection)
	if err != nil {
		log.Fatalf("Failed to get seat allocation: %v", err)
	}
	fmt.Printf("Seat Allocation :\n%v\n", seatAllocation.Allocations)

	// Modify Seat
	modifyRequest := &pb.ModifySeatRequest{
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com", // Provide the user's email
		},
		Section: &pb.Section{
			Section:    pb.SectionTypes_Class_A, // Corrected enum value
			SeatNumber: 42,                      // Provide the new seat number
		},
	}
	modifyResponse, err := client.ModifySeat(context.Background(), modifyRequest)
	if err != nil {
		log.Fatalf("Failed to modify seat: %v", err)
	}
	fmt.Printf("Seat modified: %v\n", modifyResponse)

	// Remove User
	removeResponse, err := client.RemoveUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to remove user: %v", err)
	}
	fmt.Printf("User removed: %v\n", removeResponse)
}
