package handler

import (
	"fmt"
	"net/http"

	"github.com/CRAYON-2024/worker/internal/usecase"
	"github.com/labstack/echo"
	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	UserUseCase *usecase.UserUseCase
	Trace trace.Tracer
}

func NewUserHandler(userUseCase *usecase.UserUseCase, trace trace.Tracer) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
		Trace: trace,
	}
}


func (uh *UserHandler) GetUsers(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	users, err := uh.UserUseCase.GetUser(ctx)

	if err != nil {
		fmt.Println("failed to get users", err)
		return c.JSON(500, map[string]interface{}{
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, users)
}