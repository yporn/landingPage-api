package bannersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/banners"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindBannerBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*banners.Banner
	Count() int
	PrintQuery()
}

type findBannerBuilder struct {
	db             *sqlx.DB
	req            *banners.BannerFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindBannerBuilder(db *sqlx.DB, req *banners.BannerFilter) IFindBannerBuilder {
	return &findBannerBuilder{
		db:  db,
		req: req,
	}
}

func (b *findBannerBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findBannerBuilder) initQuery() {
	b.query += `
		SELECT
			"b"."id",
			"b"."index",
			"b"."delay",
			"b"."display",
			"b"."created_at",
			"b"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "banner_images" "i"
					WHERE "i"."banner_id" = "b"."id"
				) AS "it"
			) AS "images"
		FROM "banners" "b"
		WHERE 1 = 1`
}

func (b *findBannerBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "banners" "b"
		WHERE 1 = 1`
}

func (b *findBannerBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "b"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("b"."index") LIKE ?)`)
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

func (b *findBannerBuilder) sort() {
	orderByMap := map[string]string{
		"id":    "\"b\".\"id\"",
		"index": "\"b\".\"index\"",

	}
	if orderByMap[b.req.OrderBy] == "" {
		b.req.OrderBy = orderByMap["index"]
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

func (b *findBannerBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findBannerBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findBannerBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findBannerBuilder) Result() []*banners.Banner {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	bannersData := make([]*banners.Banner, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find banners failed: %v\n", err)
		return make([]*banners.Banner, 0)
	}

	if err := json.Unmarshal(bytes, &bannersData); err != nil {
		log.Printf("unmarshal banners failed: %v\n", err)
		return make([]*banners.Banner, 0)
	}
	b.resetQuery()
	return bannersData
}

func (b *findBannerBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count banners failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findBannerBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findBannerEngineer struct {
	builder IFindBannerBuilder
}

func FindBannerEngineer(builder IFindBannerBuilder) *findBannerEngineer {
	return &findBannerEngineer{builder: builder}
}

func (en *findBannerEngineer) FindBanner() IFindBannerBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findBannerEngineer) CountBanner() IFindBannerBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}