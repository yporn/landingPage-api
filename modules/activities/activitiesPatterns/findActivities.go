package activitiesPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yporn/sirarom-backend/modules/activities"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindActivityBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*activities.Activity
	Count() int
	PrintQuery()
}

type findActivityBuilder struct {
	db             *sqlx.DB
	req            *activities.ActivityFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindActivityBuilder(db *sqlx.DB, req *activities.ActivityFilter) IFindActivityBuilder {
	return &findActivityBuilder{
		db:  db,
		req: req,
	}
}

func (b *findActivityBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findActivityBuilder) initQuery() {
	b.query += `
		SELECT
			"a"."id",
			"a"."index",
			"a"."heading",
			"a"."description",
			"a"."start_date",
			"a"."end_date",
			"a"."display",
			"a"."video_link",
			"a"."created_at",
			"a"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "activities_images" "i"
					WHERE "i"."activity_id" = "a"."id"
				) AS "it"
			) AS "images"
		FROM "activities" "a"
		WHERE 1 = 1`
}

func (b *findActivityBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "activities" "a"
		WHERE 1 = 1`
}

func (b *findActivityBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "a"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("a"."heading") LIKE ?)`)
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

func (b *findActivityBuilder) sort() {
	orderByMap := map[string]string{
		"id":         "\"a\".\"id\"",
		"heading":    "\"a\".\"heading\"",
		"created_at": "\"a\".\"created_at\"",
	}
	if orderByMap[b.req.OrderBy] == "" {
		b.req.OrderBy = orderByMap["created_at"]
	} else {
		b.req.OrderBy = orderByMap[b.req.OrderBy]
	}

	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[b.req.Sort] == "" {
		b.req.Sort = sortMap["ASC"]
	} else {
		b.req.Sort = sortMap[strings.ToUpper(b.req.Sort)]
	}

	b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
		ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
	b.lastStackIndex = len(b.values)
}

func (b *findActivityBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findActivityBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findActivityBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findActivityBuilder) Result() []*activities.Activity {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	activitiesData := make([]*activities.Activity, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find activities failed: %v\n", err)
		return make([]*activities.Activity, 0)
	}

	if err := json.Unmarshal(bytes, &activitiesData); err != nil {
		log.Printf("unmarshal banners failed: %v\n", err)
		return make([]*activities.Activity, 0)
	}
	b.resetQuery()
	return activitiesData
}

func (b *findActivityBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count activities failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findActivityBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findActivityEngineer struct {
	builder IFindActivityBuilder
}

func FindActivityEngineer(builder IFindActivityBuilder) *findActivityEngineer {
	return &findActivityEngineer{builder: builder}
}

func (en *findActivityEngineer) FindActivity() IFindActivityBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findActivityEngineer) CountActivity() IFindActivityBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
