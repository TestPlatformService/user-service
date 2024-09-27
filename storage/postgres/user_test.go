package postgres

import (
	"context"
	"testing"

	pb "user/genproto/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Register(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	user := &pb.RegisterRequest{
		HhId:        "20389",
		Firstname:   "Sanjarbek",
		Lastname:    "Abduraxmonov",
		Password:    "1111",
		Phone:       "+998940375108",
		DateOfBirth: "2007/05/16",
		Gender:      "male",
	}

	_, err = repo.Register(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

}

func TestUserRepo_Login(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	login := &pb.LoginRequest{
		HhId:     "20388",
		Password: "1111",
	}

	res, err := repo.Login(context.Background(), login)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, res)
}

func TestUserRepo_GetProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	profile := &pb.GetProfileRequest{
		Id: "1ffc468f-bac4-4935-aa7f-0159cc38e22f",
	}

	res, err := repo.GetProfile(context.Background(), profile)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, res)
}

func TestUserRepo_UpdateProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	req := &pb.UpdateProfileRequest{
		Id:             "1ffc468f-bac4-4935-aa7f-0159cc38e22f",
		Password:       "1111",
	}

	_, err = repo.UpdateProfile(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_UpdateProfileAdmin(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepo(db)
	req := &pb.UpdateProfileAdminRequest{
		Id:          "1ffc468f-bac4-4935-aa7f-0159cc38e22f",
		Firstname:   "Sanjarbek",
		Lastname:    "Abduraxmonov",
		Password:    "1111",
		Phone:       "+998940375107",
		DateOfBirth: "2007/05/16",
		Gender:      "male",
	}

	_, err = repo.UpdateProfileAdmin(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

}

func TestUserRepo_DeleteProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec("UPDATE users SET deleted_at").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.DeleteProfileRequest{
		Id: "1",
	}

	_, err = repo.DeleteProfile(context.Background(), req)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
