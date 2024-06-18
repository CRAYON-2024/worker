package main

import (
	"fmt"
	"log"
	"time"

	"github.com/CRAYON-2024/worker/internal"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	initialTime := time.Now()
	// User
	go fetchAllUser()

	// Post
	go fetchAllPost()

	finalTime := time.Now()
	fmt.Printf("Time taken: %v\n", finalTime.Sub(initialTime))
}

func fetchAllUser() {
	// Loop until 10th page ( no validation, as the inferred requirement )
	for page := 0; page < 10; page++ {
		userResponse, err := internal.FetchUsers(page)

		if err != nil {
			log.Fatalf("Error fetching users: %v", err)
		}

		for _, user := range userResponse.Data {
			userDetail, err := internal.FetchUserDetail(user.ID)

			if err != nil {
				log.Fatalf("Error fetching user detail: %v", err)
			}

			fmt.Printf("User name %s %s %s %s %s\n", user.Title, user.FirstName, user.LastName, userDetail.Email, userDetail.Gender)
		}
	}
}

func fetchAllPost() {
	// Loop until 10th page ( no validation, as the inferred requirement )

	for page := 0; page < 10; page++ {
		postResponse, err := internal.FetchPosts(page)

		if err != nil {
			log.Fatalf("Error fetching posts: %v", err)
		}

		for _, post := range postResponse.Data {
			fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.Owner.FirstName, post.Owner.LastName, post.Text, post.Likes, post.Tags, post.PublishDate.Format("2006-01-02 15:04:05"))
		}
	}
}
