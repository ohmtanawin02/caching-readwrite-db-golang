package server

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/gofiber/swagger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"golang-fiber/config"
	_ "golang-fiber/docs"
	productRouter "golang-fiber/internal/product/interface/http"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

func NewServer(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return common.ResponseJsonWithCode(c, code, uuid.New(),
				constants.CodeInternalError,
				constants.MessageENSomethingWentWrong,
				constants.MessageTHSomethingWentWrong,
				nil)
		},
		ReadBufferSize: 4 * 1024 * 1024,
		BodyLimit:      10 * 1024 * 1024,
	})

	app.Use(cors.New())

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &logger}))
	app.Use(func(c *fiber.Ctx) error {
		return common.SetLogCtx(c, logger)
	})

	readDB, err := cfg.DBRead.Connect(cfg.Debug)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect read db")
	}
	writeDB, err := cfg.DBWrite.Connect(cfg.Debug)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to connect write db")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, _ := readDB.DB()
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy", "db": err.Error(),
			})
		}
		if err := rdb.Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy", "redis": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	validate := validator.New()

	// Routes
	api := app.Group("/api/v1")

	productRouter.NewProductRouterCfg{
		App:      api,
		ReadDB:   readDB,
		WriteDB:  writeDB,
		Redis:    rdb,
		Logger:   logger,
		Validate: validate,
	}.NewProductRouter()

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
			constants.CodeNotFound,
			constants.MessageENNotFound,
			constants.MessageTHNotFound,
			nil)
	})

	return app
}
