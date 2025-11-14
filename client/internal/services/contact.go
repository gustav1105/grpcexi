package services

import (
	"context"
	"time"

	pb "grpcexi/protos"
)

type ContactService struct {
	client pb.ContactServiceClient
}

func NewContactService(client pb.ContactServiceClient) *ContactService {
	return &ContactService{client: client}
}

func (s *ContactService) AddContact(ctx context.Context, contact *pb.Contact) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := s.client.AddContact(ctx, contact)
	return err
}

func (s *ContactService) StreamContacts(ctx context.Context) (<-chan *pb.Contact, error) {
	ch := make(chan *pb.Contact)
	stream, err := s.client.GetContacts(ctx, &pb.All{})
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err != nil {
				//error
				return
			}
			for _, c := range resp.Contacts {
				ch <- c
			}
		}
	}()

	return ch, nil
}
