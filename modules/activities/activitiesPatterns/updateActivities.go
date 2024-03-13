package activitiesPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/activities"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
)

type IUpdateActivityBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateActivity() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateActivityBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *activities.Activity
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateActivityBuilder(db *sqlx.DB, req *activities.Activity, filesUsecases filesUsecases.IFilesUsecase) IUpdateActivityBuilder {
	return &updateActivityBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateActivityEngineer struct {
	builder IUpdateActivityBuilder
}

func (b *updateActivityBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateActivityBuilder) initQuery() {
	b.query += `
	UPDATE "activities" SET`
}

func (b *updateActivityBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Index != 0 {
		b.values = append(b.values, b.req.Index)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"index" = $%d`, b.lastStackIndex))
	}

	if b.req.Heading != "" {
		b.values = append(b.values, b.req.Heading)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"heading" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.StartDate != "" {
		b.values = append(b.values, b.req.StartDate)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"start_date" = $%d`, b.lastStackIndex))
	}

	if b.req.EndDate != "" {
		b.values = append(b.values, b.req.EndDate)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"end_date" = $%d`, b.lastStackIndex))
	}

	if b.req.VideoLink != "" {
		b.values = append(b.values, b.req.VideoLink)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"video_link" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateActivityBuilder) insertImages() error {
	query := `
	INSERT INTO "activities_images" (
		"filename",
		"url",
		"activity_id"
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

func (b *updateActivityBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "activities_images"
	WHERE "activity_id" = $1;`

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

func (b *updateActivityBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "activities_images"
	WHERE "activity_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/activity/%s", img.FileName),
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

func (b *updateActivityBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateActivityBuilder) updateActivity() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update activity failed: %v", err)
	}
	return nil
}

func (b *updateActivityBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateActivityBuilder) getValues() []any         { return b.values }
func (b *updateActivityBuilder) getQuery() string         { return b.query }
func (b *updateActivityBuilder) setQuery(query string)    { b.query = query }
func (b *updateActivityBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateActivityBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateActivityEngineer(b IUpdateActivityBuilder) *updateActivityEngineer {
	return &updateActivityEngineer{builder: b}
}

func (en *updateActivityEngineer) sumQueryFields() {
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

func (en *updateActivityEngineer) UpdateActivity() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Update banner
	if err := en.builder.updateActivity(); err != nil {
		return err
	}

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