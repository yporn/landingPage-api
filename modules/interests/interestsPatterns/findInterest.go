package interestsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindInterestBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*interests.Interest
	Count() int
	PrintQuery()
}

type findInterestBuilder struct {
	db             *sqlx.DB
	req            *interests.InterestFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindInterestBuilder(db *sqlx.DB, req *interests.InterestFilter) IFindInterestBuilder {
	return &findInterestBuilder{
		db:  db,
		req: req,
	}
}

func (b *findInterestBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findInterestBuilder) initQuery() {
	b.query += `
		SELECT
			"bi".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "interest_images" "i"
					WHERE "i"."interest_id" = "bi"."id"
				) AS "it"
			) AS "images"
		FROM "interests" "bi"
		WHERE 1 = 1`
}

func (b *findInterestBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "interests" "bi"
		WHERE 1 = 1`
}

func (b *findInterestBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "bi"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("bi"."bank_name") LIKE ?)`)
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

func (b *findInterestBuilder) sort() {
	orderByMap := map[string]string{
		"id":         "\"id\"",
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

func (b *findInterestBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findInterestBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findInterestBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findInterestBuilder) Result() []*interests.Interest {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	interestsData := make([]*interests.Interest, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find interests failed: %v\n", err)
		return make([]*interests.Interest, 0)
	}

	if err := json.Unmarshal(bytes, &interestsData); err != nil {
		log.Printf("unmarshal interests failed: %v\n", err)
		return make([]*interests.Interest, 0)
	}
	b.resetQuery()
	return interestsData
}

func (b *findInterestBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count interests failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findInterestBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findInterestEngineer struct {
	builder IFindInterestBuilder
}

func FindInterestEngineer(builder IFindInterestBuilder) *findInterestEngineer {
	return &findInterestEngineer{builder: builder}
}

func (en *findInterestEngineer) FindInterest() IFindInterestBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}


func (en *findInterestEngineer) CountInterest() IFindInterestBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
