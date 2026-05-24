package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	domain "golang-fiber/internal/user/domain"
	"golang-fiber/internal/user/interface/http/dto"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

type RegisterHandlerCfg struct {
	UserCommand domain.UserApplicationCommand
	Validator   *validator.Validate
}

// Register godoc
// @Summary      Register user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Register data"
// @Success      201   {object}  common.JSONResponse{data=dto.UserResponse}
// @Failure      400   {object}  common.JSONResponse
// @Failure      409   {object}  common.JSONResponse
// @Router       /auth/register [post]
func Register(cfg RegisterHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateRegisterRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		user, err := cfg.UserCommand.Register(c.UserContext(), dto.ToRegisterDomainInput(req))
		if err != nil {
			if errors.Is(err, domain.ErrDuplicateUsername) || errors.Is(err, domain.ErrDuplicateEmail) {
				return common.ResponseJsonWithCode(c, fiber.StatusConflict, uuid.New(),
					constants.CodeConflict, err.Error(), err.Error(), nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusCreated, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToUserResponse(user))
	}
}
