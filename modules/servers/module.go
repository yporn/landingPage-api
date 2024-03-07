package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoHandlers"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoRepositories"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoUsecases"
	"github.com/yporn/sirarom-backend/modules/general/generalHandlers"
	"github.com/yporn/sirarom-backend/modules/general/generalRepositories"
	"github.com/yporn/sirarom-backend/modules/general/generalUsecases"
	"github.com/yporn/sirarom-backend/modules/interests/interestsHandlers"
	"github.com/yporn/sirarom-backend/modules/interests/interestsRepositories"
	"github.com/yporn/sirarom-backend/modules/interests/interestsUsecases"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsHandlers"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsRepositories"
	"github.com/yporn/sirarom-backend/modules/jobs/jobsUsecases"
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
	AppinfoModule()
	JobModule()
	FilesModule() IFilesModule
	GeneralModule()
	InterestModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUp)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", m.mid.JwtAuth(), handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)

	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
}

func (m *moduleFactory) JobModule() {
	repository := jobsRepositories.JobsRepository(m.s.db, m.s.cfg)
	usecase := jobsUsecases.JobsUsecase(repository)
	handler := jobsHandlers.JobsHandler(m.s.cfg, usecase)

	router := m.r.Group("/jobs")

	router.Get("/:job_id", m.mid.JwtAuth(), handler.FindOneJob)
	router.Get("/", m.mid.JwtAuth(), handler.FindJob)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddJob)
	router.Patch("/update/:job_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateJob)
	router.Delete("/:job_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteJob)
}

func (m *moduleFactory) GeneralModule() {
	repository := generalRepositories.GeneralRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := generalUsecases.GeneralUsecase(repository)
	handler := generalHandlers.GeneralHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/data_setting")

	router.Get("/:general_id", m.mid.JwtAuth(), handler.FindOneGeneral)
	router.Patch("/update/:general_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateGeneral)
}

func (m *moduleFactory) InterestModule() {
	repository := interestsRepositories.InterestsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := interestsUsecases.InterestsUsecase(repository)
	handler := interestsHandlers.InterestsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/interests")

	router.Get("/:interest_id", m.mid.ApiKeyAuth(), handler.FindOneInterest)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddInterest)
	// router.Patch("/update/:job_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateJob)
	// router.Delete("/:job_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteJob)
}
