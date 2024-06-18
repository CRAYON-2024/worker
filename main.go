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

type UserJob struct {
	Page int
	User entity.UserPreview
}

type PostJob struct {
	Page int
	Post entity.PostPreview
}

const (
	numOfWorkers = 5
)

func main() {
	var (
		userJobs = make(chan UserJob)
		postJobs = make(chan PostJob)
		wg       sync.WaitGroup
	)

	err := godotenv.Load()
	if err != nil {
		log.Panic("error loading .env file")
	}

	initialTime := time.Now()

	// Assign workers to user and posts
	wg.Add(numOfWorkers)

	for worker := 0; worker < numOfWorkers; worker++ {
		go userWorker(&wg, worker, userJobs)
		go postWorker(&wg, worker, postJobs)
	}

	// Adding jobs for users
	go fetchAllUsers(&wg, userJobs)
	// Adding jobs for posts
	go fetchAllPosts(&wg, postJobs)

	wg.Wait()
	// User
	finalTime := time.Now()
	fmt.Printf("Time taken: %v\n", finalTime.Sub(initialTime))
}

func userWorker(wg *sync.WaitGroup, worker int, userJob <-chan UserJob) {
	// job here are read only
	defer wg.Done()

	for job := range userJob {
		go func(user entity.UserPreview) {
			if err := printUserDetail(user); err != nil {
				log.Fatalf("Error printing user detail: %v", err)
			}
		}(job.User)
	}
}

func postWorker(wg *sync.WaitGroup, worker int, postJob <-chan PostJob) {
	defer wg.Done()

	for job := range postJob {
		go func(post entity.PostPreview) {
			if err := printPostDetail(post); err != nil {
				log.Fatalf("Error printing user detail: %v", err)
			}
		}(job.Post)
	}
}

func fetchAllUsers(wg *sync.WaitGroup, userJob chan<- UserJob) error {
	// Loop until 10th page (no validation, as the inferred requirement)
	defer close(userJob)

	for page := 0; page < 10; page++ {
		if err := fetchUsersPage(page, userJob); err != nil {
			log.Fatalf("Error fetching users page: %v", err)
		}
	}

	return nil
}

func fetchUsersPage(page int, userJob chan<- UserJob) error {
	userResponse, err := internal.FetchUsers(page)

	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	for _, user := range userResponse.Data {
		go func(u entity.UserPreview) {
			userJob <- UserJob{Page: page, User: u}
		}(user)
	}

	return nil
}

func printUserDetail(user entity.UserPreview) error {
	userDetail, err := internal.FetchUserDetail(user.ID)

	if err != nil {
		return fmt.Errorf("error printing user detail: %w", err)
	}

	fmt.Printf("User name %s %s %s %s %s \n", user.Title, user.FirstName, user.LastName, userDetail.Email, userDetail.Gender)

	return nil
}

func fetchAllPosts(wg *sync.WaitGroup, postJob chan<- PostJob) error {
	defer close(postJob)
	// Loop until 10th page ( no validation, as the inferred requirement )

	for page := 0; page < 10; page++ {
		if err := fetchPostsPage(page, postJob); err != nil {
			log.Fatalf("Error fetching posts page: %v", err)
		}
	}

	return nil
}

func fetchPostsPage(page int, postJob chan<- PostJob) error {
	postResponse, err := internal.FetchPosts(page)

	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	for _, post := range postResponse.Data {
		postJob <- PostJob{Page: page, Post: post}
	}

	return nil
}

func printPostDetail(post entity.PostPreview) error {
	fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.Owner.FirstName, post.Owner.LastName, post.Text, post.Likes, post.Tags, post.PublishDate.Format("2006-01-02 15:04:05"))

	return nil
}
