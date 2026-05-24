package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/application/commands"
	"golang-fiber/internal/product/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type CreateHandlerCfg struct {
	ProductCommand domain.ProductApplicationCommand
	Validator      *validator.Validate
}

// CreateProduct godoc
// @Summary      Create product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateProductRequest  true  "Product data"
// @Success      201   {object}  common.JSONResponse{data=dto.ProductResponse}
// @Failure      400   {object}  common.JSONResponse
// @Failure      409   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /products [post]
func CreateProduct(cfg CreateHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateCreateRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		product, err := cfg.ProductCommand.Create(c.UserContext(), dto.ToCreateDomainInput(req))
		if err != nil {
			if errors.Is(err, commands.ErrDuplicateProductName) {
				return common.ResponseJsonWithCode(c, fiber.StatusConflict, uuid.New(),
					constants.CodeConflict, constants.MessageENConflict, constants.MessageTHConflict, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusCreated, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToProductResponse(*product))
	}
}
