package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/users"
)

type IInsertUser interface {
	initTransaction() error
	insertUser() error
	insertUserImage() error
	insertRole() error
	commit() error
	getUserId() string
	Result() (*users.UserPassport, error)
}

type insertUserBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *users.User
}

func InsertUserBuilder(db *sqlx.DB, req *users.User) IInsertUser {
	return &insertUserBuilder{
		db:  db,
		req: req,
	}
}

type insertUserEngineer struct {
	builder IInsertUser
}

func (b *insertUserBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertUserBuilder) insertUser() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"name",
		"tel",
		"display"
	)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING "id";`

	if err := b.db.QueryRowContext(
		ctx,
		query,
		b.req.Email,
		b.req.Password,
		b.req.Username,
		b.req.Name,
		b.req.Tel,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return fmt.Errorf("email has been used")
		default:
			return fmt.Errorf("insert user failed: %v", err)
		}
	}

	return nil
}

func (b *insertUserBuilder) insertUserImage() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "user_images" (
		"filename",
		"url",
		"user_id"
	)
	VALUES`

	valueStack := make([]interface{}, 0) // Change any to interface{}
	var index int
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *insertUserBuilder) insertRole() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "user_roles" (
		"user_id",
		"role_id"
	)
	VALUES ($1, $2);
	`

	for _, role := range b.req.UserRole {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			role.RoleId,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert user roles failed: %v", err)
		}
	}

	return nil
}

func (b *insertUserBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertUserBuilder) getUserId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertUserEngineer(b IInsertUser) *insertUserEngineer {
	return &insertUserEngineer{builder: b}
}

func (en *insertUserEngineer) InsertUser() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertUser(); err != nil {
		return "", err
	}
	if err := en.builder.insertUserImage(); err != nil {
		return "", err
	}
	if err := en.builder.insertRole(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getUserId(), nil
}

func (b *insertUserBuilder) Result() (*users.UserPassport, error) {
	query := `
	SELECT
		json_build_object(
			'user',"t",
			'token', NULL
		)
		FROM (
			SELECT 
				"u"."id",
				"u"."email",
				"u"."username"
			FROM "users" "u"
			WHERE "u"."id" = $1
		) AS "t"
		`

	data := make([]byte, 0)
	if err := b.db.Get(&data, query, &b.req.Id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed %v", err)
	}
	return user, nil
}
