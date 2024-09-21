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
		User:   postgres.NewIstorage(db),
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

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := s.User.User().Login(ctx, req)
	if err != nil {
		s.Logger.Error("failed to login", "error", err)
		return nil, err
	}

	return &pb.LoginResponse{
		Id:   res.Id,
		Role: res.Role,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	res, err := s.User.User().GetProfile(ctx, req)
	if err != nil {
		s.Logger.Error("Failed to get profile", "error", err)
		return nil, err
	}

	return &pb.GetProfileResponse{
		HhId:        res.HhId,
		Firstname:   res.Firstname,
		Lastname:    res.Lastname,
		Password:    res.Password,
		Phone:       res.Phone,
		DateOfBirth: res.DateOfBirth,
		Gender:      res.Gender,
		Id:          res.Id,
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	res, err := s.User.User().GetAllUsers(ctx, req)
	if err != nil {
		s.Logger.Error("Failed to GetAllUsers", "error", err)
		return nil, err
	}

	return &pb.GetAllUsersResponse{
		Users:      res.Users,
		TotalCount: res.TotalCount,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Void, error) {
	_, err := s.User.User().UpdateProfile(ctx, req)
	if err != nil {
		s.Logger.Error("Failed to update profile", "error", err)
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *UserService) UpdateProfileAdmin(ctx context.Context, req *pb.UpdateProfileAdminRequest) (*pb.Void, error) {
	_, err := s.User.User().UpdateProfileAdmin(ctx, req)
	if err != nil {
		s.Logger.Error("Failed to UpdateProfileAdmin", "error", err)
		return nil, err
	}

	return &pb.Void{}, nil
}

func (s *UserService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.Void, error) {
	_, err := s.User.User().DeleteProfile(ctx, req)
	if err != nil {
		s.Logger.Error("Failed to Delete profile", "error", err)
		return nil, err
	}

	return &pb.Void{}, nil
}
