package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/CRAYON-2024/worker/internal/entity"
	"github.com/CRAYON-2024/worker/internal/repository"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type UserUseCase struct {
	Topic      string
	Repository *repository.UserRepository
	Producer   *kafka.Writer
	Trace      trace.Tracer
}

type UserUsecaseParam struct {
	UserUseCase
}

func NewUserUseCase(userUsecaseParam UserUsecaseParam) *UserUseCase {
	return &UserUseCase{
		Topic:      userUsecaseParam.Topic,
		Repository: userUsecaseParam.Repository,
		Producer:   userUsecaseParam.Producer,
		Trace:      userUsecaseParam.Trace,
	}
}

func (uc *UserUseCase) GetUser(ctx context.Context) ([]entity.CustomUser, error) {
	log.Println("user usecase")

	ctx, span := uc.Trace.Start(ctx, "use case get user")
	defer span.End()

	users, err := uc.Repository.GetUsers(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.CustomUser{}, nil
		}
		return nil, err
	}

	var res = []entity.CustomUser{}
	for _, val := range users {
		res = append(res, entity.CustomUser{
			ID:   val.ID,
			Name: val.Name,
		})
	}

	byt, err := json.Marshal(res)
	if err != nil {
		log.Println("failed to marshal obj", err)
		return nil, err
	}

	log.Println(uc.Topic)
	message := kafka.Message{
		Partition: -1,
		Value:     []byte(string(byt)),
	}

	if err := uc.Producer.WriteMessages(ctx, message); err != nil {
		log.Fatalf("Failed to write message: %v", err)
	} else {
		log.Println("Message sent successfully")
	}

	return res, nil
}
