package hpRouter

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"golang-fiber/internal/supplier/application/commands"
	"golang-fiber/internal/supplier/application/queries"
	"golang-fiber/internal/supplier/infrastructure/repository"
	"golang-fiber/internal/supplier/interface/http/handler"
	"golang-fiber/pkg/cache"
)

type NewSupplierRouterCfg struct {
	App      fiber.Router
	ReadDB   *gorm.DB
	WriteDB  *gorm.DB
	Redis    *redis.Client
	Logger   zerolog.Logger
	Validate *validator.Validate
}

func (cfg NewSupplierRouterCfg) NewSupplierRouter() {
	queryRepo := repository.NewSupplierQueryRepository(repository.SupplierQueryRepositoryCfg{
		DB:     cfg.ReadDB,
		Logger: cfg.Logger,
	})
	cmdRepo := repository.NewSupplierCommandRepository(repository.SupplierCommandRepositoryCfg{
		DB:     cfg.WriteDB,
		Logger: cfg.Logger,
	})

	redisCache := cache.NewRedisCache(cfg.Redis)

	// use cases
	supplierQuery := queries.NewSupplierQuery(queries.SupplierQueryCfg{
		Repo:  queryRepo,
		Cache: redisCache,
	})
	supplierCommand := commands.NewSupplierCommand(commands.SupplierCommandCfg{
		Repo:  cmdRepo,
		Cache: redisCache,
	})

	// register routes
	g := cfg.App.Group("/suppliers")
	g.Get("/", handler.FindAllSuppliers(handler.FindAllHandlerCfg{SupplierQuery: supplierQuery, Validator: cfg.Validate}))
	g.Get("/:id", handler.FindSupplierByID(handler.FindByIDHandlerCfg{SupplierQuery: supplierQuery}))
	g.Post("/", handler.CreateSupplier(handler.CreateHandlerCfg{SupplierCommand: supplierCommand, Validator: cfg.Validate}))
	g.Put("/:id", handler.UpdateSupplier(handler.UpdateHandlerCfg{SupplierCommand: supplierCommand, Validator: cfg.Validate}))
	g.Delete("/:id", handler.SoftDeleteSupplier(handler.SoftDeleteHandlerCfg{SupplierCommand: supplierCommand}))
}
