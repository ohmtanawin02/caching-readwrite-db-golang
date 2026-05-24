package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type SoftDeleteHandlerCfg struct {
	ProductCommand domain.ProductApplicationCommand
}

// SoftDeleteProduct godoc
// @Summary      Soft delete product
// @Description  Set deleted_at — product จะไม่แสดงใน query ปกติ แต่ข้อมูลยังอยู่ใน DB
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  common.JSONResponse
// @Failure      400  {object}  common.JSONResponse
// @Failure      404  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Router       /products/{id} [delete]
func SoftDeleteProduct(cfg SoftDeleteHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		if err := cfg.ProductCommand.SoftDelete(c.UserContext(), uint(id)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
					constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, nil)
	}
}
