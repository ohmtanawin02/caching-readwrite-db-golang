package commands

import (
	"context"
	"errors"
	"time"

	domain "golang-fiber/internal/user/domain"
	"golang-fiber/internal/user/domain/entity"
	"golang-fiber/pkg/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserCommand struct {
	queryRepo domain.UserQueryRepository
	cmdRepo   domain.UserCommandRepository
	jwtSecret string
	jwtTTL    time.Duration
}

type UserCommandCfg struct {
	QueryRepo domain.UserQueryRepository
	CmdRepo   domain.UserCommandRepository
	JWTSecret string
	JWTTTL    time.Duration
}

func NewUserCommand(cfg UserCommandCfg) domain.UserApplicationCommand {
	return &UserCommand{
		queryRepo: cfg.QueryRepo,
		cmdRepo:   cfg.CmdRepo,
		jwtSecret: cfg.JWTSecret,
		jwtTTL:    cfg.JWTTTL,
	}
}

func (c *UserCommand) Register(ctx context.Context, input domain.RegisterInput) (*entity.User, error) {
	// duplicate checks go to writeDB (via cmdRepo) to avoid replication lag
	if _, err := c.cmdRepo.FindByUsername(ctx, input.Username); err == nil {
		return nil, domain.ErrDuplicateUsername
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := c.cmdRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, domain.ErrDuplicateEmail
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashed),
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Phone:     input.Phone,
	}

	if err := c.cmdRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = "" // ไม่ส่ง hash กลับไป
	return user, nil
}

func (c *UserCommand) Login(ctx context.Context, input domain.LoginInput) (*domain.LoginResult, error) {
	user, err := c.queryRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, c.jwtSecret, c.jwtTTL)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &domain.LoginResult{Token: token, User: user}, nil
}
