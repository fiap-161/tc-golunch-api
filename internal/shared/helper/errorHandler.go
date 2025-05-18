package helper

import (
	"net/http"

	appError "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "Internal Server Error"

	switch err.(type) {
	case *appError.ValidationError:
		status = http.StatusBadRequest
		message = "Validation failed"
	case *appError.UnauthorizedError:
		status = http.StatusUnauthorized
		message = "Unauthorized"
	case *appError.NotFoundError:
		status = http.StatusBadRequest
		message = "Invalid resource"
	}

	c.JSON(status, appError.ErrorDTO{
		Message:      message,
		MessageError: err.Error(),
	})
}
