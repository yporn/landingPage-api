package bannersPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/banners"
)

type IInsertBannerBuidler interface {
	initTransaction() error
	insertBanner() error
	insertAttachment() error
	commit() error
	getBannerId() string
}

type insertBannerBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *banners.Banner
}

func InsertBannerBuilder(db *sqlx.DB, req *banners.Banner) IInsertBannerBuidler {
	return &insertBannerBuilder{
		db:  db,
		req: req,
	}
}

type insertBannerEngineer struct {
	builder IInsertBannerBuidler
}

func (b *insertBannerBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertBannerBuilder) insertBanner() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "banners" (
		"index",
		"delay",
		"display"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Index,
		b.req.Delay,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert banner failed: %v", err)
	}
	return nil
}

func (b *insertBannerBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

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
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *insertBannerBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertBannerBuilder) getBannerId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertBannerEngineer(b IInsertBannerBuidler) *insertBannerEngineer {
	return &insertBannerEngineer{builder: b}
}

func (en *insertBannerEngineer) InsertBanner() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertBanner(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getBannerId(), nil
}