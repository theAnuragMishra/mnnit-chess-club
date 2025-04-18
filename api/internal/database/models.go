// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"time"
)

type CsrfToken struct {
	SessionID string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Game struct {
	ID               string
	BaseTime         int32
	Increment        int32
	WhiteID          *int32
	BlackID          *int32
	Fen              string
	GameLength       int16
	Result           string
	CreatedAt        time.Time
	EndTimeLeftWhite *int32
	EndTimeLeftBlack *int32
	ResultReason     *string
	RatingW          int32
	RatingB          int32
	ChangeW          *int32
	ChangeB          *int32
}

type Move struct {
	ID           int32
	GameID       string
	MoveNumber   int32
	PlayerID     *int32
	MoveNotation string
	Orig         string
	Dest         string
	MoveFen      string
	CreatedAt    time.Time
}

type Session struct {
	ID        string
	UserID    int32
	CreatedAt time.Time
	ExpiresAt time.Time
}

type User struct {
	ID         int32
	Email      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Username   *string
	AvatarUrl  *string
	GoogleID   string
	Rating     float64
	Rd         float64
	Volatility float64
}
