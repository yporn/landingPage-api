package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/users"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type IFindUserBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*users.User
	Count() int
	PrintQuery()
}

type findUserBuilder struct {
	db             *sqlx.DB
	req            *users.UserFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindUserBuilder(db *sqlx.DB, req *users.UserFilter) IFindUserBuilder {
	return &findUserBuilder{
		db:  db,
		req: req,
	}
}

func (b *findUserBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findUserBuilder) initQuery() {
	b.query += `
		SELECT
			"u".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "user_images" "i"
					WHERE "i"."user_id" = "u"."id"
				) AS "it"
			) AS "images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("r")), '[]'::json)
				FROM (
					SELECT
						"r"."id",
						"r"."user_id",
						"r"."role_id"
					FROM "user_roles" "r"
					WHERE "r"."user_id" = "u"."id"
				) AS "r"
			) AS "roles"
		FROM "users" "u"
		WHERE 1 = 1`
}

func (b *findUserBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "users" "u"
		WHERE 1 = 1`
}
func (b *findUserBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "u"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("u"."name") LIKE ? OR LOWER("u"."username") LIKE ? OR LOWER("u"."email") LIKE ?)`)
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

func (b *findUserBuilder) sort() {
	orderByMap := map[string]string{
		"id":       "\"id\"",
		"name":     "\"name\"",
		"username": "\"username\"",
		"email":    "\"email\"",
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

func (b *findUserBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}

func (b *findUserBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findUserBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findUserBuilder) Result() []*users.User {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	usersData := make([]*users.User, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find users failed: %v\n", err)
		return make([]*users.User, 0)
	}

	if err := json.Unmarshal(bytes, &usersData); err != nil {
		log.Printf("unmarshal users failed: %v\n", err)
		return make([]*users.User, 0)
	}
	b.resetQuery()
	return usersData
}

func (b *findUserBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count users failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}

func (b *findUserBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findUserEngineer struct {
	builder IFindUserBuilder
}

func FindUserEngineer(builder IFindUserBuilder) *findUserEngineer {
	return &findUserEngineer{builder: builder}
}

func (en *findUserEngineer) FindUser() IFindUserBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findUserEngineer) CountUser() IFindUserBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}