package common

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type loggerKey struct{}

type logCtxData struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
}

func SetLogCtx(c *fiber.Ctx, logger zerolog.Logger) error {
	ctx := logCtxData{IP: c.IP()}

	if username, ok := c.Locals("username").(string); ok {
		ctx.Username = username
	}

	log := logger.
		With().
		Str("request_id", uuid.NewString()).
		Any("context", ctx).
		Logger()

	c.SetUserContext(context.WithValue(c.Context(), loggerKey{}, log))
	return c.Next()
}

func logFromCtx(ctx context.Context) zerolog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(zerolog.Logger); ok {
		return logger
	}
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func NewAppLogger(ctx context.Context, source string) zerolog.Logger {
	return logFromCtx(ctx).With().Str("source", source).Type("type", "APP").Logger()
}

func NewRepoLogger(ctx context.Context, source string) zerolog.Logger {
	return logFromCtx(ctx).With().Str("source", source).Type("type", "DATA").Logger()
}
