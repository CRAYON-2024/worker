package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"sync"

	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/CRAYON-2024/worker/internal/api"
	"github.com/CRAYON-2024/worker/internal/entity"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UserJob struct {
	Page int
	User entity.UserPreview
}

type PostJob struct {
	Page int
	Post entity.PostPreview
}

var (
	numOfWorkers int
)

func WorkerCommand(container *bootstrap.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Command to fetch users and posts concurrently",
		Run: func(cmd *cobra.Command, args []string) {
			if numOfWorkers < 1 {
				log.Fatalf("Number of workers must be greater than 0")
			}

			initialTime := time.Now()
			// container.GetKafkaProducer()
			runWorker(
				viper.GetString("kafka.topic.worker"),
				container.GetKafkaProducer(),
			)

			finalTime := time.Now()
			fmt.Printf("Time taken: %v\n", finalTime.Sub(initialTime))
		},
	}

	cmd.Flags().IntVarP(&numOfWorkers, "worker", "w", 5, "Define how many worker to do these tasks")

	return cmd
}

func runWorker(
	topic string,
	producer *kafka.Writer,
) {
	var (
		userJobs = make(chan UserJob)
		postJobs = make(chan PostJob)
		wg       sync.WaitGroup
		context  = context.Background()
	)

	err := godotenv.Load()
	if err != nil {
		log.Panic("error loading .env file")
	}

	// Assign workers to user and posts
	wg.Add(numOfWorkers)

	for worker := 0; worker < numOfWorkers; worker++ {
		go userWorker(context, &wg, userJobs, producer, topic)
		go postWorker(context, &wg, postJobs, producer, topic)
	}

	// Adding jobs for users
	go fetchAllUsers(userJobs)
	// Adding jobs for posts
	go fetchAllPosts(postJobs)

	wg.Wait()
}

func userWorker(ctx context.Context, wg *sync.WaitGroup, userJob <-chan UserJob, producer *kafka.Writer, topic string) {
	// job here are read only
	defer wg.Done()

	for job := range userJob {
		message, err := printUserDetail(job.User)
		if err != nil {
			log.Fatalf("Error printing user detail: %v", err)
		}

		kafkaMessage := kafka.Message{
			Topic:     topic,
			Partition: -1,
			Value:     []byte(message),
		}

		if err := producer.WriteMessages(ctx, kafkaMessage); err != nil {
			log.Fatalf("Failed to write message: %v", err)
		}
	}
}

func postWorker(ctx context.Context, wg *sync.WaitGroup, postJob <-chan PostJob, producer *kafka.Writer, topic string) {
	defer wg.Done()

	for job := range postJob {
		message, err := printPostDetail(job.Post)
		if err != nil {
			log.Fatalf("Error printing post detail: %v", err)
		}

		kafkaMessage := kafka.Message{
			Topic:     topic,
			Partition: -1,
			Value:     []byte(message),
		}

		if err := producer.WriteMessages(ctx, kafkaMessage); err != nil {
			log.Fatalf("Failed to write message: %v", err)
		}
	}
}

func fetchAllUsers(userJob chan<- UserJob) error {
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
	userResponse, err := api.FetchUsers(page)

	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	for _, user := range userResponse.Data {
		userJob <- UserJob{Page: page, User: user}
	}

	return nil
}

func printUserDetail(user entity.UserPreview) (string, error) {
	userDetail, err := api.FetchUserDetail(user.ID)

	if err != nil {
		return "", fmt.Errorf("error fetching user detail: %w", err)
	}

	message := fmt.Sprintf("User name %s %s %s %s %s \n", user.Title, user.FirstName, user.LastName, userDetail.Email, userDetail.Gender)

	return message, nil
}

func fetchAllPosts(postJob chan<- PostJob) error {
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
	postResponse, err := api.FetchPosts(page)

	if err != nil {
		return fmt.Errorf("error fetching posts: %w", err)
	}

	for _, post := range postResponse.Data {
		postJob <- PostJob{Page: page, Post: post}
	}

	return nil
}

func printPostDetail(post entity.PostPreview) (string, error) {
	var message = "Post by " + post.Owner.FirstName + " " + post.Owner.LastName + ":\n" + post.Text + "\n\nLikes " + fmt.Sprint(post.Likes) + " Tags " + fmt.Sprint(post.Tags) + "\nDate posted " + post.PublishDate.Format("2006-01-02 15:04:05")

	return message, nil
}
