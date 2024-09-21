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

	_, err := u.DB.Exec(query, req.HhId, req.Firstname, req.Lastname, req.Password, req.Phone, req.Gender, req.DateOfBirth)

	if err != nil {
		u.Log.Error("Error while Registering", "error", err)
		return nil, errors.ErrUnsupported
	}

	return &pb.Void{}, nil
}

func (u *UserRepo) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	var user pb.LoginResponse

	query := `SELECT id, role 
	FROM users 
	WHERE hh_id = $1 and password_hash = $2 and deleted_at IS NULL`

	err := u.DB.QueryRow(query, req.HhId, req.Password).Scan(&user.Id, &user.Role)
	if err == sql.ErrNoRows {
		u.Log.Error("No user found", "ID", req.HhId)
		return nil, errors.New("no user found")
	} else if err != nil {
		u.Log.Error("Error getting user by ID", "err", err)
		return nil, err
	}

	return &pb.LoginResponse{
		Id: user.Id,
		Role: user.Role,
	}, nil
}

func (u *UserRepo) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var user pb.GetProfileResponse

	query := `SELECT hh_id, first_name, last_name, password_hash, phone_number, date_of_birth, gender 
	FROM users 
	WHERE id = $1 and deleted_at IS NULL`

	err := u.DB.QueryRow(query, req.Id).Scan(&user.HhId, &user.Firstname, &user.Lastname, &user.Password, &user.Phone, &user.DateOfBirth, &user.Gender)
	if err == sql.ErrNoRows {
		u.Log.Error("No user found", "ID", req.Id)
		return nil, errors.New("no user found")
	} else if err != nil {
		u.Log.Error("Error getting user by ID", "err", err)
		return nil, err
	}

	return &pb.GetProfileResponse{
		HhId: user.HhId,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Password: user.Password,
		Phone: user.Phone,
		DateOfBirth: user.DateOfBirth,
		Gender: user.Gender,
		Id: user.Id,
	}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
    var allUsers pb.GetAllUsersResponse

    countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
    var totalCount int64
    err := u.DB.QueryRow(countQuery).Scan(&totalCount)
    if err != nil {
        u.Log.Error("Error while counting users", "error", err)
        return nil, err
    }
    allUsers.TotalCount = totalCount

    offset := (req.Page - 1) * req.Limit
    if req.Page <= 0 {
        offset = 0
    }

    query := `
        SELECT role, "group", subject, teacher, hh_id, phone_number, gender 
        FROM users 
        WHERE deleted_at IS NULL 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    `

    rows, err := u.DB.Query(query, req.Limit, offset)
    if err != nil {
        u.Log.Error("Error fetching users", "error", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user pb.AllUsers
        err := rows.Scan(&user.Role, &user.Group, &user.Subject, &user.Teacher, &user.HhId, &user.PhoneNumber, &user.Gender)
        if err != nil {
            u.Log.Error("Error scanning user", "error", err)
            return nil, err
        }
        allUsers.Users = append(allUsers.Users, &user)
    }

    if err := rows.Err(); err != nil {
        u.Log.Error("Row iteration error", "error", err)
        return nil, err
    }

    return &allUsers, nil
}


func (u *UserRepo) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Void, error) {
	query := `UPDATE users SET profile_image = $1, password_hash = $2, updated_at = CURRENT_TIMESTAMP 
	WHERE id = $3 AND deleted_at IS NULL`

	_, err := u.DB.Exec(query, req.ProfilePicture, req.Password, req.Id)
	if err != nil {
		u.Log.Error("Error updating user profile", "ID", req.Id, "error", err)
		return nil, errors.New("failed to update profile")
	}

	return &pb.Void{}, nil
}

func (u *UserRepo) UpdateProfileAdmin(ctx context.Context, req *pb.UpdateProfileAdminRequest) (*pb.Void, error) {
	query := `UPDATE users SET first_name = $1, last_name = $2, password_hash = $3, phone_number = $4, date_of_birth = $5, gender = $6, "group" = $7, updated_at = CURRENT_TIMESTAMP 
	WHERE id = $8 AND deleted_at IS NULL`

	_, err := u.DB.Exec(query, req.Firstname, req.Lastname, req.Password, req.Phone, req.DateOfBirth, req.Gender, req.Group, req.Id)
	if err != nil {
		u.Log.Error("Error updating user by admin", "ID", req.Id, "error", err)
		return nil, errors.New("failed to update user profile by admin")
	}

	return &pb.Void{}, nil
}

func (u *UserRepo) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.Void, error) {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	_, err := u.DB.Exec(query, req.Id)
	if err != nil {
		u.Log.Error("Error deleting user profile", "ID", req.Id, "error", err)
		return nil, errors.New("failed to delete profile")
	}

	return &pb.Void{}, nil
}
