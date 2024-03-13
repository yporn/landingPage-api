package activitiesPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/activities"
)

type IInsertActivityBuidler interface {
	initTransaction() error
	insertActivity() error
	insertAttachment() error
	commit() error
	getActivityId() string
}

type insertActivityBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *activities.Activity
}

func InsertActivityBuilder(db *sqlx.DB, req *activities.Activity) IInsertActivityBuidler {
	return &insertActivityBuilder{
		db:  db,
		req: req,
	}
}

type insertActivityEngineer struct {
	builder IInsertActivityBuidler
}

func (b *insertActivityBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertActivityBuilder) insertActivity() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "activities" (
		"index",
		"heading",
		"description",
		"start_date",
		"end_date",
		"video_link",
		"display"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Index,
		b.req.Heading,
		b.req.Description,
		b.req.StartDate,
		b.req.EndDate,
		b.req.VideoLink,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert activity failed: %v", err)
	}
	return nil
}

func (b *insertActivityBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

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
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *insertActivityBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertActivityBuilder) getActivityId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertActivityEngineer(b IInsertActivityBuidler) *insertActivityEngineer {
	return &insertActivityEngineer{builder: b}
}

func (en *insertActivityEngineer) InsertActivity() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertActivity(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getActivityId(), nil
}