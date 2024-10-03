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
        INSERT INTO users (hh_id, first_name, last_name, password_hash, phone_number, gender, date_of_birth, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err = u.DB.ExecContext(ctx, query,
		req.HhId,
		req.Firstname,
		req.Lastname,
		hashedPassword,
		req.Phone,
		req.Gender,
		req.DateOfBirth,
		req.Role,
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
	var userID string
	var hashedPassword string
	var role string
	var isAdmin bool

	// Check if the user exists in the 'users' table
	err := u.DB.QueryRowContext(ctx, "SELECT id, password_hash, role FROM users WHERE hh_id = $1", req.HhId).Scan(&userID, &hashedPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found in the 'users' table, check 'admin' table
			err = u.DB.QueryRowContext(ctx, "SELECT id, password, role FROM admin WHERE hh_id = $1", req.HhId).Scan(&userID, &hashedPassword, &role)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, fmt.Errorf("user with hh_id %s not found", req.HhId)
				}
				return nil, fmt.Errorf("error querying admin table: %w", err)
			}
			isAdmin = true
		} else {
			return nil, fmt.Errorf("error querying users table: %w", err)
		}
	}

	// Compare passwords
	if isAdmin {
		// For admin, compare the plain text password (since it's not hashed)
		if hashedPassword != req.Password {
			return nil, fmt.Errorf("invalid password")
		}
	} else {
		// For regular users, compare the hashed password
		err = comparePassword(hashedPassword, req.Password)
		if err != nil {
			return nil, fmt.Errorf("invalid password")
		}
	}

	// Return the user's ID and role
	return &pb.LoginResponse{
		Id:   userID,
		Role: role,
	}, nil
}

func (u *UserRepo) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var user pb.GetProfileResponse
	var profileImage *string // Use a pointer to handle NULL values

	query := `SELECT hh_id, first_name, last_name, password_hash, phone_number, date_of_birth, gender, profile_image 
	FROM users 
	WHERE id = $1 AND deleted_at IS NULL`

	err := u.DB.QueryRow(query, req.Id).Scan(&user.HhId, &user.Firstname, &user.Lastname, &user.Password, &user.Phone, &user.DateOfBirth, &user.Gender, &profileImage)
	if err == sql.ErrNoRows {
		u.Log.Error("No user found", "ID", req.Id)
		return nil, errors.New("no user found")
	} else if err != nil {
		u.Log.Error("Error getting user by ID", "err", err)
		return nil, err
	}

	// Handle the case where profile_image might be NULL
	if profileImage != nil {
		user.Photo = *profileImage // Dereference if not NULL
	} else {
		user.Photo = "" // Or set to a default value
	}

	return &pb.GetProfileResponse{
		HhId:        user.HhId,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Password:    user.Password,
		Phone:       user.Phone,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
		Id:          req.Id,
		Photo:       user.Photo, // Ensure photo is included
	}, nil
}

func (u *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	query := `SELECT hh_id, first_name, last_name, phone_number, date_of_birth, gender, id, role, profile_image 
        FROM users 
        WHERE deleted_at IS NULL`

	var params []interface{}
	var conditions []string
	paramIndex := 1

	if req.Role != "" {
		conditions = append(conditions, fmt.Sprintf("role = $%d", paramIndex))
		params = append(params, req.Role)
		paramIndex++
	}
	if req.Group != "" {
		query += ` JOIN student_groups sg ON sg.student_hh_id = users.hh_id 
                JOIN groups g ON g.id = sg.group_id`
		conditions = append(conditions, fmt.Sprintf("g.name = $%d", paramIndex))
		params = append(params, req.Group)
		paramIndex++
	}
	if req.Subject != "" {
		query += ` JOIN groups g ON g.id = sg.group_id`
		conditions = append(conditions, fmt.Sprintf("g.subject_id = $%d", paramIndex))
		params = append(params, req.Subject)
		paramIndex++
	}
	if req.Teacher != "" {
		query += ` JOIN teacher_groups tg ON tg.group_id = g.id`
		conditions = append(conditions, fmt.Sprintf("tg.teacher_id = $%d", paramIndex))
		params = append(params, req.Teacher)
		paramIndex++
	}
	if req.HhId != "" {
		conditions = append(conditions, fmt.Sprintf("hh_id = $%d", paramIndex))
		params = append(params, req.HhId)
		paramIndex++
	}
	if req.PhoneNumber != "" {
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", paramIndex))
		params = append(params, req.PhoneNumber)
		paramIndex++
	}
	if req.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", paramIndex))
		params = append(params, req.Gender)
		paramIndex++
	}
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIndex)
		params = append(params, req.Limit)
		paramIndex++
	}
	offset := (req.Page - 1) * req.Limit
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", paramIndex)
		params = append(params, (offset-1)*req.Limit)
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
		var profileImage *string
		err := rows.Scan(&user.HhId, &user.Firstname, &user.Lastname, &user.Phone, &user.DateOfBirth, &user.Gender, &user.Id, &user.Role, &profileImage)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if profileImage != nil {
			user.Photo = *profileImage
		} else {
			user.Photo = ""
		}

		users = append(users, &user)
	}

	// Prepare count query with conditions
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Only use the original params if there are conditions
	var countParams []interface{}
	if len(conditions) > 0 {
		countParams = params // Use the full params if there are conditions
	} else {
		countParams = []interface{}{} // No params needed for count query
	}

	// Execute count query
	var totalCount int64
	err = u.DB.QueryRowContext(ctx, countQuery, countParams...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	return &pb.GetAllUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

func (u *UserRepo) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.Void, error) {
	// Initialize the base query


	// Conditionally update password if provided
	if req.Id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Initialize the query parts
	query := `UPDATE users SET `
	var params []interface{}
	var updates []string
	paramIndex := 1

	// Check if password is provided
	if req.Password != "" {
		// Hash the new password
		hashedPassword, err := hashPassword(req.Password) // Assume `hashPassword` is a function to hash the password
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates = append(updates, fmt.Sprintf("password_hash = $%d", paramIndex))
		params = append(params, hashedPassword)
		paramIndex++
	}

	// Always update the updated_at field
	updates = append(updates, fmt.Sprintf("updated_at = $%d", paramIndex))
	params = append(params, time.Now())
	paramIndex++

	// Ensure there are updates to be made
	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates provided")
	}

	// Complete the query
	query += strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", paramIndex)
	params = append(params, req.Id)

	// Execute the update query
	_, err := u.DB.ExecContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Return success
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

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *UserRepo) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.Void, error) {
	// Validate that both user ID and photo are provided
	if req.Id == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if req.Photo == "" {
		return nil, fmt.Errorf("photo is required")
	}

	// Update the user's profile_image in the database
	query := `UPDATE users SET profile_image = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := u.DB.ExecContext(ctx, query, req.Photo, time.Now(), req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile image: %w", err)
	}

	// Return success
	return &pb.Void{}, nil
}

func (u *UserRepo) DeletePhoto(ctx context.Context, req *pb.DeletePhotoRequest) (*pb.Void, error) {
	// Validate that the user ID is provided
	if req.Id == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Update the user's profile_image to NULL
	query := `UPDATE users SET profile_image = NULL, updated_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := u.DB.ExecContext(ctx, query, time.Now(), req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete profile image: %w", err)
	}

	// Return success
	return &pb.Void{}, nil
}
