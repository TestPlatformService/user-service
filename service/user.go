package service

import (
	"context"
	"database/sql"
	"log/slog"
	pb "user/genproto/user"
	"user/storage"
	"user/storage/postgres"
)

type UserService struct {
	pb.UnimplementedUsersServer
	User storage.IStorage

	Logger *slog.Logger
}

func NewUserService(db *sql.DB, Logger *slog.Logger) *UserService {
	return &UserService{
		User:   postgres.NewUserRepo(db),
		Logger: Logger,
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Void, error) {
	_, err := s.User.User().Register(ctx, req)
	if err != nil {
		s.Logger.Error("failed to create user", "error", err)
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *UserService) 