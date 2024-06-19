package cmd

import (
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/CRAYON-2024/worker/internal/api"
	"github.com/CRAYON-2024/worker/internal/entity"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func WorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Command to fetch users and posts concurrently",
		Run: func(cmd *cobra.Command, args []string) {
			initialTime := time.Now()
			runWorker()
			finalTime := time.Now()
			fmt.Printf("Time taken: %v\n", finalTime.Sub(initialTime))
		},
	}
}

func runWorker() {
	var (
		wg = sync.WaitGroup{}
	)

	err := godotenv.Load()
	if err != nil {
		log.Panic("error loading .env file")
	}

	wg.Add(2)
	// User

	go func() {
		defer wg.Done()
		if err := fetchAllUsers(); err != nil {
			log.Printf("error in fetchAllUsers: %v", err)
		}
	}()

	// Post
	go func() {
		defer wg.Done()
		if err := fetchAllPosts(); err != nil {
			log.Printf("Error in fetchAllPosts: %v", err)
		}
	}()

	wg.Wait()
}

func fetchAllUsers() error {
	var wg sync.WaitGroup

	// Loop until 10th page (no validation, as the inferred requirement)

	for page := 0; page < 10; page++ {
		wg.Add(1)
		go func(pg int) {
			defer wg.Done()
			if err := fetchUsersPage(pg); err != nil {
				log.Fatalf("Error fetching users page: %v", err)
			}
		}(page)
	}
	wg.Wait()

	return nil
}

func fetchUsersPage(page int) error {
	userResponse, err := api.FetchUsers(page)

	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	var innerWg sync.WaitGroup

	for _, user := range userResponse.Data {
		innerWg.Add(1)
		go func(u entity.UserPreview) {
			defer innerWg.Done()
			if error := printUserDetail(u); error != nil {
				log.Printf("Error printing user details for user %s: %v", u.ID, err)
			}
		}(user)
	}

	innerWg.Wait()

	return nil
}

func printUserDetail(user entity.UserPreview) error {
	userDetail, err := api.FetchUserDetail(user.ID)

	if err != nil {
		return fmt.Errorf("error printing user detail: %w", err)
	}

	fmt.Printf("User name %s %s %s %s %s \n", user.Title, user.FirstName, user.LastName, userDetail.Email, userDetail.Gender)

	return nil
}

func fetchAllPosts() error {
	var wg sync.WaitGroup
	// Loop until 10th page ( no validation, as the inferred requirement )

	for page := 0; page < 10; page++ {
		wg.Add(1)

		go func(pg int) {
			defer wg.Done()
			if err := fetchPostsPage(pg); err != nil {
				log.Fatalf("Error fetching posts page: %v", err)
			}
		}(page)
	}

	wg.Wait()

	return nil
}

func fetchPostsPage(page int) error {
	postResponse, err := api.FetchPosts(page)

	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	for _, post := range postResponse.Data {
		fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.Owner.FirstName, post.Owner.LastName, post.Text, post.Likes, post.Tags, post.PublishDate.Format("2006-01-02 15:04:05"))
	}

	return nil
}
