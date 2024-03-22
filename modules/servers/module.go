package servers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yporn/sirarom-backend/modules/activities/activitiesHandlers"
	"github.com/yporn/sirarom-backend/modules/activities/activitiesRepositories"
	"github.com/yporn/sirarom-backend/modules/activities/activitiesUsecases"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoHandlers"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoRepositories"
	"github.com/yporn/sirarom-backend/modules/appinfo/appinfoUsecases"
	"github.com/yporn/sirarom-backend/modules/banners/bannersHandlers"
	"github.com/yporn/sirarom-backend/modules/banners/bannersRepositories"
	"github.com/yporn/sirarom-backend/modules/banners/bannersUsecases"
	"github.com/yporn/sirarom-backend/modules/general/generalHandlers"
	"github.com/yporn/sirarom-backend/modules/general/generalRepositories"
	"github.com/yporn/sirarom-backend/modules/general/generalUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsHandlers"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsRepositories"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsUsecases"
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
	"github.com/yporn/sirarom-backend/modules/projects/projectsHandlers"
	"github.com/yporn/sirarom-backend/modules/projects/projectsRepositories"
	"github.com/yporn/sirarom-backend/modules/projects/projectsUsecases"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsHandlers"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsRepositories"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsUsecases"
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
	BannerModule()
	ActivityModule()
	ProjectModule()
	HouseModelModule()
	PromotionModule()
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

	router.Post("/signup", m.mid.JwtAuth(), handler.SignUp)
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

	router.Get("/:job_id", handler.FindOneJob)
	router.Get("/", handler.FindJob)
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

	router.Get("/:interest_id", handler.FindOneInterest)
	router.Get("/", handler.FindInterest)

	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddInterest)
	router.Patch("/update/:interest_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateInterest)
	router.Delete("/:interest_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteInterest)
}

func (m *moduleFactory) BannerModule() {
	repository := bannersRepositories.BannersRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := bannersUsecases.BannersUsecase(repository)
	handler := bannersHandlers.BannersHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/banners")

	router.Get("/:banner_id", handler.FindOneBanner)
	router.Get("/", handler.FindBanner)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddBanner)
	router.Patch("/update/:banner_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateBanner)
	router.Delete("/:banner_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteBanner)
}

func (m *moduleFactory) ActivityModule() {
	repository := activitiesRepositories.ActivitiesRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := activitiesUsecases.ActivitiesUsecase(repository)
	handler := activitiesHandlers.ActivitiesHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/activities")

	router.Get("/:activity_id", handler.FindOneActivity)
	router.Get("/", handler.FindActivity)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddActivity)
	router.Patch("/update/:activity_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateActivity)
	router.Delete("/:activity_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteActivity)
}

func (m *moduleFactory) ProjectModule() {
	repository := projectsRepositories.ProjectsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := projectsUsecases.ProjectsUsecase(repository)
	handler := projectsHandlers.ProjectsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/projects")

	router.Get("/:project_id", handler.FindOneProject)
	router.Get("/", handler.FindProject)
	router.Get("/:project_id/house_models", handler.FindProjectHouseModel)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProject)
	router.Patch("/update/:project_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProject)
	router.Delete("/:project_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProject)
}

func (m *moduleFactory) HouseModelModule() {
	repository := houseModelsRepositories.HouseModelsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := houseModelsUsecases.HouseModelsUsecases(repository)
	handler := houseModelsHandlers.HouseModelsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/house_models")

	router.Get("/:house_model_id", handler.FindOneHouseModel)
	router.Get("/projects/:project_id", handler.FindHouseModel)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddHouseModel)
	router.Patch("/update/:house_model_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateHouseModel)
	router.Delete("/:house_model_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteHouseModel)
}

func (m *moduleFactory) PromotionModule() {
	repository := promotionsRepositories.PromotionsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := promotionsUsecases.PromotionsUsecase(repository)
	handler := promotionsHandlers.PromotionsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	router := m.r.Group("/promotions")

	router.Get("/", handler.FindPromotion)
	router.Get("/:promotion_id", handler.FindOnePromotion)
	router.Post("/create", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddPromotion)
	router.Patch("/update/:promotion_id", m.mid.JwtAuth(), handler.UpdatePromotion)
	// router.Delete("/:house_model_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteHouseModel)
}