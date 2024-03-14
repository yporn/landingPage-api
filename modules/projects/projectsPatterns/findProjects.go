package projectsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/projects"
)

type IFindProjectBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	// buildWhereDate()
	buildSort()
	buildPaginate()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDb() *sqlx.DB
	reset()
}

type findProjectBuilder struct {
	db        *sqlx.DB
	req       *projects.ProjectFilter
	query     string
	values    []any
	lastIndex int
}

func FindProjectBuilder(db *sqlx.DB, req *projects.ProjectFilter) IFindProjectBuilder {
	return &findProjectBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

type findProjectEngineer struct {
	builder IFindProjectBuilder
}

func FindProjectEngineer(b IFindProjectBuilder) *findProjectEngineer {
	return &findProjectEngineer{builder: b}
}

func (b *findProjectBuilder) initQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("at"))
	FROM (
		SELECT
			"p".*,
			(
				SELECT
				COALESCE(array_to_json(array_agg("htit")), '[]'::json)
				FROM (
					SELECT
						"hti".*
					FROM "project_house_type_items" "hti"
					WHERE "hti"."project_id" = "p"."id"
					
				) AS "htit"
			) AS "house_type_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("dait")), '[]'::json)
				FROM (
					SELECT
						"dai".*
					FROM "project_desc_area_items" "dai"
					WHERE "dai"."project_id" = "p"."id"
					
				) AS "dait"
			) AS "area_items",
			(
				SELECT
				COALESCE(array_to_json(array_agg("cit")), '[]'::json)
				FROM (
					SELECT
						"ci".*
					FROM "project_comfortable_items" "ci"
					WHERE "ci"."project_id" = "p"."id"
					
				) AS "cit"
			) AS "facilities_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "project_images" "i"
					WHERE "i"."project_id" = "p"."id"
				) AS "it"
			) AS "images"
			FROM "projects" "p"
			WHERE 1 = 1
	`
}

func (b *findProjectBuilder) initCountQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "projects" "p"
		WHERE 1 = 1`
}

func (b *findProjectBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := fmt.Sprintf(`
		AND (
			LOWER("o"."name") LIKE $%d OR
			LOWER("o"."type_project") LIKE $%d OR
			LOWER("o"."location") LIKE $%d
		)`,
			b.lastIndex+1,
			b.lastIndex+2,
			b.lastIndex+3,
		)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findProjectBuilder) buildWhereStatus() {
	if b.req.StatusProject != "" {
		b.values = append(
			b.values,
			strings.ToLower(b.req.StatusProject),
		)

		query := fmt.Sprintf(`
		AND "p"."status_project" = $%d`,
			b.lastIndex+1,
		)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findProjectBuilder) buildSort() {
	b.values = append(b.values, b.req.OrderBy)

	// Fix
	b.query += fmt.Sprintf(`
		ORDER BY $%d %s`, b.lastIndex+1, b.req.Sort)

	b.lastIndex = len(b.values)
}

func (b *findProjectBuilder) buildPaginate() {
	b.values = append(
		b.values,
		(b.req.Page-1)*b.req.Limit,
		b.req.Limit,
	)

	// Fix
	b.query += fmt.Sprintf(`
		OFFSET $%d LIMIT $%d`, b.lastIndex+1, b.lastIndex+2)

	b.lastIndex = len(b.values)
}

func (b *findProjectBuilder) closeQuery() {
	b.query += `
	) AS "at"`
}

func (b *findProjectBuilder) getQuery() string { return b.query }

func (b *findProjectBuilder) setQuery(query string) { b.query = query }

func (b *findProjectBuilder) getValues() []any { return b.values }

func (b *findProjectBuilder) setValues(data []any) { b.values = data }

func (b *findProjectBuilder) setLastIndex(n int) { b.lastIndex = n }

func (b *findProjectBuilder) getDb() *sqlx.DB { return b.db }

func (b *findProjectBuilder) reset() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

func (en *findProjectEngineer) FindProject() []*projects.Project {
	_, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	en.builder.initQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	// en.builder.buildWhereDate()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	raw := make([]byte, 0)
	if err := en.builder.getDb().Get(&raw, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("get projects failed: %v\n", err)
		return make([]*projects.Project, 0)
	}

	projectsData := make([]*projects.Project, 0)
	if err := json.Unmarshal(raw, &projectsData); err != nil {
		log.Printf("unmarshal projects failed: %v\n", err)
	}

	en.builder.reset()
	return projectsData
}

func (en *findProjectEngineer) CountOrder() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	en.builder.initCountQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	// en.builder.buildWhereDate()

	var count int
	if err := en.builder.getDb().Get(&count, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("count projects failed: %v\n", err)
		return 0
	}

	en.builder.reset()
	return count
}