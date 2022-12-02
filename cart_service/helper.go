package main

import (
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) validateUUID(ids ...string) error {
	for _, id := range ids {
		_, err := uuid.FromString(id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
		}
	}

	return nil
}
