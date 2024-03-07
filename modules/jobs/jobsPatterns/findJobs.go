package jobsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/jobs"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindJobBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*jobs.Job
	Count() int
	PrintQuery()
}

type findJobBuilder struct {
	db             *sqlx.DB
	req            *jobs.JobFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindJobBuilder(db *sqlx.DB, req *jobs.JobFilter) IFindJobBuilder {
	return &findJobBuilder{
		db:  db,
		req: req,
	}
}

func (b *findJobBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}
func (b *findJobBuilder) initQuery() {
	b.query += `
		SELECT 
			"id",
			"position",
			"amount",
			"location",
			"start_date",
			"end_date",
			"display",
			"created_at",
			"updated_at"
		FROM "careers"
		WHERE 1=1
	`
}
func (b *findJobBuilder) countQuery() {
	b.query += `
	SELECT
		COUNT(*) AS "count"
	FROM "careers"
	WHERE 1 = 1`
}
func (b *findJobBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "id" = ?`)
	}

	// Search Check
	if b.req.Search != "" {
		b.values = append(b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)
		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("position") LIKE ? OR LOWER("location") LIKE ?)`)
	}

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
		} else {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
		}
	}
	// Last stack record
	b.lastStackIndex = len(b.values)

	// Summary query
	b.query += queryWhere

}
func (b *findJobBuilder) sort() {
	orderByMap := map[string]string{
		"id":         "\"id\"",
		"position":   "\"position\"",
		"location":   "\"location\"",
		"created_at": "\"created_at\"",
	}

	orderBy := orderByMap[b.req.OrderBy]
	if orderBy == "" {
		orderBy = orderByMap["created_at"]
	} else {
		orderBy = orderByMap[b.req.OrderBy]
	}

	sortOrder := strings.ToUpper(b.req.Sort)
	if sortOrder == "" {
		b.req.Sort = "ASC"
	}

	// b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
		ORDER BY %s %s`, orderBy, b.req.Sort)
	b.lastStackIndex = len(b.values)
}
func (b *findJobBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}
func (b *findJobBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}
func (b *findJobBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}
func (b *findJobBuilder) Result() []*jobs.Job {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	jobsData := make([]*jobs.Job, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find jobs failed: %v\n", err)
		return make([]*jobs.Job, 0)
	}

	if err := json.Unmarshal(bytes, &jobsData); err != nil {
		log.Printf("unmarshal jobs failed: %v\n", err)
		return make([]*jobs.Job, 0)
	}
	b.resetQuery()
	return jobsData
}
func (b *findJobBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count jobs failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}
func (b *findJobBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findJobEngineer struct {
	builder IFindJobBuilder
}

func FindJobEngineer(builder IFindJobBuilder) *findJobEngineer {
	return &findJobEngineer{builder: builder}
}

func (en *findJobEngineer) FindJob() IFindJobBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findJobEngineer) CountJob() IFindJobBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
