package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	domain "golang-fiber/internal/supplier/domain"
	"golang-fiber/internal/supplier/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type FindAllHandlerCfg struct {
	SupplierQuery domain.SupplierApplicationQuery
	Validator     *validator.Validate
}

// FindAllSuppliers godoc
// @Summary      List suppliers
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page        query  int    false  "Page number"          default(1)
// @Param        limit       query  int    false  "Items per page"       default(20)
// @Param        sort_by     query  string false  "Sort field"           Enums(id,name,created_at) default(id)
// @Param        sort_order  query  string false  "Sort direction"       Enums(asc,desc) default(asc)
// @Success      200  {object}  common.JSONResponse{data=dto.SupplierListResponse}
// @Failure      401  {object}  common.JSONResponse
// @Failure      500  {object}  common.JSONResponse
// @Router       /suppliers [get]
func FindAllSuppliers(cfg FindAllHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateFindAllRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		domainReq := dto.ToFindAllDomainRequest(req)
		result, err := cfg.SupplierQuery.FindAll(c.UserContext(), domainReq)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess,
			dto.ToSupplierListResponse(result, domainReq))
	}
}
