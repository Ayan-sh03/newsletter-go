package models

type Newsletter struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      int64  `json:"author"`
}
