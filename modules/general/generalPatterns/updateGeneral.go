package generalPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/general"
)

type IUpdateGeneralBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	closeQuery()
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	commit() error
	updateGeneral() error
	insertImages() error 
	getOldImages() []*entities.Image
	getImagesLen() int
	deleteOldImages() error
}

type updateGeneralBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *general.General
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateGeneralBuilder(db *sqlx.DB, req *general.General, filesUsecases filesUsecases.IFilesUsecase) IUpdateGeneralBuilder {
	return &updateGeneralBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}


type updateGeneralEngineer struct {
	builder IUpdateGeneralBuilder
}


func (b *updateGeneralBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}


func (b *updateGeneralBuilder) initQuery() {
	b.query += `
	UPDATE "data_settings" SET`
}

func (b *updateGeneralBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Tel != "" {
		b.values = append(b.values, b.req.Tel)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"tel" = $%d`, b.lastStackIndex))
	}

	if b.req.Email != "" {
		b.values = append(b.values, b.req.Email)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"email" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkFacebook != "" {
		b.values = append(b.values, b.req.LinkFacebook)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_facebook" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkInstagram != "" {
		b.values = append(b.values, b.req.LinkInstagram)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_instagram" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkTwitter != "" {
		b.values = append(b.values, b.req.LinkTwitter)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_twitter" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkTikTok != "" {
		b.values = append(b.values, b.req.LinkTikTok)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_tiktok" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkLine != "" {
		b.values = append(b.values, b.req.LinkLine)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_line" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkWebsite != "" {
		b.values = append(b.values, b.req.LinkWebsite)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_website" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateGeneralBuilder) insertImages() error {
	query := `
	INSERT INTO "data_setting_images" (
		"filename",
		"url",
		"data_setting_id"
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

func (b *updateGeneralBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "data_setting_images"
	WHERE "data_setting_id" = $1;`

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

func (b *updateGeneralBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "data_setting_images"
	WHERE "data_setting_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/data_setting/%s", img.FileName),
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

func (b *updateGeneralBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateGeneralBuilder) updateGeneral() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update general failed: %v", err)
	}
	return nil
}

func (b *updateGeneralBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateGeneralBuilder) getValues() []any         { return b.values }
func (b *updateGeneralBuilder) getQuery() string         { return b.query }
func (b *updateGeneralBuilder) setQuery(query string)    { b.query = query }
func (b *updateGeneralBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateGeneralBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateGeneralEngineer(b IUpdateGeneralBuilder) *updateGeneralEngineer {
	return &updateGeneralEngineer{builder: b}
}

func (en *updateGeneralEngineer) sumQueryFields() {
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


func (en *updateGeneralEngineer) UpdateGeneral() error {
	en.builder.initTransaction()
	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Update job
	if err := en.builder.updateGeneral(); err != nil {
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