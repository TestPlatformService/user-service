package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	pb "user/genproto/group"
	"user/storage"
)

type GroupService struct {
	pb.UnimplementedGroupServiceServer
	Storage storage.IStorage
	Logger  *slog.Logger
}

func NewGroupService(db *sql.DB, Logger *slog.Logger, istorage storage.IStorage) *GroupService {
	return &GroupService{
		Storage: istorage,
		Logger:  Logger,
	}
}

func (g *GroupService) CreateGroup(ctx context.Context, req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {

	res, err := g.Storage.Group().CreateGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error creating Group: %v", err))
		return nil, err
	}
	g.Logger.Info("CreateGroup rpc method finished")
	return res, nil
}

func (g *GroupService) UpdateGroup(ctx context.Context, req *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {

	res, err := g.Storage.Group().UpdateGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error updating group: %v", err))
		return nil, err
	}
	g.Logger.Info("UpdateGroup rpc method finished")

	return res, nil
}

func (g *GroupService) DeleteGroup(ctx context.Context, req *pb.GroupId) (*pb.DeleteResp, error) {

	res, err := g.Storage.Group().DeleteGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteGroup rpc method finished")

	return res, nil
}

func (g *GroupService) GetGroupById(ctx context.Context, req *pb.GroupId) (*pb.Group, error) {

	res, err := g.Storage.Group().GetGroupById(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error GetGroupById: %v", err))
		return nil, err
	}
	g.Logger.Info("GetGroupById rpc method finished")

	return res, nil
}

func (g *GroupService) GetAllGroups(ctx context.Context, req *pb.GetAllGroupsReq) (*pb.GetAllGroupsResp, error) {

	res, err := g.Storage.Group().GetAllGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error GetAllGroups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetAllGroups rpc method finished")

	return res, nil
}

func (g *GroupService) AddStudentToGroup(ctx context.Context, req *pb.AddStudentReq) (*pb.AddStudentResp, error) {
	res, err := g.Storage.Group().AddStudentToGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error adding Student to group: %v", err))
		return nil, err
	}
	g.Logger.Info("AddStudentToGroup rpc method finished")

	return res, nil
}

func (g *GroupService) DeleteStudentFromGroup(ctx context.Context, req *pb.DeleteStudentReq) (*pb.DeleteResp, error) {
	res, err := g.Storage.Group().DeleteStudentFromGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting student from group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteStudentFromGroup rpc method finished")

	return res, nil
}

func (g *GroupService) AddTeacherToGroup(ctx context.Context, req *pb.AddTeacherReq) (*pb.AddTeacherResp, error) {
	res, err := g.Storage.Group().AddTeacherToGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error adding teacher to group: %v", err))
		return nil, err
	}
	g.Logger.Info("AddTeacherToGroup rpc method finished")

	return res, nil
}

func (g *GroupService) DeleteTeacherFromGroup(ctx context.Context, req *pb.DeleteTeacherReq) (*pb.DeleteResp, error) {
	res, err := g.Storage.Group().DeleteTeacherFromGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteTeacherFromGroup rpc method finished")

	return res, nil
}

func (g *GroupService) GetStudentGroups(ctx context.Context, req *pb.StudentId) (*pb.StudentGroups, error) {
	res, err := g.Storage.Group().GetStudentGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error getting Student groups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetStudentGroups rpc method finished")

	return res, nil
}

func (g *GroupService) GetTeacherGroups(ctx context.Context, req *pb.TeacherId) (*pb.TeacherGroups, error) {
	res, err := g.Storage.Group().GetTeacherGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error getting Teacher groups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetTeacherGroups rpc method finished")

	return res, nil
}

func (g *GroupService) GetGroupStudents(ctx context.Context, req *pb.GroupId) (*pb.GroupStudents, error){
	resp, err := g.Storage.Group().GetGroupStudents(req)
	if err != nil{
		g.Logger.Error(fmt.Sprintf("Error getting group students: %v", err))
		return nil, err
	}
	g.Logger.Info("GetGroupStudents rpc method finished")
	return resp, nil
}

func (g *GroupService) CreateGroupDay(ctx context.Context, req *pb.CreateGroupDayReq) (*pb.CreateGroupDayResp, error){
	resp, err := g.Storage.Group().CreateGroupDay(req)
	if err != nil{
		g.Logger.Error(fmt.Sprintf("Error create group day: %v", err))
		return nil, err
	}
	g.Logger.Info("CreateGroupDay rpc method finished")
	return resp, nil
}

func (g *GroupService) DeleteGroupDay(ctx context.Context, req *pb.DeleteGroupDayReq) (*pb.DeleteGroupDayResp, error){
	resp, err := g.Storage.Group().DeleteGroupDay(req)
	if err != nil{
		g.Logger.Error(fmt.Sprintf("Error delete group day: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteGroupDay rpc method finishid")
	return resp, nil
}