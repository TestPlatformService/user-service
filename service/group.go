package service

import (
	"database/sql"
	"fmt"
	"log/slog"
	pb "user/genproto/group"
	"user/storage"
	"user/storage/postgres"
)

type GroupService struct {
	pb.UnimplementedGroupServiceServer
	Storage storage.IStorage
	Logger  *slog.Logger
}

func NewGroupService(db *sql.DB, Logger *slog.Logger) *GroupService {
	return &GroupService{
		Storage: postgres.NewIstorage(db),
		Logger:  Logger,
	}
}

func (g *GroupService) CreateGroup(req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {

	res, err := g.Storage.Group().CreateGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error creating Group: %v", err))
		return nil, err
	}
	g.Logger.Info("CreateGroup rpc method finished")
	return res, nil
}

func (g *GroupService) UpdateGroup(req *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {

	res, err := g.Storage.Group().UpdateGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error updating group: %v", err))
		return nil, err
	}
	g.Logger.Info("UpdateGroup rpc method finished")
	
	return res, nil
}

func (g *GroupService) DeleteGroup(req *pb.GroupId) (*pb.DeleteResp, error) {

	res, err := g.Storage.Group().DeleteGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteGroup rpc method finished")

	return res, nil
}

func (g *GroupService) GetGroupById(req *pb.GroupId) (*pb.Group, error) {

	res, err := g.Storage.Group().GetGroupById(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error GetGroupById: %v", err))
		return nil, err
	}
	g.Logger.Info("GetGroupById rpc method finished")

	return res, nil
}

func (g *GroupService) GetAllGroups(req *pb.GetAllGroupsReq) (*pb.GetAllGroupsResp, error) {

	res, err := g.Storage.Group().GetAllGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error GetAllGroups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetAllGroups rpc method finished")

	return res, nil
}

func (g *GroupService) AddStudentToGroup(req *pb.AddStudentReq) (*pb.AddStudentResp, error) {
	res, err := g.Storage.Group().AddStudentToGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error adding Student to group: %v", err))
		return nil, err
	}
	g.Logger.Info("AddStudentToGroup rpc method finished")

	return res, nil
}

func (g *GroupService) DeleteStudentFromGroup(req *pb.DeleteStudentReq) (*pb.DeleteResp, error) {
	res, err := g.Storage.Group().DeleteStudentFromGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting student from group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteStudentFromGroup rpc method finished")

	return res, nil
}

func (g *GroupService) AddTeacherToGroup(req *pb.AddTeacherReq) (*pb.AddTeacherResp, error) {
	res, err := g.Storage.Group().AddTeacherToGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error adding teacher to group: %v", err))
		return nil, err
	}
	g.Logger.Info("AddTeacherToGroup rpc method finished")

	return res, nil
}

func (g *GroupService) DeleteTeacherFromGroup(req *pb.DeleteTeacherReq) (*pb.DeleteResp, error) {
	res, err := g.Storage.Group().DeleteTeacherFromGroup(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error deleting group: %v", err))
		return nil, err
	}
	g.Logger.Info("DeleteTeacherFromGroup rpc method finished")

	return res, nil
}

func (g *GroupService) GetStudentGroups(req *pb.StudentId) (*pb.StudentGroups, error) {
	res, err := g.Storage.Group().GetStudentGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error getting Student groups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetStudentGroups rpc method finished")

	return res, nil
}

func (g *GroupService) GetTeacherGroups(req *pb.TeacherId) (*pb.TeacherGroups, error) {
	res, err := g.Storage.Group().GetTeacherGroups(req)
	if err != nil {
		g.Logger.Error(fmt.Sprintf("Error getting Teacher groups: %v", err))
		return nil, err
	}
	g.Logger.Info("GetTeacherGroups rpc method finished")

	return res, nil
}