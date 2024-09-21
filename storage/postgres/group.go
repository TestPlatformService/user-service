package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "user/genproto/group"

	"github.com/google/uuid"
)

type GroupRepo interface {
}

type groupImpl struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewGroupRepo(db *sql.DB, logger *slog.Logger) GroupRepo {
	return &groupImpl{
		DB:     db,
		Logger: logger,
	}
}

func (G *groupImpl) CreateGroup(req *pb.CreateGroupReq) (*pb.CreateGroupResp, error) {
	id := uuid.NewString()
	query := `
        INSERT INTO groups(
          id, name, subject_id, room, start_time, end_time, started_at)
        VALUES
          ($1, $2, $3, $4, $5, $6, $7)`
	_, err := G.DB.Exec(query, id, req.Name, req.SubjectId, req.Room, req.StartTime, req.EndTime, req.StartedAt)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}

	return &pb.CreateGroupResp{
		Id:        id,
		CreatedAt: time.Now().String(),
	}, nil
}

func (G *groupImpl) UpdateGroup(req *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	query := `
        UPDATE groups SET 
          name = $1, room = $2, start_time = $3, end_time = $4, started_at = $5
        WHERE 
          id = $6 AND deleted_at IS NULL`
	_, err := G.DB.Exec(query, req.Name, req.Room, req.StartTime, req.EndTime, req.StartTime, req.Id)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	return &pb.UpdateGroupResp{
		Id:        req.Id,
		UpdatedAt: time.November.String(),
	}, nil
}

func (G *groupImpl) DeleteGroup(req *pb.GroupId) (*pb.DeleteResp, error) {
	query := `
        UPDATE groups SET
          deleted_at = $1
        WHERE 
          id = $2 AND deleted_at IS NULL`
	_, err := G.DB.Exec(query, time.Now().String(), req.Id)
	if err != nil {
		G.Logger.Error(err.Error())
		return &pb.DeleteResp{
			Status: "Group o'chirilmadi",
		}, err
	}
	return &pb.DeleteResp{
		Status: "Group o'chirildi",
	}, nil
}

func (G *groupImpl) GetGroupById(req *pb.GroupId) (*pb.Group, error) {
	resp := pb.Group{Id: req.Id}
	query := `
        SELECT 
          name, subject_id, room, start_time, end_time, started_at
        FROM 
          groups
        WHERE
          id = $1 AND deleted_at IS NULL`
	err := G.DB.QueryRow(query, req.Id).Scan(&resp.Name, &resp.SubjectId, &resp.Room, &resp.StartTime, &resp.EndTime, &resp.StartedAt)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	query = `
        SELECT 
          teacher_id
        FROM 
          teacher_groups
        WHERE
          group_id = $1 AND deleted_at IS NULL`
	err = G.DB.QueryRow(query, req.Id).Scan(&resp.TeacherId)
	if err != nil {
		G.Logger.Error(err.Error())
	}
	return &resp, nil
}

func (G *groupImpl) GetAllGroups(req *pb.GetAllGroupsReq) (*pb.GetAllGroupsResp, error) {
	groups := []*pb.Group{}
	query := `
        SELECT 
          id, name, subject_id, room, start_time, end_time, started_at
        FROM 
          groups
        WHERE 
          deleted_at IS NULL`
	if len(req.SubjectId) > 0 {
		query += fmt.Sprintf(" AND subject_id = %s", req.SubjectId)
	}
	if len(req.Room) > 0 {
		query += fmt.Sprintf(" AND room = %s", req.Room)
	}
	query += fmt.Sprintf(" limit = %d offset = %d", req.Limit, req.Offset)

	rows, err := G.DB.Query(query)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		var group *pb.Group
		err = rows.Scan(&group.Id, &group.Name, &group.SubjectId, &group.Room, &group.StartTime, &group.EndTime, &group.StartedAt)
		if err != nil {
			G.Logger.Error(err.Error())
			return nil, err
		}
		groups = append(groups, group)
	}
	return &pb.GetAllGroupsResp{
		Groups: groups,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

func (G *groupImpl) AddStudentToGroup(req *pb.AddStudentReq) (*pb.AddStudentResp, error) {
	id := uuid.NewString()
	query := `
		  INSERT INTO student_groups(
			id, student_hh_id, group_id)
		  VALUES
			($1, $2, $3)`
	_, err := G.DB.Exec(query, id, req.StudentId, req.GroupId)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	return &pb.AddStudentResp{
		Id:        id,
		CreatedAt: time.Now().String(),
	}, nil
}

func (G *groupImpl) DeleteStudentFromGroup(req *pb.DeleteStudentReq) (*pb.DeleteResp, error) {
	query := `
		  UPDATE user_groups SET 
			deleted_at = $1
		  WHERE 
			group_id = $2 AND student_hh_id = $3 AND deleted_at IS NULL`
	_, err := G.DB.Exec(query, time.Now().String(), req.GroupId, req.UserId)
	if err != nil {
		G.Logger.Error(err.Error())
		return &pb.DeleteResp{
			Status: "User guruhdan o'chirilmadi",
		}, err
	}
	return &pb.DeleteResp{
		Status: "User guruhdan muvaffaqiyatli o'chirildi",
	}, nil
}

func (G *groupImpl) AddTeacherToGroup(req *pb.AddTeacherReq) (*pb.AddTeacherResp, error) {
	id := uuid.NewString()
	query := `
		  INSERT INTO teacher_groups(
			id, teacher_id, group_id)
		  VALUES
			($1, $2, $3)`
	_, err := G.DB.Exec(query, id, req.TeacherId, req.GroupId)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	return &pb.AddTeacherResp{
		Id:        id,
		CreatedAt: time.Now().String(),
	}, nil
}

func (G *groupImpl) DeleteTeacherFromGroup(req *pb.DeleteTeacherReq) (*pb.DeleteResp, error) {
	query := `
		  UPDATE teacher_groups SET 
			deleted_at = $1
		  WHERE 
			teacher_id = $2 AND group_id = $3 AND deleted_at IS NULL`
	_, err := G.DB.Exec(query, time.Now().String(), req.TeacherId, req.GroupId)
	if err != nil {
		G.Logger.Error(err.Error())
		return &pb.DeleteResp{
			Status: "Teacher guruhdan o'chirilmadi",
		}, err
	}
	return &pb.DeleteResp{
		Status: "Teacher guruhdan muvaffaqiyatli o'chirildi",
	}, nil
}

func (G *groupImpl) GetStudentGroups(req *pb.StudentId) (*pb.StudentGroups, error) {
	var groups []*pb.Group
	query := `
		  SELECT 
			G.id, G.name, G.subject_id, G.teacher_id, G.room, G.start_time, G.end_time, G.started_at
		  FROM 
			groups AS G
		  JOIN 
			student_groups AS SG
		  ON 
			SG.group_id = G.id
		  WHERE
			SG.student_id = $1 AND SG.deleted_at IS NULL`
	rows, err := G.DB.Query(query, req.Id)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		var group *pb.Group
		err = rows.Scan(&group.Id, &group.Name, &group.SubjectId, &group.TeacherId, &group.Room, &group.StartTime, &group.EndTime, &group.StartedAt)
		if err != nil {
			G.Logger.Error(err.Error())
			return nil, err
		}
		groups = append(groups, group)
	}
	return &pb.StudentGroups{
		Groups: groups,
	}, nil
}

func (G *groupImpl) GetTeacherGroups(req *pb.TeacherId) (*pb.TeacherGroups, error) {
	var groups []*pb.Group
	query := `
		  SELECT 
			G.id, G.name, G.subject_id, G.teacher_id, G.room, G.start_time, G.end_time, G.started_at
		  FROM 
			groups AS G
		  JOIN 
			teacher_groups AS TG
		  ON 
			TG.group_id = G.id
		  WHERE
			TG.teacher_id = $1 AND SG.deleted_at IS NULL`
	rows, err := G.DB.Query(query, req.Id)
	if err != nil {
		G.Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		var group *pb.Group
		err = rows.Scan(&group.Id, &group.Name, &group.SubjectId, &group.TeacherId, &group.Room, &group.StartTime, &group.EndTime, &group.StartedAt)
		if err != nil {
			G.Logger.Error(err.Error())
			return nil, err
		}
		groups = append(groups, group)
	}
	return &pb.TeacherGroups{
		Groups: groups,
	}, nil
}