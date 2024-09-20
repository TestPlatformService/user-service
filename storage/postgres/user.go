package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	pb "user/genproto/user"
	"user/logs"
	"user/storage"
)



type UserRepo struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewUserRepo(DB *sql.DB) storage.IUserStorage {
	return &UserRepo{DB: DB, Log: logs.NewLogger()}
}

func (u *UserRepo) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Void, error) {

	query := `INSERT INTO users (hh_id, first_name, last_name, password_hash, phone_number, gender, date_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	err := u.DB.QueryRow(query, req.HhId, req.Firstname, req.Lastname, req.Password, req.Phone, req.Gender, req.DateOfBirth)
	if err != nil {
		u.Log.Error("Error while Registering", "error", err)
		return nil, errors.ErrUnsupported
	}

	return &pb.Void{}, nil
}

func (u *UserRepo) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	var user pb.LoginResponse

	query := `SELECT id, role FROM users WHERE hh_id = $1 and password_hash = $2 and deleted_at IS NULL`

	err := u.DB.QueryRow(query, req.HhId, req.Password).Scan(&user.HhId, &user.Role)
	if err == sql.ErrNoRows {
		u.Log.Error("No user found", "ID", req.HhId)
		return nil, errors.New("no user found")
	} else if err != nil {
		u.Log.Error("Error getting user by ID", "err", err)
		return nil, err
	}

	return &pb.LoginResponse{
		HhId: user.HhId,
		Role: user.Role,
	}, nil
}

func (u *UserRepo) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var user pb.GetProfileResponse

	query := `SELECT hh_id, first_name, last_name, password_hash, phone_number, date_of_birth, gender`

	err := u.DB.QueryRow(query, )
}