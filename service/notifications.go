package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	pb "user/genproto/notification"
	"user/storage"
)

type NotificationsService struct {
	pb.UnimplementedNotificationsServer
	Storage storage.IStorage
	Logger  *slog.Logger
}

func NewNotificationsService(db *sql.DB, Logger *slog.Logger, istorage storage.IStorage) *NotificationsService {
	return &NotificationsService{
		Storage: istorage,
		Logger:  Logger,
	}
}

func (s *NotificationsService) CreateNotification(ctx context.Context, req *pb.CreateNotificationsReq) (*pb.CreateNotificationsRes, error) {
	s.Logger.Info("CreateNotifications rpc method is working")
	resp, err := s.Storage.Notifications().CreateNotifications(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error creating notification: %v", err))
		return nil, err
	}
	s.Logger.Info("CreateNotifications rpc method finished")
	return resp, nil
}

func (s *NotificationsService) GetAllNotifications(ctx context.Context, req *pb.GetNotificationsReq) (*pb.GetNotificationsResponse, error) {
	s.Logger.Info("GetAllNotifications rpc method is working")
	resp, err := s.Storage.Notifications().GetAllNotifications(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error getting notifications: %v", err))
		return nil, err
	}
	s.Logger.Info("GetAllNotifications rpc method finished")
	return resp, nil
}

func (s *NotificationsService) GetAndMarkNotificationAsRead(ctx context.Context, req *pb.GetAndMarkNotificationAsReadReq) (*pb.GetAndMarkNotificationAsReadRes, error) {
	s.Logger.Info("GetAndMarkNotificationAsRead rpc method is working")
	resp, err := s.Storage.Notifications().GetAndMarkNotificationAsRead(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error getting and marking notification as read: %v", err))
		return nil, err
	}
	s.Logger.Info("GetAndMarkNotificationAsRead rpc method finished")
	return resp, nil
}

func (s *NotificationsService) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadReq) (*pb.Void, error) {
	s.Logger.Info("MarkNotificationAsRead rpc method is working")
	_, err := s.Storage.Notifications().MarkNotificationAsRead(ctx, req)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error marking notification as read: %v", err))
		return nil, err
	}
	s.Logger.Info("MarkNotificationAsRead rpc method finished")
	return &pb.Void{}, nil
}
