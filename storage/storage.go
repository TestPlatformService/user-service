package storage

import (
	"context"
	pb1 "user/genproto/notification"
	pb "user/genproto/user"
)

type IStorage interface {
	User() IUserStorage
	Notifications() INotificationStorage
	Close()
}

type IUserStorage interface {
	Register(context.Context, *pb.RegisterRequest) (*pb.Void, error)
	Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)
	GetProfile(context.Context, *pb.GetProfileRequest) (*pb.GetProfileResponse, error)
	GetAllUsers(context.Context, *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error)
	UpdateProfile(context.Context, *pb.UpdateProfileRequest) (*pb.Void, error)
	UpdateProfileAdmin(context.Context, *pb.UpdateProfileAdminRequest) (*pb.Void, error)
	DeleteProfile(context.Context, *pb.DeleteProfileRequest) (*pb.Void, error)
}

type INotificationStorage interface {
	CreateNotifications(context.Context, *pb1.CreateNotificationsReq) (*pb1.CreateNotificationsRes, error)
	GetAllNotifications(context.Context, *pb1.GetNotificationsReq) (*pb1.GetNotificationsResponse, error)
	GetAndMarkNotificationAsRead(context.Context, *pb1.GetAndMarkNotificationAsReadReq) (*pb1.GetAndMarkNotificationAsReadRes, error)
}