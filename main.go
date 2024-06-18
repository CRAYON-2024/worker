package main

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/CRAYON-2024/worker/internal"
	"github.com/CRAYON-2024/worker/internal/entity"
	"github.com/joho/godotenv"
)

func main() {
	var (
		wg = sync.WaitGroup{}
	)

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	initialTime := time.Now()

	wg.Add(2)
	// User

	go func() {
		defer wg.Done()
		fetchAllUser()
	}()

	// Post
	go func() {
		defer wg.Done()
		fetchAllPost()
	}()

	wg.Wait()

	finalTime := time.Now()
	fmt.Printf("Time taken: %v\n", finalTime.Sub(initialTime))
}

func fetchAllUser() {
	var wg sync.WaitGroup

	// Loop until 10th page (no validation, as the inferred requirement)
	for page := 0; page < 10; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			userResponse, err := internal.FetchUsers(page)

			if err != nil {
				log.Fatalf("Error fetching users: %v", err)
			}

			var innerWg sync.WaitGroup
			for _, user := range userResponse.Data {
				innerWg.Add(1)
				go func(user entity.UserPreview) {
					defer innerWg.Done()
					userDetail, err := internal.FetchUserDetail(user.ID)

					if err != nil {
						log.Fatalf("Error fetching user detail: %v", err)
					}

					fmt.Printf("User name %s %s %s %s %s \n", user.Title, user.FirstName, user.LastName, userDetail.Email, userDetail.Gender)
				}(user)
			}
			innerWg.Wait()
		}(page)
	}
	wg.Wait()
}


func fetchAllPost() {
	var wg sync.WaitGroup
	// Loop until 10th page ( no validation, as the inferred requirement )

	for page := 0; page < 10; page++ {
		wg.Add(1)

		go func (pg int) {
			defer wg.Done()
			postResponse, err := internal.FetchPosts(pg)

			if err != nil {
				log.Fatalf("Error fetching posts: %v", err)
			}

			for _, post := range postResponse.Data {
				fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.Owner.FirstName, post.Owner.LastName, post.Text, post.Likes, post.Tags, post.PublishDate.Format("2006-01-02 15:04:05"))
			}
		}(page)
	}

	wg.Wait()
}
