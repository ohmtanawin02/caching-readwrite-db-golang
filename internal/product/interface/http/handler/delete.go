package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type DeleteHandlerCfg struct {
	ProductCommand domain.ProductApplicationCommand
}

// DeleteProduct godoc
// @Summary      Delete product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      400  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Router       /products/{id} [delete]
func DeleteProduct(cfg DeleteHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		if err := cfg.ProductCommand.Delete(c.UserContext(), uint(id)); err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, nil)
	}
}
