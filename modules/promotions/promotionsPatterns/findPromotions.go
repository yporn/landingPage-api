package promotionsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/promotions"
	"github.com/yporn/sirarom-backend/pkg/utils"
)



type IFindPromotionBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*promotions.Promotion
	Count() int
	PrintQuery()
}

type findPromotionBuilder struct {
	db             *sqlx.DB
	req            *promotions.PromotionFilter
	query          string
	values         []any
	lastStackIndex int
}

func FindPromotionBuilder(db *sqlx.DB, req *promotions.PromotionFilter) IFindPromotionBuilder {
	return &findPromotionBuilder{
		db:  db,
		req: req,
	}
}

func (b *findPromotionBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findPromotionBuilder) initQuery() {
	b.query += `
		SELECT
			"p".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "promotion_images" "i"
					WHERE "i"."promotion_id" = "p"."id"
				) AS "it"
			) AS "promotion_images"
		FROM "promotions" "p"
		WHERE 1 = 1`
}

func (b *findPromotionBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "promotions" "p"
		WHERE 1 = 1`
}

func (b *findPromotionBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "p"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("p"."heading") LIKE ?)`)
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

func (b *findPromotionBuilder) sort() {
	orderByMap := map[string]string{
		"id":         "\"id\"",
		"heading":   "\"heading\"",
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

func (b *findPromotionBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findPromotionBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findPromotionBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findPromotionBuilder) Result() []*promotions.Promotion {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	promotionsData := make([]*promotions.Promotion, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find promotions failed: %v\n", err)
		return make([]*promotions.Promotion, 0)
	}

	if err := json.Unmarshal(bytes, &promotionsData); err != nil {
		log.Printf("unmarshal promotions failed: %v\n", err)
		return make([]*promotions.Promotion, 0)
	}
	b.resetQuery()
	return promotionsData
}

func (b *findPromotionBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count promotions failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findPromotionBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findPromotionEngineer struct {
	builder IFindPromotionBuilder
}

func FindPromotionEngineer(builder IFindPromotionBuilder) *findPromotionEngineer {
	return &findPromotionEngineer{builder: builder}
}

func (en *findPromotionEngineer) FindPromotion() IFindPromotionBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findPromotionEngineer) CountPromotion() IFindPromotionBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
