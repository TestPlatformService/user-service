package postgres

import (
	"database/sql"
	"testing"
	pb "user/genproto/group"
	"user/logs"
)

var logger = logs.NewLogger()

func DB() *sql.DB {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	return db
}

func Test_CreateGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.CreateGroup(&pb.CreateGroupReq{
		Name:      "GO11",
		SubjectId: "44f8f54d-ee40-4e7e-9d5c-8dc5f9b136d9",
		Room:      "Osmondagi bolalar",
		StartTime: "14:30",
		EndTime:   "20:00",
		StartedAt: "2024-04-22",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_UpdateGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.UpdateGroup(&pb.UpdateGroupReq{
		Id:        "405965b2-3717-4f01-b331-693ddb4adf0a",
		Name:      "Foundation N54",
		Room:      "Yandex",
		StartTime: "17:00",
		EndTime:   "20:00",
		StartedAt: "2023-10-10",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_DeleteGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.DeleteGroup(&pb.GroupId{
		Id: "405965b2-3717-4f01-b331-693ddb4adf0a",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetGroupById(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.GetGroupById(&pb.GroupId{
		Id: "de6419ff-85ab-45b1-8d09-9559896159b1",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetAllGroups(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.GetAllGroups(&pb.GetAllGroupsReq{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_AddStudentToGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.AddStudentToGroup(&pb.AddStudentReq{
		GroupId:     "0272cf9c-fd83-4ab7-abe3-a6b1b543cc52",
		StudentHhId: "20389",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_DeleteStudentFromGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.DeleteStudentFromGroup(&pb.DeleteStudentReq{
		GroupId:     "4188773f-6e09-4d4a-9d21-0d593f29116f",
		StudentHhId: "20388",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_AddTeacherToGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.AddTeacherToGroup(&pb.AddTeacherReq{
		GroupId:   "4188773f-6e09-4d4a-9d21-0d593f29116f",
		TeacherId: "9abad94c-ab88-4191-92b8-fa98a2903400",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_DeleteTeacherFromGroup(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.DeleteTeacherFromGroup(&pb.DeleteTeacherReq{
		Id:        "06535a29-4d85-4a11-8b72-978e263ad24a0",
		GroupId:   "4188773f-6e09-4d4a-9d21-0d593f29116f",
		TeacherId: "9abad94c-ab88-4191-92b8-fa98a2903400",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetStudentGroups(t *testing.T) {
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.GetStudentGroups(&pb.StudentId{
		HhId: "20388",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_GetTeacherGroups(t *testing.T){
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.GetTeacherGroups(&pb.TeacherId{
		Id: "9abad94c-ab88-4191-92b8-fa98a2903400",
	})
	if err != nil{
		t.Fatalf(err.Error())
	}
}

func Test_GetGroupStudents(t *testing.T){
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.GetGroupStudents(&pb.GroupId{
		Id: "0272cf9c-fd83-4ab7-abe3-a6b1b543cc52",
	})
	if err != nil{
		t.Fatal(err)
	}
}
