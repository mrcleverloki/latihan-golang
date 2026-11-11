package model

import (
	"encoding/json"
	"time"
)

// Nama-nama event yang dicatat ke outbox.
const (
	EventEnrollmentCreated = "enrollment.created"
	EventEnrollmentDropped = "enrollment.dropped"
	EventEnrollmentSwapped = "enrollment.swapped"
)

// TaskSendNotification adalah nama task type, QueueNotifications adalah nama queue.
const (
	TaskSendNotification = "notification:send"
	QueueNotifications   = "notifications"
)

// EnrollmentEvent adalah isi event KRS.
type EnrollmentEvent struct {
	UserID         int    `json:"user_id"`
	CourseID       int    `json:"course_id"`
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	FromCourseCode string `json:"from_course_code,omitempty"`
	FromCourseName string `json:"from_course_name,omitempty"`
}

// OutboxEvent merepresentasikan satu baris tabel outbox.
type OutboxEvent struct {
	ID        int64
	EventType string
	Payload   []byte
	CreatedAt time.Time
}

// NotificationJob adalah isi task yang dikirim ke Redis.
type NotificationJob struct {
	OutboxID  int64           `json:"outbox_id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// Message adalah pesan umum yang siap dikirim.
type Message struct {
	To      string
	Subject string
	Body    string
}

// Notification adalah data notifikasi yang tersimpan di DB.
type Notification struct {
	ID       int64     `json:"id"`
	OutboxID int64     `json:"outbox_id"`
	UserID   int       `json:"user_id"`
	Channel  string    `json:"channel"`
	Subject  string    `json:"subject"`
	Body     string    `json:"body"`
	SentAt   time.Time `json:"sent_at"`
}

// QueueStats adalah ringkasan monitoring queue.
type QueueStats struct {
	OutboxPending     int        `json:"outbox_pending"`
	OutboxPublished   int        `json:"outbox_published"`
	OldestPendingAge  string     `json:"oldest_pending_age"`
	NotificationsSent int        `json:"notifications_sent"`
	Queue             *QueueInfo `json:"queue"`
	QueueError        string     `json:"queue_error,omitempty"`
}

// QueueInfo salinan info antrean dari Asynq.
type QueueInfo struct {
	Name      string `json:"name"`
	Pending   int    `json:"pending"`
	Active    int    `json:"active"`
	Scheduled int    `json:"scheduled"`
	Retry     int    `json:"retry"`
	Archived  int    `json:"archived"`
	Completed int    `json:"completed"`
	Processed int    `json:"processed_today"`
	Failed    int    `json:"failed_today"`
}