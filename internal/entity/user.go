package entity

import (
	"time"
)

type (
	UserResponse struct {
		Data  []UserPreview `json:"data"`
		Total int           `json:"total"`
		Page  int           `json:"page"`
		Limit int           `json:"limit"`
	}

	UserPreview struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Picture   string `json:"picture"`
	}

	User struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		FirstName   string    `json:"firstName"`
		LastName    string    `json:"lastName"`
		Picture     string    `json:"picture"`
		Gender      string    `json:"gender"`
		Email       string    `json:"email"`
		DateOfBirth time.Time `json:"dateOfBirth"`
		Phone       string    `json:"phone"`
		Location    struct {
			Street   string `json:"street"`
			City     string `json:"city"`
			State    string `json:"state"`
			Country  string `json:"country"`
			Timezone string `json:"timezone"`
		} `json:"location"`
		RegisterDate time.Time `json:"registerDate"`
		UpdatedDate  time.Time `json:"updatedDate"`
	}

	CustomUser struct {
		ID   int
		Name string
	}
	
)
