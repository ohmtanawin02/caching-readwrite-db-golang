package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	domain "golang-fiber/internal/product/domain"
	"golang-fiber/internal/product/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type FindAllHandlerCfg struct {
	ProductQuery domain.ProductApplicationQuery
	Validator    *validator.Validate
}

// FindAllProducts godoc
// @Summary      List products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page        query  int    false  "Page number"          default(1)
// @Param        limit       query  int    false  "Items per page"       default(20)
// @Param        supplier_id query  int    false  "Filter by supplier"
// @Param        sort_by     query  string false  "Sort field"           Enums(id,name,price,stock,created_at) default(id)
// @Param        sort_order  query  string false  "Sort direction"       Enums(asc,desc) default(asc)
// @Success      200  {object}  common.JSONResponse{data=[]dto.ProductResponse}
// @Failure      500  {object}  common.JSONResponse
// @Router       /products [get]
func FindAllProducts(cfg FindAllHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateFindAllRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		products, err := cfg.ProductQuery.FindAll(c.UserContext(), dto.ToFindAllDomainRequest(req))
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToProductResponses(products))
	}
}
