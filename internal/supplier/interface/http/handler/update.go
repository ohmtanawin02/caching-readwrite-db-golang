package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"golang-fiber/internal/supplier/application/commands"
	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type UpdateHandlerCfg struct {
	SupplierCommand domain.SupplierApplicationCommand
	Validator       *validator.Validate
}

// UpdateSupplier godoc
// @Summary      Update supplier
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                       true  "Supplier ID"
// @Param        body  body      dto.UpdateSupplierRequest  true  "Supplier data"
// @Success      200   {object}  common.JSONResponse{data=internal_supplier_interface_http_dto.SupplierResponse}
// @Failure      400   {object}  common.JSONResponse
// @Failure      404   {object}  common.JSONResponse
// @Failure      409   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /suppliers/{id} [put]
func UpdateSupplier(cfg UpdateHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		req, err := dto.ValidateUpdateRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		supplier, err := cfg.SupplierCommand.Update(c.UserContext(), uint(id), dto.ToUpdateDomainInput(req))
		if err != nil {
			if errors.Is(err, commands.ErrDuplicateSupplierName) {
				return common.ResponseJsonWithCode(c, fiber.StatusConflict, uuid.New(),
					constants.CodeConflict, constants.MessageENConflict, constants.MessageTHConflict, nil)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return common.ResponseJsonWithCode(c, fiber.StatusNotFound, uuid.New(),
					constants.CodeNotFound, constants.MessageENNotFound, constants.MessageTHNotFound, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToSupplierResponse(*supplier))
	}
}
