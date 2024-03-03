// server.go

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "ticket/pb"
)

type ticketServer struct {
	seatAllocations map[pb.SectionTypes]map[int32]pb.User
	pb.UnimplementedTicketServiceServer
}

func newTicketServer() *ticketServer {
	return &ticketServer{
		seatAllocations: make(map[pb.SectionTypes]map[int32]pb.User),
	}
}

func (s *ticketServer) PurchaseTicket(ctx context.Context, ticket *pb.TicketRequest) (*pb.Receipt, error) {

	if s.seatAllocations[pb.SectionTypes_Class_A] == nil {
		s.seatAllocations[pb.SectionTypes_Class_A] = make(map[int32]pb.User)
	}

	// Assign a seat in Class_A
	seatNumber := len(s.seatAllocations[pb.SectionTypes_Class_A]) + 1
	s.seatAllocations[pb.SectionTypes_Class_A][int32(seatNumber)] = *ticket.User

	return &pb.Receipt{
		From:      ticket.From,
		To:        ticket.To,
		User:      ticket.User,
		PricePaid: 20.0,
	}, nil
}

func (s *ticketServer) GetReceipt(ctx context.Context, user *pb.User) (*pb.Receipt, error) {

	for _, seats := range s.seatAllocations {
		for _, u := range seats {
			if u.Email == user.Email {
				return &pb.Receipt{
					From:      "London",
					To:        "France",
					User:      user,
					PricePaid: 20.0,
				}, nil

			}
		}
	}

	return nil, fmt.Errorf("user not found")

}

func (s *ticketServer) GetSeatAllocation(ctx context.Context, section *pb.Section) (*pb.SeatAllocationList, error) {
	// Return the seat allocation for the requested section
	allocations := make([]*pb.SeatAllocation, 0)

	for _, user := range s.seatAllocations[section.GetSection()] {
		allocations = append(allocations, &pb.SeatAllocation{
			User:    &user,
			Section: section,
		})
	}

	return &pb.SeatAllocationList{Allocations: allocations}, nil
}

func (s *ticketServer) RemoveUser(ctx context.Context, user *pb.User) (*pb.RemoveUserResponse, error) {
	// Remove the user from the seat allocation
	for _, seats := range s.seatAllocations {
		for seatNumber, u := range seats {
			if u.Email == user.Email {
				delete(seats, seatNumber)
				return &pb.RemoveUserResponse{Status: "DELETION SUCCEESS"}, nil
			}
		}
	}

	return &pb.RemoveUserResponse{Status: "USER_NOT_FOUND"}, nil
}

func (s *ticketServer) ModifySeat(ctx context.Context, request *pb.ModifySeatRequest) (*pb.SeatAllocation, error) {

	if s.seatAllocations[request.Section.GetSection()] == nil {
		s.seatAllocations[request.Section.GetSection()] = make(map[int32]pb.User)
	}

	// Modify the seat for the given user
	for seatNumber, u := range s.seatAllocations[request.Section.GetSection()] {
		if u.GetEmail() == request.User.GetEmail() {
			updatedUser := u

			// Delete the old entry
			delete(s.seatAllocations[request.Section.GetSection()], seatNumber)
			fmt.Println("333")

			// Add the updated entry with the new seat number
			s.seatAllocations[request.Section.GetSection()][request.Section.GetSeatNumber()] = updatedUser

			return &pb.SeatAllocation{User: &updatedUser, Section: request.Section}, nil
		}
	}

	// Print seat allocations if the user is not found
	fmt.Println("User not found. Current Seat Allocations:", s.seatAllocations)

	return nil, fmt.Errorf("user not found")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterTicketServiceServer(server, newTicketServer())
	fmt.Println("Server listening on port 50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
