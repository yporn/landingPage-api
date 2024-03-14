package bannersPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/banners"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
)

type IUpdateBannerBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateBanner() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateBannerBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *banners.Banner
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateBannerBuilder(db *sqlx.DB, req *banners.Banner, filesUsecases filesUsecases.IFilesUsecase) IUpdateBannerBuilder {
	return &updateBannerBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateBannerEngineer struct {
	builder IUpdateBannerBuilder
}

func (b *updateBannerBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateBannerBuilder) initQuery() {
	b.query += `
	UPDATE "banners" SET`
}

func (b *updateBannerBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Index != 0 {
		b.values = append(b.values, b.req.Index)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"index" = $%d`, b.lastStackIndex))
	}

	if b.req.Delay != 0 {
		b.values = append(b.values, b.req.Delay)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"delay" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateBannerBuilder) insertImages() error {
	query := `
	INSERT INTO "banner_images" (
		"filename",
		"url",
		"banner_id"
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

func (b *updateBannerBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "banner_images"
	WHERE "banner_id" = $1;`

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

func (b *updateBannerBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "banner_images"
	WHERE "banner_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/banners/%s", img.FileName),
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

func (b *updateBannerBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateBannerBuilder) updateBanner() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update banner failed: %v", err)
	}
	return nil
}

func (b *updateBannerBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateBannerBuilder) getValues() []any         { return b.values }
func (b *updateBannerBuilder) getQuery() string         { return b.query }
func (b *updateBannerBuilder) setQuery(query string)    { b.query = query }
func (b *updateBannerBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateBannerBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateBannerEngineer(b IUpdateBannerBuilder) *updateBannerEngineer {
	return &updateBannerEngineer{builder: b}
}

func (en *updateBannerEngineer) sumQueryFields() {
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

func (en *updateBannerEngineer) UpdateBanner() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Update banner
	if err := en.builder.updateBanner(); err != nil {
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