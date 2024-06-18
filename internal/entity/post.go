package entity

import "time"

type (
	PostResponse struct {
		Data  []PostPreview `json:"data"`
		Total int           `json:"total"`
		Page  int           `json:"page"`
		Limit int           `json:"limit"`
	}

	PostPreview struct {
		ID          string      `json:"id"`
		Image       string      `json:"image"`
		Likes       int         `json:"likes"`
		Tags        []string    `json:"tags"`
		Text        string      `json:"text"`
		PublishDate time.Time   `json:"publishDate"`
		Owner       UserPreview `json:"owner"`
	}
)
