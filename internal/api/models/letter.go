package models

type Letter struct {
	Id           int64  `json:"id"`
	NewsletterId int64  `json:"newsletter_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
