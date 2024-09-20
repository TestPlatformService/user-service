package storage

import (
	"context"
	pb "user/genproto/user"
)

type IStorage interface {
	User() IUserStorage
	Notifications() INotificationStorage
	Close()
}

type IUserStorage interface {
	Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)
	GetProfile(context.Context, *pb.GetProfileRequest) (*pb.GetProfileResponse, error)
	GetAllUsers(context.Context, *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error)
	UpdateProfile(context.Context, *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error)
	UpdateProfileAdmin(context.Context, *pb.UpdateProfileAdminRequest) (*pb.UpdateProfileAdminResponse, error)
	DeleteProfile(context.Context, *pb.DeleteProfileRequest) (*pb.DeleteProfileResponse, error)
}

type INotificationStorage interface {
	CreateNotifications(context.Context, *pb.CreateNotificationsReq) (*pb.CreateNotificationsRes, error)
	GetAllNotifications(context.Context, *pb.GetNotificationsReq) (*pb.GetNotificationsResponse, error)
	GetAndMarkNotificationAsRead(context.Context, *pb.GetAndMarkNotificationAsReadReq) (*pb.GetAndMarkNotificationAsReadRes, error)
}
