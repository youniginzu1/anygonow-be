package notification

import "context"

type NotificationRepo interface {
	InsertNotification(context.Context, *Notification) (*Notification, error)
	SelectNotification(context.Context, *Search) (*Notification, error)
	UpdateNotification(context.Context, *Search, *Notification) error
	UpsertNotification(context.Context, *Notification) error
}
