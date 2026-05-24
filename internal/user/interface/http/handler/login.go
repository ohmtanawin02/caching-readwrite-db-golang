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

type LoginHandlerCfg struct {
	UserCommand domain.UserApplicationCommand
	Validator   *validator.Validate
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login data"
// @Success      200   {object}  common.JSONResponse{data=dto.LoginResponse}
// @Failure      400   {object}  common.JSONResponse
// @Failure      401   {object}  common.JSONResponse
// @Router       /auth/login [post]
func Login(cfg LoginHandlerCfg) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := dto.ValidateLoginRequest(c, cfg.Validator)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusBadRequest, uuid.New(),
				constants.CodeBadRequest, constants.MessageENBadRequest, constants.MessageTHBadRequest, nil)
		}

		result, err := cfg.UserCommand.Login(c.UserContext(), dto.ToLoginDomainInput(req))
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredentials) {
				return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
					constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
			}
			return common.ResponseJsonWithCode(c, fiber.StatusInternalServerError, uuid.New(),
				constants.CodeInternalError, constants.MessageENSomethingWentWrong, constants.MessageTHSomethingWentWrong, nil)
		}

		return common.ResponseJsonWithCode(c, fiber.StatusOK, uuid.Nil,
			constants.CodeOK, constants.MessageENSuccess, constants.MessageTHSuccess, dto.ToLoginResponse(result))
	}
}
