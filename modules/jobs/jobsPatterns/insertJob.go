package jobsPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/jobs"
)

type IInsertJobBuilder interface {
	initTransaction() error
	insertJob() error
	getJobId() string
	commit() error
}


type insertJobBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *jobs.Job
}

func InsertJobBuilder(db *sqlx.DB, req *jobs.Job) IInsertJobBuilder {
	return &insertJobBuilder{
		db:  db,
		req: req,
	}
}

type insertJobEngineer struct {
	builder IInsertJobBuilder
}

func InsertJobEngineer(b IInsertJobBuilder) *insertJobEngineer {
	return &insertJobEngineer{builder: b}
}

func (b *insertJobBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertJobBuilder) insertJob() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "careers" (
		"position",
		"amount",
		"location",
		"description",
		"qualification",
		"start_date",
		"end_date",
		"status",
		"display"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Position,
		b.req.Amount,
		b.req.Location,
		b.req.Description,
		b.req.Qualification,
		b.req.StartDate,
		b.req.EndDate,
		b.req.Status,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert jobs failed: %v", err)
	}
	return nil
}

func (b *insertJobBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertJobBuilder) getJobId() string {
	return strconv.Itoa(b.req.Id)
}

func (en *insertJobEngineer) InsertJob() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertJob(); err != nil {
		return "", err
	}
	
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getJobId(), nil
}