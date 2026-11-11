package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"api-students/app/model"
	"api-students/config"
	"api-students/database"
)

func main() {
	config.LoadEnv()

	outboxID := flag.Int64("outbox", 0, "ID baris outbox yang ingin dikirim ulang")
	copies := flag.Int("copies", 1, "Jumlah salinan task yang dikirim")
	flag.Parse()

	if *outboxID <= 0 {
		log.Fatal("flag -outbox wajib diisi dengan ID valid")
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("koneksi db gagal: %v", err)
	}
	defer pool.Close()

	var event model.OutboxEvent
	err = pool.QueryRow(ctx, `
		SELECT id, event_type, payload, created_at
		FROM outbox WHERE id = $1`, *outboxID).
		Scan(&event.ID, &event.EventType, &event.Payload, &event.CreatedAt)
	if err != nil {
		log.Fatalf("gagal mengambil outbox ID %d: %v", *outboxID, err)
	}

	client := asynq.NewClient(config.RedisOpt())
	defer client.Close()

	job := model.NotificationJob{
		OutboxID:  event.ID,
		EventType: event.EventType,
		Payload:   event.Payload,
	}

	payloadBytes, err := json.Marshal(job)
	if err != nil {
		log.Fatalf("marshal payload gagal: %v", err)
	}

	fmt.Printf("Mengirim %d salinan task untuk outbox ID %d ke Redis...\n", *copies, *outboxID)

	for i := 1; i <= *copies; i++ {
		// Menggunakan TaskID unik per salinan agar tidak ditolak ErrTaskIDConflict oleh Redis saat pengujian idempotency
		task := asynq.NewTask(
			model.TaskSendNotification,
			payloadBytes,
			asynq.TaskID(fmt.Sprintf("redeliver-%d-copy-%d-%d", *outboxID, i, time.Now().UnixNano())),
			asynq.Queue(model.QueueNotifications),
			asynq.MaxRetry(5),
			asynq.Timeout(30*time.Second),
			asynq.Retention(24*time.Hour),
		)

		info, err := client.EnqueueContext(ctx, task)
		if err != nil {
			log.Fatalf("gagal enqueue salinan ke-%d: %v", i, err)
		}
		fmt.Printf(" [%d/%d] Enqueued TaskID: %s\n", i, *copies, info.ID)
	}

	fmt.Println("Selesai.")
}