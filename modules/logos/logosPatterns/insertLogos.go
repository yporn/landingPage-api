package logosPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/logos"
)

type IInsertLogoBuidler interface {
	initTransaction() error
	insertLogo() error
	insertAttachment() error
	commit() error
	getLogoId() string
}

type insertLogoBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *logos.Logo
}

func InsertLogoBuilder(db *sqlx.DB, req *logos.Logo) IInsertLogoBuidler {
	return &insertLogoBuilder{
		db:  db,
		req: req,
	}
}

type insertLogoEngineer struct {
	builder IInsertLogoBuidler
}

func (b *insertLogoBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertLogoBuilder) insertLogo() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "logos" (
		"index",
		"name",
		"display"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Index,
		b.req.Name,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert logo failed: %v", err)
	}
	return nil
}

func (b *insertLogoBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "logo_images" (
		"filename",
		"url",
		"logo_id"
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
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *insertLogoBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertLogoBuilder) getLogoId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertLogoEngineer(b IInsertLogoBuidler) *insertLogoEngineer {
	return &insertLogoEngineer{builder: b}
}

func (en *insertLogoEngineer) InsertLogo() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertLogo(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getLogoId(), nil
}