package postgres

import (
	"context"
	"database/sql"
	"fmt"
	pb "user/genproto/notification"
	"user/storage"
)

type NotificationsRepository struct {
	Db *sql.DB
}

func NewNotificationsRepository(db *sql.DB) storage.INotificationStorage {
	return &NotificationsRepository{Db: db}
}

func (r *NotificationsRepository) CreateNotifications(ctx context.Context, req *pb.CreateNotificationsReq) (*pb.CreateNotificationsRes, error) {
	var id string
	err := r.Db.QueryRowContext(ctx, `
		INSERT INTO notifications (user_id, messages, is_read) 
		VALUES ($1, $2, $3) RETURNING id`,
		req.GetUserId(), req.GetMessage(), false,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert notification: %w", err)
	}

	return &pb.CreateNotificationsRes{Id: id}, nil
}

func (r *NotificationsRepository) GetAllNotifications(ctx context.Context, req *pb.GetNotificationsReq) (*pb.GetNotificationsResponse, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT id, user_id, messages, created_at 
		FROM notifications 
		WHERE user_id = $1`, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*pb.Notification
	for rows.Next() {
		var notification pb.Notification
		if err := rows.Scan(&notification.Id, &notification.UserId, &notification.Message, &notification.Date); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &pb.GetNotificationsResponse{Notifications: notifications}, nil
}

func (r *NotificationsRepository) GetAndMarkNotificationAsRead(ctx context.Context, req *pb.GetAndMarkNotificationAsReadReq) (*pb.GetAndMarkNotificationAsReadRes, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT id, user_id, messages, created_at 
		FROM notifications 
		WHERE user_id = $1 AND is_read = FALSE`, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*pb.Notification
	for rows.Next() {
		var notification pb.Notification
		if err := rows.Scan(&notification.Id, &notification.UserId, &notification.Message, &notification.Date); err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, notif := range notifications {
		_, err := tx.ExecContext(ctx, `
			UPDATE notifications 
			SET is_read = TRUE 
			WHERE id = $1`, notif.GetId())
		if err != nil {
			return nil, fmt.Errorf("failed to update notification %s: %w", notif.GetId(), err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.GetAndMarkNotificationAsReadRes{Notifications: notifications}, nil
}
