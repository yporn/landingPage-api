package jobsPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/jobs"
)

type IUpdateJobBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	closeQuery()
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	commit() error
	updateJob() error 
}

type updateJobBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *jobs.Job
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateJobBuilder(db *sqlx.DB, req *jobs.Job) IUpdateJobBuilder {
	return &updateJobBuilder{
		db:            db,
		req:           req,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateJobEngineer struct {
	builder IUpdateJobBuilder
}

func (b *updateJobBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateJobBuilder) initQuery() {
	b.query += `
	UPDATE "careers" SET`
}

func (b *updateJobBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Position != "" {
		b.values = append(b.values, b.req.Position)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"position" = $%d`, b.lastStackIndex))
	}

	if b.req.Amount != "" {
		b.values = append(b.values, b.req.Amount)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"amount" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.Location != "" {
		b.values = append(b.values, b.req.Location)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"location" = $%d`, b.lastStackIndex))
	}

	if b.req.Qualification != "" {
		b.values = append(b.values, b.req.Qualification)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"qualification" = $%d`, b.lastStackIndex))
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

	if b.req.Status != "" {
		b.values = append(b.values, b.req.Status)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"status" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}



func (b *updateJobBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateJobBuilder) updateJob() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update job failed: %v", err)
	}
	return nil
}

func (b *updateJobBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateJobBuilder) getValues() []any         { return b.values }
func (b *updateJobBuilder) getQuery() string         { return b.query }
func (b *updateJobBuilder) setQuery(query string)    { b.query = query }

func (b *updateJobBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateJobEngineer(b IUpdateJobBuilder) *updateJobEngineer {
	return &updateJobEngineer{builder: b}
}

func (en *updateJobEngineer) sumQueryFields() {
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

func (en *updateJobEngineer) UpdateJob() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Update job
	if err := en.builder.updateJob(); err != nil {
		return err
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}

