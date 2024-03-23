package houseModelsPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels"
)

type IUpdateHouseModelBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateHouseModel() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateHouseModelBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *houseModels.HouseModel
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateHouseModelBuilder(db *sqlx.DB, req *houseModels.HouseModel, filesUsecases filesUsecases.IFilesUsecase) IUpdateHouseModelBuilder {
	return &updateHouseModelBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateHouseModelEngineer struct {
	builder IUpdateHouseModelBuilder
}

func (b *updateHouseModelBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateHouseModelBuilder) initQuery(){
	b.query += `
	UPDATE "house_models" SET`
}

func (b *updateHouseModelBuilder) updateQuery() {
	setStatements := []string{}

	if b.req.Name != "" {
		b.values = append(b.values, b.req.Name)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"name" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkVideo != "" {
		b.values = append(b.values, b.req.LinkVideo)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_video" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkVirtualTour != "" {
		b.values = append(b.values, b.req.LinkVirtualTour)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_virtual_tour" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	if b.req.Index != 0 {
		b.values = append(b.values, b.req.Index)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"index" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateHouseModelBuilder) insertImages() error {
	query := `
	INSERT INTO "house_model_images" (
		"filename",
		"url",
		"house_model_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}
func (b *updateHouseModelBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "house_model_images"
	WHERE "house_model_id" = $1;`

	images := make([]*entities.Image, 0)
	if err := b.db.Select(
		&images,
		query,
		b.req.Id,
	); err != nil {
		return make([]*entities.Image, 0)
	}
	return images
}
func (b *updateHouseModelBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "house_model_images"
	WHERE "house_model_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/house_models/%s", img.FileName),
			})
		}
		b.filesUsecases.DeleteFileOnStorage(deleteFileReq)
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete images failed: %v", err)
	}
	return nil
}

func (b *updateHouseModelBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateHouseModelBuilder) updateHouseModel() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update house model failed: %v", err)
	}
	return nil
}

func (b *updateHouseModelBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateHouseModelBuilder) getValues() []any         { return b.values }
func (b *updateHouseModelBuilder) getQuery() string         { return b.query }
func (b *updateHouseModelBuilder) setQuery(query string)    { b.query = query }
func (b *updateHouseModelBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateHouseModelBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateHouseModelEngineer(b IUpdateHouseModelBuilder) *updateHouseModelEngineer {
	return &updateHouseModelEngineer{builder: b}
}

func (en *updateHouseModelEngineer) sumQueryFields() {
	en.builder.updateQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
}

func (en *updateHouseModelEngineer) UpdateHouseModel() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}
