package handler

import (
	"log/slog"
	"user/genproto/user"
)

type Handler struct {
	User user.UsersClient
	Log  *slog.Logger
}
