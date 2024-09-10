package response

import "github.com/gofiber/fiber/v2"

type ErrorMessage struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      bool   `json:"error"`
}

type Response struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
}

func HandleError(ctx *fiber.Ctx, err error, message string, statusCode int) error {
	if message == "" {
		message = err.Error()
	}
	return ctx.Status(statusCode).JSON(ErrorMessage{
		Error:      true,
		Message:    message,
		StatusCode: statusCode,
	})
}

func HandleSuccess(ctx *fiber.Ctx, message string, data interface{}, statusCode int) error {
	if data != nil {
		return ctx.Status(statusCode).JSON(Response{
			Message:    message,
			Data:       data,
			StatusCode: statusCode,
		})
	}
	return ctx.Status(statusCode).JSON(fiber.Map{
		"message":    message,
		"statusCode": statusCode,
	})
}
