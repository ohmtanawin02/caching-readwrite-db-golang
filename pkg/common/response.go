package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type response struct {
	Status  int        `json:"status"`
	ErrorID *uuid.UUID `json:"error_id,omitempty"`
	Message message    `json:"message"`
	Data    any        `json:"data"`
}

type message struct {
	Code string `json:"code"`
	EN   string `json:"en"`
	TH   string `json:"th"`
}

type JSONResponse struct {
	Status  int        `json:"status"`
	ErrorID *uuid.UUID `json:"error_id,omitempty"`
	Message message    `json:"message"`
	Data    any        `json:"data"`
}

func ResponseJsonWithCode(c *fiber.Ctx, statusCode int, errorID uuid.UUID, code, messageEN, messageTH string, data any) error {
	r := response{
		Status: statusCode,
		Message: message{
			Code: code,
			EN:   messageEN,
			TH:   messageTH,
		},
		Data: data,
	}

	if errorID != uuid.Nil {
		r.ErrorID = &errorID
	}

	return c.Status(statusCode).JSON(r)
}
