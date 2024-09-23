package postgres

import (
	"database/sql"
	"testing"
	pb "user/genproto/group"
	"user/logs"
)

var logger = logs.NewLogger()

func DB()(*sql.DB){
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
		Name: "GO11",
		SubjectId: "44f8f54d-ee40-4e7e-9d5c-8dc5f9b136d9",
		Room: "Osmondagi bolalar",
		StartTime: "14:30",
		EndTime: "20:00",
		StartedAt: "2024-04-22",
	})
	if err != nil{
		t.Fatalf(err.Error())
	}
}

func Test_UpdateGroup(t *testing.T){
	db := DB()
	defer db.Close()

	group := NewGroupRepo(db, logger)

	_, err := group.UpdateGroup(&pb.UpdateGroupReq{
		Id: "405965b2-3717-4f01-b331-693ddb4adf0a",
		Name: "Foundation N54",
		Room: "Yandex",
		StartTime: "17:00",
		EndTime: "20:00",
		StartedAt: "2023-10-10",
	})
	if err != nil{
		t.Fatalf(err.Error())
	}
}
