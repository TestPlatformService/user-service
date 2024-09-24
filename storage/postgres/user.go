package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	pb "user/genproto/user"
	"user/logs"
	"user/storage"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewUserRepo(DB *sql.DB) storage.IUserStorage {
	return &UserRepo{DB: DB, Log: logs.NewLogger()}
}

func (u *UserRepo) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Void, error) {
	// Check if the phone number or hh_id already exists
	var existingUserCount int
	err := u.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE phone_number = $1 OR hh_id = $2", req.Phone, req.HhId).Scan(&existingUserCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing users: %w", err)
	}
	if existingUserCount > 0 {
		return nil, fmt.Errorf("user with this phone number or hh_id already exists")
	}

	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert the new user into the database
	query := `
        INSERT INTO users (hh_id, first_name, last_name, password_hash, phone_number, gender, date_of_birth, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err = u.DB.ExecContext(ctx, query,
		req.HhId,
		req.Firstname,
		req.Lastname,
		hashedPassword,
		req.Phone,
		req.Gender,
		req.DateOfBirth,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Return success
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
		Id:   user.Id,
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
		HhId:        user.HhId,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Password:    user.Password,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
		Id:          user.Id,
	}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	// Start with base query
	query := `SELECT hh_id, first_name, last_name, phone_number, date_of_birth, gender, id, role 
		FROM users 
		WHERE deleted_at IS NULL`

	// List of query parameters
	var params []interface{}
	var conditions []string
	paramIndex := 1

	// Filter by role
	if req.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", paramIndex))
		params = append(params, req.Role)
		paramIndex++
	}

	// Filter by group (join with groups if necessary)
	if req.Group != "" {
		query += ` JOIN student_groups sg ON sg.student_hh_id = users.hh_id 
			 JOIN groups g ON g.id = sg.group_id`
		conditions = append(conditions, fmt.Sprintf("g.name = $%d", paramIndex))
		params = append(params, req.Group)
		paramIndex++
	}

	// Filter by subject (assuming subject is linked to groups)
	if req.Subject != "" {
		query += ` JOIN groups g ON g.id = sg.group_id`
		conditions = append(conditions, fmt.Sprintf("g.subject_id = $%d", paramIndex))
		params = append(params, req.Subject)
		paramIndex++
	}

	// Filter by teacher (join with teacher_groups)
	if req.Teacher != "" {
		query += ` JOIN teacher_groups tg ON tg.group_id = g.id`
		conditions = append(conditions, fmt.Sprintf("tg.teacher_id = $%d", paramIndex))
		params = append(params, req.Teacher)
		paramIndex++
	}

	// Filter by hh_id
	if req.HhId != "" {
		conditions = append(conditions, fmt.Sprintf("hh_id = $%d", paramIndex))
		params = append(params, req.HhId)
		paramIndex++
	}

	// Filter by phone_number
	if req.PhoneNumber != "" {
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", paramIndex))
		params = append(params, req.PhoneNumber)
		paramIndex++
	}

	// Filter by gender
	if req.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", paramIndex))
		params = append(params, req.Gender)
		paramIndex++
	}

	// If there are conditions, append them to the query
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Add limit and offset for pagination
	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIndex)
		params = append(params, req.Limit)
		paramIndex++
	}
	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", paramIndex)
		params = append(params, req.Offset)
		paramIndex++
	}

	// Execute the query
	rows, err := u.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Initialize response
	var users []*pb.GetProfileResponse
	for rows.Next() {
		var user pb.GetProfileResponse
		err := rows.Scan(&user.HhId, &user.Firstname, &user.Lastname, &user.Phone, &user.DateOfBirth, &user.Gender, &user.Id, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, &user)
	}

	// Get total count of users (without limit and offset)
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}
	err = u.DB.QueryRowContext(ctx, countQuery, params[:len(params)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Return the response
	return &pb.GetAllUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       req.Offset/req.Limit + 1,
		Limit:      req.Limit,
	}, nil
}

func (u *UserRepo) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Void, error) {
	// Initialize the base query
	query := "UPDATE users SET "
	var params []interface{}
	var setClauses []string
	paramIndex := 1

	// Conditionally update profile picture if provided
	if req.ProfilePicture != "" {
		setClauses = append(setClauses, fmt.Sprintf("profile_image = $%d", paramIndex))
		params = append(params, req.ProfilePicture)
		paramIndex++
	}

	// Conditionally update password if provided
	if req.Password != "" {
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", paramIndex))
		hashedPassword, err := hashPassword(req.Password) // Use the hashPassword function
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		params = append(params, hashedPassword)
		paramIndex++
	}

	// If no fields to update, return an error
	if len(setClauses) == 0 {
		return nil, fmt.Errorf("nothing to update")
	}

	// Append updated_at field
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramIndex))
	params = append(params, time.Now()) // Update the `updated_at` timestamp
	paramIndex++

	// Build the final query
	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, req.Id)

	// Execute the query
	_, err := u.DB.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Return a success response
	return &pb.Void{}, nil
}

func (u *UserRepo) UpdateProfileAdmin(ctx context.Context, req *pb.UpdateProfileAdminRequest) (*pb.Void, error) {
	// Start building the query
	query := "UPDATE users SET"
	params := []interface{}{}
	paramCounter := 1

	// Update first name if provided
	if req.Firstname != "" {
		query += fmt.Sprintf(" first_name = $%d,", paramCounter)
		params = append(params, req.Firstname)
		paramCounter++
	}

	// Update last name if provided
	if req.Lastname != "" {
		query += fmt.Sprintf(" last_name = $%d,", paramCounter)
		params = append(params, req.Lastname)
		paramCounter++
	}

	// Update password if provided (and hash it)
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		query += fmt.Sprintf(" password_hash = $%d,", paramCounter)
		params = append(params, hashedPassword)
		paramCounter++
	}

	// Update phone number if provided
	if req.Phone != "" {
		query += fmt.Sprintf(" phone_number = $%d,", paramCounter)
		params = append(params, req.Phone)
		paramCounter++
	}

	// Update date of birth if provided
	if req.DateOfBirth != "" {
		query += fmt.Sprintf(" date_of_birth = $%d,", paramCounter)
		params = append(params, req.DateOfBirth)
		paramCounter++
	}

	// Update gender if provided
	if req.Gender != "" {
		query += fmt.Sprintf(" gender = $%d,", paramCounter)
		params = append(params, req.Gender)
		paramCounter++
	}

	// If no fields were updated, return an error
	if paramCounter == 1 {
		return nil, fmt.Errorf("no valid fields provided for update")
	}

	// Remove the trailing comma and add the WHERE clause
	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", paramCounter)
	params = append(params, req.Id)

	// Execute the query
	_, err := u.DB.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
