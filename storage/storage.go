package storage

import (
	"context"
	pb2 "user/genproto/group"
	pb1 "user/genproto/notification"
	pb "user/genproto/user"
)

type IStorage interface {
	User() IUserStorage
	Notifications() INotificationStorage
	Group() IGroupStorage
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

type IGroupStorage interface {
	CreateGroup(*pb2.CreateGroupReq) (*pb2.CreateGroupResp, error)
	UpdateGroup(*pb2.UpdateGroupReq) (*pb2.UpdateGroupResp, error)
	DeleteGroup(*pb2.GroupId) (*pb2.DeleteResp, error)
	GetGroupById(*pb2.GroupId) (*pb2.Group, error)
	GetAllGroups(*pb2.GetAllGroupsReq) (*pb2.GetAllGroupsResp, error)
	AddStudentToGroup(*pb2.AddStudentReq) (*pb2.AddStudentResp, error)
	DeleteStudentFromGroup(*pb2.DeleteStudentReq) (*pb2.DeleteResp, error)
	AddTeacherToGroup(*pb2.AddTeacherReq) (*pb2.AddTeacherResp, error)
	DeleteTeacherFromGroup(*pb2.DeleteTeacherReq) (*pb2.DeleteResp, error)
	GetStudentGroups(*pb2.StudentId) (*pb2.StudentGroups, error)
	GetTeacherGroups(*pb2.TeacherId) (*pb2.TeacherGroups, error)
	GetGroupStudents(req *pb2.GroupId)(*pb2.GroupStudents, error)
}
