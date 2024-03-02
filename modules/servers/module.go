package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/modules/middlewares/middlewaresHandlers"
	"github.com/yporn/sirarom-backend/modules/middlewares/middlewaresRepositories"
	"github.com/yporn/sirarom-backend/modules/middlewares/middlewaresUsecases"
	"github.com/yporn/sirarom-backend/modules/monitor/monitorHandlers"
	"github.com/yporn/sirarom-backend/modules/users/usersHandlers"
	"github.com/yporn/sirarom-backend/modules/users/usersRepositories"
	"github.com/yporn/sirarom-backend/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	// route
	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUp)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)

	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}
