package usersHandlers

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/users"
	"github.com/yporn/sirarom-backend/modules/users/usersUsecases"
	"github.com/yporn/sirarom-backend/pkg/auth"
)

type userHandlersErrCode string

const (
	SignUpErr             userHandlersErrCode = "users-001"
	SignInErr             userHandlersErrCode = "users-002"
	RefreshPassportErr    userHandlersErrCode = "users-003"
	SignOutErr            userHandlersErrCode = "users-004"
	GenerateAdminTokenErr userHandlersErrCode = "users-005"
	UpdateUserErr         userHandlersErrCode = "users-006"
	DeleteUserErr         userHandlersErrCode = "users-007"
	FindOneUserErr        userHandlersErrCode = "users-008"
	FindUserErr           userHandlersErrCode = "users-009"
)

type IUsersHandler interface {
	FindOneUser(c *fiber.Ctx) error
	FindUser(c *fiber.Ctx) error 
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
	filesUsecase filesUsecases.IFilesUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase, filesUsecase filesUsecases.IFilesUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
		filesUsecase: filesUsecase,
	}
}

func (h *usersHandler) FindOneUser(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	user, err := h.usersUsecase.FindOneUser(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindOneUserErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, user).Res()
}

func (h *usersHandler) FindUser(c *fiber.Ctx) error {
	req := &users.UserFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindUserErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 100000
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}
	if req.Sort == "" {
		req.Sort = "DESC"
	}

	users := h.usersUsecase.FindUser(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, users).Res()
}

func (h *usersHandler) SignUp(c *fiber.Ctx) error {
	// Request body parser
	req := &users.User{
		Images:   make([]*entities.Image, 0),
		UserRole: make([]*users.UserRole, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(SignUpErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignInErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(
		auth.Admin,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(GenerateAdminTokenErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *usersHandler) UpdateUser(c *fiber.Ctx) error {
	userIdStr := strings.Trim(c.Params("user_id"), " ")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateUserErr),
			err.Error(),
		).Res()
	}

	req := &users.User{
		Images:   make([]*entities.Image, 0),
		UserRole: make([]*users.UserRole, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateUserErr),
			err.Error(),
		).Res()
	}
	req.Id = userId

	user, err := h.usersUsecase.UpdateUser(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(UpdateUserErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, user).Res()
}

func (h *usersHandler) DeleteUser(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	// Retrieve the user by ID
	user, err := h.usersUsecase.FindOneUser(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteUserErr),
			err.Error(),
		).Res()
	}

	// Construct requests to delete associated files
	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, img := range user.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("users/%s", path.Base(img.Url)),
		})
	}

	// Delete associated files from storage
	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteUserErr),
			err.Error(),
		).Res()
	}

	// Delete the user
	if err := h.usersUsecase.DeleteUser(userId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteUserErr),
			err.Error(),
		).Res()
	}

	// Return success response
	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
