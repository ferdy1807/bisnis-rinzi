// File: packages/backend/outbox/outbox.go
package outbox

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Status tipe khusus untuk membatasi status literal outbox agar aman dari typo.
type EventStatus string

const (
	StatusPending EventStatus = "PENDING"
	StatusSuccess EventStatus = "PROCESSED"
	StatusFailed  EventStatus = "FAILED"
)

// Event merepresentasikan struktur standar tabel outbox_events lintas-skema.
type Event struct {
	ID            string      `json:"id"`
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	EventType     string      `json:"event_type"`
	Payload       []byte      `json:"payload"`
	Status        EventStatus `json:"status"`
	ErrorMessage  *string     `json:"error_message,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	ProcessedAt   *time.Time  `json:"processed_at,omitempty"`
}

// CreateEvent membungkus entitas domain menjadi objek Outbox Event standar terenkapsulasi JSON.
func CreateEvent(aggregateType, aggregateID, eventType string, payloadData interface{}) (*Event, error) {
	binaryPayload, err := json.Marshal(payloadData)
	if err != nil {
		return nil, fmt.Errorf("outbox_lib: gagal melakukan serialisasi payload ke json: %w", err)
	}

	now := time.Now().UTC()
	return &Event{
		ID:            generateUUIDv4(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       binaryPayload,
		Status:        StatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// SaveEventTx menyimpan entri event outbox ke dalam skema database lokal pembawa transaksi bisnis.
// Harus dijalankan di dalam blok transaksi pgx.Tx yang sama dengan mutasi bisnis untuk menjamin ACID.
func SaveEventTx(ctx context.Context, tx pgx.Tx, e *Event) error {
	query := `INSERT INTO outbox_events (id, aggregate_type, aggregate_id, event_type, payload, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id) DO NOTHING`
	_, err := tx.Exec(ctx, query, e.ID, e.AggregateType, e.AggregateID, e.EventType, e.Payload, string(e.Status), e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("outbox_lib: gagal menyimpan rekor event ke database: %w", err)
	}
	return nil
}

// FetchPendingEvents mengambil daftar antrean event berstatus PENDING secara berurutan (FIFO).
func FetchPendingEvents(ctx context.Context, tx pgx.Tx, limit int) ([]*Event, error) {
	query := `SELECT id, aggregate_type, aggregate_id, event_type, payload, status, error_message, created_at, updated_at, processed_at 
	          FROM outbox_events 
	          WHERE status = 'PENDING' 
	          ORDER BY created_at ASC 
	          LIMIT $1 FOR UPDATE SKIP LOCKED`

	rows, err := tx.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("outbox_lib: gagal mengeksekusi kueri pengambilan antrean: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		var statusStr string
		err := rows.Scan(&e.ID, &e.AggregateType, &e.AggregateID, &e.EventType, &e.Payload, &statusStr, &e.ErrorMessage, &e.CreatedAt, &e.UpdatedAt, &e.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("outbox_lib: gagal memindai baris data outbox: %w", err)
		}
		e.Status = EventStatus(statusStr)
		events = append(events, &e)
	}
	return events, nil
}

// UpdateEventStatusTx memperbarui status penyelesaian event pasca diproses atau didistribusikan oleh worker.
func UpdateEventStatusTx(ctx context.Context, tx pgx.Tx, id string, status EventStatus, errMsg *string) error {
	now := time.Now().UTC()
	var query string

	if status == StatusSuccess {
		query = `UPDATE outbox_events SET status = $1, error_message = NULL, updated_at = $2, processed_at = $3 WHERE id = $4`
		_, err := tx.Exec(ctx, query, string(status), now, now, id)
		return err
	}

	query = `UPDATE outbox_events SET status = $1, error_message = $2, updated_at = $3, processed_at = $4 WHERE id = $5`
	_, err := tx.Exec(ctx, query, string(status), errMsg, now, now, id)
	return err
}

func generateUUIDv4() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Set UUID version to 4
	b[8] = (b[8] & 0x3f) | 0x80 // Set UUID variant to RFC4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
