package service

import (
	"context"

	"grpcexi/internal/repo"
	pb "grpcexi/protos"

	"google.golang.org/grpc"
)

type ContactService struct {
	pb.UnimplementedContactServiceServer
	repo repo.ContactRepo
}

func NewContactService(repo repo.ContactRepo) *ContactService {
	return &ContactService{
		repo: repo,
	}
}

func (s *ContactService) Register(server *grpc.Server) {
	pb.RegisterContactServiceServer(server, s)
}

func (s *ContactService) GetContacts(all *pb.All, stream pb.ContactService_GetContactsServer) error {
	contacts, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	for _, c := range contacts {
		contactMsg := &pb.Contacts{Contacts: []*pb.Contact{c}}
		if err := stream.Send(contactMsg); err != nil {
			return err
		}
	}
	return nil
}

func (s *ContactService) AddContact(ctx context.Context, contact *pb.Contact) (*pb.Contact, error) {
	s.repo.Create(contact)
	return contact, nil
}
