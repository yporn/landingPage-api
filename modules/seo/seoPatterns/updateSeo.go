package seoPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/seo"
)

type IUpdateSeoBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	closeQuery()
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	commit() error
	updateSeo() error
}

type updateSeoBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *seo.Seo
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateSeoBuilder(db *sqlx.DB, req *seo.Seo) IUpdateSeoBuilder {
	return &updateSeoBuilder{
		db:            db,
		req:           req,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateSeoEngineer struct {
	builder IUpdateSeoBuilder
}

func (b *updateSeoBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateSeoBuilder) initQuery() {
	b.query += `
	UPDATE "seo" SET`
}

func (b *updateSeoBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"title" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.Keyword != "" {
		b.values = append(b.values, b.req.Keyword)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"keyword" = $%d`, b.lastStackIndex))
	}

	if b.req.Robot != "" {
		b.values = append(b.values, b.req.Robot)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"robot" = $%d`, b.lastStackIndex))
	}

	if b.req.GoogleBot != "" {
		b.values = append(b.values, b.req.GoogleBot)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"google_bot" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateSeoBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}


func (b *updateSeoBuilder) updateSeo() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update seo failed: %v", err)
	}
	return nil
}

func (b *updateSeoBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateSeoBuilder) getValues() []any         { return b.values }
func (b *updateSeoBuilder) getQuery() string         { return b.query }
func (b *updateSeoBuilder) setQuery(query string)    { b.query = query }
func (b *updateSeoBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateSeoEngineer(b IUpdateSeoBuilder) *updateSeoEngineer {
	return &updateSeoEngineer{builder: b}
}

func (en *updateSeoEngineer) sumQueryFields() {
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

func (en *updateSeoEngineer) UpdateSeo() error {
	en.builder.initTransaction()
	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	// Update job
	if err := en.builder.updateSeo(); err != nil {
		return err
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}