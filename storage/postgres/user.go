package postgres

import (
	"database/sql"
	"log/slog"
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

