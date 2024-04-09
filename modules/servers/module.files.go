package servers

import (
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/files/filesHandlers"
)

type IFilesModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFilesHandler
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFilesHandler
}

func (m *moduleFactory) FilesModule() IFilesModule {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.r.Group("/files")
	router.Post("/upload", f.mid.JwtAuth(), f.handler.UploadFiles)
	router.Patch("/delete", f.mid.JwtAuth(), f.handler.DeleteFile)

}

func (f *filesModule) Usecase() filesUsecases.IFilesUsecase { return f.usecase }
func (f *filesModule) Handler() filesHandlers.IFilesHandler { return f.handler }
