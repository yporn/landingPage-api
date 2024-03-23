package logosPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/logos"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindLogoBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*logos.Logo
	Count() int
	PrintQuery()
}

type findLogoBuilder struct {
	db             *sqlx.DB
	req            *logos.LogoFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindLogoBuilder(db *sqlx.DB, req *logos.LogoFilter) IFindLogoBuilder {
	return &findLogoBuilder{
		db:  db,
		req: req,
	}
}

func (b *findLogoBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findLogoBuilder) initQuery() {
	b.query += `
		SELECT
			"l".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "logo_images" "i"
					WHERE "i"."logo_id" = "l"."id"
				) AS "it"
			) AS "images"
		FROM "logos" "l"
		WHERE 1 = 1`
}

func (b *findLogoBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "logos" "l"
		WHERE 1 = 1`
}

func (b *findLogoBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "l"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("a"."name") LIKE ?)`)
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

func (b *findLogoBuilder) sort() {
	orderByMap := map[string]string{
		"id":         "\"id\"",
		"name":   "\"name\"",
		"created_at": "\"created_at\"",
	}

	orderBy := orderByMap[b.req.OrderBy]
	if orderBy == "" {
		orderBy = orderByMap["id"]
	} else {
		orderBy = orderByMap[b.req.OrderBy]
	}

	sortOrder := strings.ToUpper(b.req.Sort)
	if sortOrder == "" {
		b.req.Sort = "desc"
	}

	// b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`
		ORDER BY %s %s`, orderBy, b.req.Sort)
	b.lastStackIndex = len(b.values)
}

func (b *findLogoBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findLogoBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findLogoBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findLogoBuilder) Result() []*logos.Logo {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	logosData := make([]*logos.Logo, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find logos failed: %v\n", err)
		return make([]*logos.Logo, 0)
	}

	if err := json.Unmarshal(bytes, &logosData); err != nil {
		log.Printf("unmarshal logos failed: %v\n", err)
		return make([]*logos.Logo, 0)
	}
	b.resetQuery()
	return logosData
}

func (b *findLogoBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count logos failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findLogoBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findLogoEngineer struct {
	builder IFindLogoBuilder
}

func FindLogoEngineer(builder IFindLogoBuilder) *findLogoEngineer {
	return &findLogoEngineer{builder: builder}
}

func (en *findLogoEngineer) FindLogo() IFindLogoBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findLogoEngineer) CountLogo() IFindLogoBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}

