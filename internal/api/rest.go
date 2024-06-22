package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CRAYON-2024/worker/bootstrap"
	"github.com/CRAYON-2024/worker/internal/handler"
	"github.com/CRAYON-2024/worker/internal/repository"
	"github.com/CRAYON-2024/worker/internal/usecase"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func APIRouter(container *bootstrap.Container) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var (
		userRepo = repository.NewUserRepository(container.GetDBRead(), container.GetTracer())
		userUC   = usecase.NewUserUseCase(usecase.UserUsecaseParam{
			UserUseCase: usecase.UserUseCase{
				Repository: userRepo,
				Producer:   container.GetKafkaProducer(),
				Topic:      viper.GetString("kafka.topic.handson"),
				Trace:      container.GetTracer(),
			},
		})
		userHandler = handler.NewUserHandler(userUC, container.GetTracer())
	)

	e.GET("/user", userHandler.GetUsers)

	go func() {
		if err := e.Start(viper.GetString("server.ip") + ":" + viper.GetString("server.port")); err != nil {
			log.Fatalln("failed to run echo router", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	logrus.Infof("shutting down the server %v", sig)
	log.Println("closing dependency", container.Terminate())

	if err := e.Shutdown(container.Ctx); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(fmt.Sprintf("failed to gracefully shut down the server %s", err))
	}
}
