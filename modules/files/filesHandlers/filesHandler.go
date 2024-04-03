package filesHandlers

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type filesHandlersErrCode string

const (
	uploadErr filesHandlersErrCode = "files-001"
	deleteErr filesHandlersErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
		"webp": "webp",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		if ext == "webp" {
			req = append(req, &files.FileReq{
				File:        file,
				Destination: destination + "/" + file.Filename,
				FileName:    file.Filename,
				Extension:   ext,
			})
		} else {
			// Convert other image formats to webp
			webPFileName := utils.RandFileName("webp")
			webPFilePath := filepath.Join(destination, webPFileName)

			// Convert to WebP directly without saving to temporary location
			if err := convertToWebP(file, webPFilePath); err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(uploadErr),
					err.Error(),
				).Res()
			}

			// Add the WebP file to the list of files to be uploaded
			req = append(req, &files.FileReq{
				File:        file,
				Destination: destination + "/" + webPFileName,
				FileName:    webPFileName,
				Extension:   "webp",
			})
		}
	}

	res, err := h.filesUsecase.UploadToStorage(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.filesUsecase.DeleteFileOnStorage(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func convertToWebP(file *multipart.FileHeader, outputPath string) error {
	//open file
	inputFile, err := file.Open()
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var img image.Image
	fileName := file.Filename

	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".png":
		img, err = png.Decode(inputFile)
		if err != nil {
			return err
		}
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(inputFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported image format")
	}

	// Create the directory if it doesn't exist
	err = os.MkdirAll("./assets/images/convert/"+filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	outputFile, err := os.Create("./assets/images/convert/" + outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = webp.Encode(outputFile, img, nil)
	if err != nil {
		return err
	}

	// Delete the directory after encoding the file to WebP
	if err := os.RemoveAll("./assets/images/convert/"); err != nil {
		return err
	}

	return nil
}
