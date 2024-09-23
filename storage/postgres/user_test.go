package postgres

import (
	"context"
	"testing"
	"time"

	pb "user/genproto/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec("INSERT INTO users").WithArgs("hh123", "John", "Doe", "password", "1234567890", "male", "1990-01-01").WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.RegisterRequest{
		HhId:        "hh123",
		Firstname:   "John",
		Lastname:    "Doe",
		Password:    "password",
		Phone:       "1234567890",
		Gender:      "male",
		DateOfBirth: "1990-01-01",
	}

	_, err = repo.Register(context.Background(), req)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepo_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	rows := sqlmock.NewRows([]string{"id", "role"}).AddRow(1, "user")
	mock.ExpectQuery("SELECT id, role FROM users").WithArgs("hh123", "password").WillReturnRows(rows)

	req := &pb.LoginRequest{
		HhId:     "hh123",
		Password: "password",
	}

	resp, err := repo.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "user", resp.Role)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepo_GetProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"hh_id", "first_name", "last_name", "password_hash", "phone_number", "date_of_birth", "gender"}).
		AddRow("hh123", "John", "Doe", "password", "1234567890", dob, "male")
	mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(1).WillReturnRows(rows)

	req := &pb.GetProfileRequest{
		Id: "1",
	}

	resp, err := repo.GetProfile(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "hh123", resp.HhId)
	assert.Equal(t, "John", resp.Firstname)
	assert.Equal(t, "Doe", resp.Lastname)
	assert.Equal(t, "password", resp.Password)
	assert.Equal(t, "1234567890", resp.Phone)
	assert.Equal(t, dob.Format("2006-01-02"), resp.DateOfBirth)
	assert.Equal(t, "male", resp.Gender)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepo_UpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec("UPDATE users SET").WithArgs("profile.jpg", "newpassword", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.UpdateProfileRequest{
		Id:             "1",
		ProfilePicture: "profile.jpg",
		Password:       "newpassword",
	}

	_, err = repo.UpdateProfile(context.Background(), req)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepo_UpdateProfileAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock yaratishda xato: %v", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	mock.ExpectExec("UPDATE users SET").WithArgs("John", "Doe", "newpassword", "1234567890", "1990-01-01", "male", "A1", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.UpdateProfileAdminRequest{
		Id:          "1",
		Firstname:   "John",
		Lastname:    "Doe",
		Password:    "newpassword",
		Phone:       "1234567890",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		Group:       "A1",
	}

	_, err = repo.UpdateProfileAdmin(context.Background(), req)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
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
