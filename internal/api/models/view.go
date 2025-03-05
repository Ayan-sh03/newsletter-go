package models

type View struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	NewsletterId int64  `json:"newsletter_id"`
	LetterId     int64  `json:"letter_id"`
	ViewedAt     string `json:"viewed_at"`
}
