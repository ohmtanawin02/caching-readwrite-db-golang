package hpRouter

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"golang-fiber/internal/product/application/commands"
	"golang-fiber/internal/product/application/queries"
	"golang-fiber/internal/product/infrastructure/repository"
	"golang-fiber/internal/product/interface/http/handler"
	"golang-fiber/pkg/cache"
)

type NewProductRouterCfg struct {
	App      fiber.Router
	ReadDB   *gorm.DB
	WriteDB  *gorm.DB
	Redis    *redis.Client
	Logger   zerolog.Logger
	Validate *validator.Validate
}

func (cfg NewProductRouterCfg) NewProductRouter() {
	queryRepo := repository.NewProductQueryRepository(repository.ProductQueryRepositoryCfg{
		DB:     cfg.ReadDB,
		Logger: cfg.Logger,
	})
	cmdRepo := repository.NewProductCommandRepository(repository.ProductCommandRepositoryCfg{
		DB:     cfg.WriteDB,
		Logger: cfg.Logger,
	})

	redisCache := cache.NewRedisCache(cfg.Redis)

	// use cases
	productQuery := queries.NewProductQuery(queries.ProductQueryCfg{
		Repo:  queryRepo,
		Cache: redisCache,
	})
	productCommand := commands.NewProductCommand(commands.ProductCommandCfg{
		Repo:  cmdRepo,
		Cache: redisCache,
	})

	// register routes
	g := cfg.App.Group("/products")
	g.Get("/", handler.FindAllProducts(handler.FindAllHandlerCfg{ProductQuery: productQuery, Validator: cfg.Validate}))
	g.Get("/:id", handler.FindProductByID(handler.FindByIDHandlerCfg{ProductQuery: productQuery}))
	g.Post("/", handler.CreateProduct(handler.CreateHandlerCfg{ProductCommand: productCommand, Validator: cfg.Validate}))
	g.Put("/:id", handler.UpdateProduct(handler.UpdateHandlerCfg{ProductCommand: productCommand, Validator: cfg.Validate}))
	g.Delete("/:id", handler.SoftDeleteProduct(handler.SoftDeleteHandlerCfg{ProductCommand: productCommand}))
}
