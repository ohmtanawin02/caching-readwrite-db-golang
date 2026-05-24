package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type FindByIDHandlerCfg struct {
	ProductQuery domain.ProductApplicationQuery
}

// FindProductByID godoc
// @Summary      Get product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  common.JSONResponse{data=dto.ProductResponse}
// @Failure      400  {object}  common.JSONResponse
// @Failure      404  {object}  common.JSONResponse
// @Router       /products/{id} [get]
func FindProductByID(cfg FindByIDHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		product, err := cfg.ProductQuery.FindByID(c.UserContext(), uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
					constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToProductResponse(*product))
	}
}
