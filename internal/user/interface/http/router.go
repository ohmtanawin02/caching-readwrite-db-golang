package userRouter

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"golang-fiber/internal/user/application/commands"
	"golang-fiber/internal/user/infrastructure/repository"
	"golang-fiber/internal/user/interface/http/handler"
)

type NewUserRouterCfg struct {
	App       fiber.Router
	ReadDB    *gorm.DB
	WriteDB   *gorm.DB
	Logger    zerolog.Logger
	Validate  *validator.Validate
	JWTSecret string
	JWTTTL    time.Duration
}

func (cfg NewUserRouterCfg) NewUserRouter() {
	queryRepo := repository.NewUserQueryRepository(repository.UserQueryRepositoryCfg{
		DB:     cfg.ReadDB,
		Logger: cfg.Logger,
	})
	cmdRepo := repository.NewUserCommandRepository(repository.UserCommandRepositoryCfg{
		DB:     cfg.WriteDB,
		Logger: cfg.Logger,
	})

	command := commands.NewUserCommand(commands.UserCommandCfg{
		QueryRepo: queryRepo,
		CmdRepo:   cmdRepo,
		JWTSecret: cfg.JWTSecret,
		JWTTTL:    cfg.JWTTTL,
	})

	g := cfg.App.Group("/auth")
	g.Post("/register", handler.Register(handler.RegisterHandlerCfg{UserCommand: command, Validator: cfg.Validate}))
	g.Post("/login", handler.Login(handler.LoginHandlerCfg{UserCommand: command, Validator: cfg.Validate}))
}
