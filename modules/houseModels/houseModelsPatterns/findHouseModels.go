package houseModelsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/houseModels"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindHouseModelBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*houseModels.HouseModel
	Count() int
	PrintQuery()
}

type findHouseModelBuilder struct {
	db             *sqlx.DB
	projectId      string
	req            *houseModels.HouseModelFilter
	query          string
	values         []any
	lastStackIndex int
}

func FindHouseModelBuilder(db *sqlx.DB, projectId string, req *houseModels.HouseModelFilter) IFindHouseModelBuilder {
	return &findHouseModelBuilder{
		db:     db,
		projectId:  projectId,
		req:    req,
		values: make([]any, 0),
	}
}

type findHouseModelEngineer struct {
	builder IFindHouseModelBuilder
}

func FindHouseModelEngineer(b IFindHouseModelBuilder) *findHouseModelEngineer {
	return &findHouseModelEngineer{builder: b}
}

func (b *findHouseModelBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findHouseModelBuilder) initQuery() {
	b.query += `
	SELECT to_jsonb("at")
	FROM (
		SELECT
			"hm".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("hmi")), '[]'::json)
				FROM (
					SELECT
						"hmi".*
					FROM "house_model_type_items" "hmi"
					WHERE "hmi"."house_model_id" = "hm"."id"
				) AS "hmi"
			) AS "type_items",
			(
				SELECT
					COALESCE(array_to_json(array_agg("ihm")), '[]'::json)
				FROM (
					SELECT
						"ihm"."id",
						"ihm"."filename",
						"ihm"."url"
					FROM "house_model_images" "ihm"
					WHERE "ihm"."house_model_id" = "hm"."id"
				) AS "ihm"
			) AS "house_images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("hmp")), '[]'::json)
				FROM (
					SELECT
						"hmp".*,
						(
							SELECT
								COALESCE(array_to_json(array_agg("hmpi")), '[]'::json)
							FROM (
								SELECT
									"hmpi".*
								FROM "house_model_plan_items" "hmpi"
								WHERE "hmpi"."house_model_plan_id" = "hmp"."id"
							) AS "hmpi"
						) AS "plan_items",
						(
							SELECT
								COALESCE(array_to_json(array_agg("ihmp")), '[]'::json)
							FROM (
								SELECT
									"ihmp"."id",
									"ihmp"."filename",
									"ihmp"."url"
								FROM "house_model_plan_images" "ihmp"
								WHERE "ihmp"."house_model_plan_id" = "hmp"."id"
							) AS "ihmp"
						) AS "plan_images"
					FROM "house_model_plans" "hmp"
					WHERE "hmp"."house_model_id" = "hm"."id"
				) AS "hmp"
			) AS "house_plan",
			(
				SELECT
					COALESCE(array_to_json(array_agg("ptm")), '[]'::json)
				FROM (
					SELECT
						"ptm".*,
						(
							SELECT
								COALESCE(array_to_json(array_agg("pt")), '[]'::json)
							FROM (
								SELECT
									"pt".*
								FROM "promotions" "pt"
								WHERE "pt"."id" = "ptm"."promotion_id"
							) AS "pt"
						) AS "promotions",
					FROM "promotion_house_models" "ptm"
					WHERE "ptm"."house_model_id" = "hm"."id"
				) AS "ptm"
			) AS "promotions"
			FROM "house_models" "hm"
		WHERE "hm"."project_id" = $1
	)
	`
}

func (b *findHouseModelBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "house_models" "hm"
		WHERE "hm"."project_id" = $1
	`
}

func (b *findHouseModelBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "hm"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("hm"."name") LIKE ?)`)
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

func (b *findHouseModelBuilder) sort() {
	orderByMap := map[string]string{
		"id":   "\"p\".\"id\"",
		"name": "\"p\".\"name\"",
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

	b.query += fmt.Sprintf(`
		ORDER BY %s %s`, orderBy, b.req.Sort)
	b.lastStackIndex = len(b.values)
}

func (b *findHouseModelBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findHouseModelBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findHouseModelBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findHouseModelBuilder) Result() []*houseModels.HouseModel {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	houseModelsData := make([]*houseModels.HouseModel, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find house models failed: %v\n", err)
		return make([]*houseModels.HouseModel, 0)
	}

	if err := json.Unmarshal(bytes, &houseModelsData); err != nil {
		log.Printf("unmarshal house models failed: %v\n", err)
		return make([]*houseModels.HouseModel, 0)
	}
	b.resetQuery()
	return houseModelsData
}

func (b *findHouseModelBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count house models failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findHouseModelBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

func (en *findHouseModelEngineer) FindHouseModel() IFindHouseModelBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findHouseModelEngineer) CountHouseModel() IFindHouseModelBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}
