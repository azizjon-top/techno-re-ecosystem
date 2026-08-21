package models

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	ChannelID   uuid.UUID `db:"channel_id" json:"channel_id"`
	OwnerID     uuid.UUID `db:"owner_id" json:"owner_id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Video struct {
	VideoID         uuid.UUID `db:"video_id" json:"video_id"`
	ChannelID       uuid.UUID `db:"channel_id" json:"channel_id"`
	Title           string    `db:"title" json:"title"`
	Description     *string   `db:"description" json:"description"`
	FilePath        string    `db:"file_path" json:"file_path"`
	FileHash        *string   `db:"file_hash" json:"file_hash"`
	DurationSeconds *int      `db:"duration_seconds" json:"duration_seconds"`
	ViewCount       int       `db:"view_count" json:"view_count"`
	P2PEnabled      bool      `db:"p2p_enabled" json:"p2p_enabled"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type VideoView struct {
	ViewID         uuid.UUID `db:"view_id" json:"view_id"`
	VideoID        uuid.UUID `db:"video_id" json:"video_id"`
	ViewerID       uuid.UUID `db:"viewer_id" json:"viewer_id"`
	WatchedSeconds int       `db:"watched_seconds" json:"watched_seconds"`
	P2PContributed bool      `db:"p2p_contributed" json:"p2p_contributed"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
