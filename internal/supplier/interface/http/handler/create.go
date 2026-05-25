package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"golang-fiber/internal/supplier/application/commands"
	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type CreateHandlerCfg struct {
	SupplierCommand domain.SupplierApplicationCommand
	Validator       *validator.Validate
}

// CreateSupplier godoc
// @Summary      Create supplier
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateSupplierRequest  true  "Supplier data"
// @Success      201   {object}  common.JSONResponse{data=internal_supplier_interface_http_dto.SupplierResponse}
// @Failure      400   {object}  common.JSONResponse
// @Failure      409   {object}  common.JSONResponse
// @Failure      500   {object}  common.JSONResponse
// @Router       /suppliers [post]
func CreateSupplier(cfg CreateHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateCreateRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		supplier, err := cfg.SupplierCommand.Create(c.UserContext(), dto.ToCreateDomainInput(req))
		if err != nil {
			if errors.Is(err, commands.ErrDuplicateSupplierName) {
				return common.ResponseJsonWithCode(c, fiber.StatusConflict, uuid.New(),
					constants.CodeConflict, constants.MessageENConflict, constants.MessageTHConflict, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusCreated, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToSupplierResponse(*supplier))
	}
}
