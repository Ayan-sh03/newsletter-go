package models

type Subscriber struct {
	Id           int64 `json:"id"`
	UserId       int64 `json:"user_id"`
	NewsletterId int64 `json:"newsletter_id"`
}
