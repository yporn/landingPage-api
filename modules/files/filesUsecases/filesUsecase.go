package filesUsecases

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	
	"strings"
	"time"

	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/files"
)

type IFilesUsecase interface {
	UploadToStorage(req []*files.FileReq) ([]*files.FileRes, error)
	DeleteFileOnStorage(req []*files.DeleteFileReq) error
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

type filesPub struct {
	destination string
	file        *files.FileRes
}

func (u *filesUsecase) uploadToStorageWorker(ctx context.Context, jobs <-chan *files.FileReq, results chan<- *files.FileRes, errs chan<- error) {
	for job := range jobs {
		container, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}
		b, err := ioutil.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}

		// Upload an object to storage
		dest := fmt.Sprintf("assets/images/%s", job.Destination)
		if err := os.WriteFile(dest, b, 0777); err != nil {
			if err := os.MkdirAll("./assets/images/"+strings.Replace(job.Destination, job.FileName, "", 1), 0777); err != nil {
				errs <- fmt.Errorf("mkdir \"./assets/images/%s\" failed: %v", err, job.Destination)
				return
			}
			if err := os.WriteFile(dest, b, 0777); err != nil {
				errs <- fmt.Errorf("write file failed: %v", err)
				return
			}
		}
		// fmt.Println("http://" + u.cfg.App().Host() + ":" + strconv.Itoa(u.cfg.App().Port()) + "/" + dest)
		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("%s/%s", u.cfg.App().AppUrl(), dest),
			},
			destination: dest,
		}

		errs <- nil
		results <- newFile.file
	}
}

func (u *filesUsecase) UploadToStorage(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *files.FileReq, len(req))
	resultsCh := make(chan *files.FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadToStorageWorker(ctx, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}
	return res, nil
}

func (u *filesUsecase) deleteFromStorageFileWorkers(ctx context.Context, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		if err := os.Remove("./assets/images/" + job.Destination); err != nil {
			errs <- fmt.Errorf("remove file: %s failed: %v", job.Destination, err)
			return
		}
		errs <- nil
	}
}

func (u *filesUsecase) DeleteFileOnStorage(req []*files.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.deleteFromStorageFileWorkers(ctx, jobsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		return err
	}
	return nil
}
