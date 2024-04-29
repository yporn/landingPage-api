package usersPatterns

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/users"

)

type IUpdateUserBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery() 
	updateRole() error
	updateImages() error
	closeQuery()
	updateUser() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateUserBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *users.User
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateUserBuilder(db *sqlx.DB, req *users.User, filesUsecases filesUsecases.IFilesUsecase) IUpdateUserBuilder {
	return &updateUserBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateUserEngineer struct {
	builder IUpdateUserBuilder
}

func (b *updateUserBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateUserBuilder) initQuery() {
	b.query += `
	UPDATE "users" SET `
}

func (b *updateUserBuilder) updateQuery() {
	setStatements := []string{}

	if b.req.Email != "" {
		b.values = append(b.values, b.req.Email)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"email" = $%d`, b.lastStackIndex))
	}

	if b.req.Username != "" {
		b.values = append(b.values, b.req.Username)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"username" = $%d`, b.lastStackIndex))
	}

	if b.req.Password != "" {
		fmt.Println("password : ",b.req.Password)
		b.values = append(b.values, b.req.Password)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"password" = $%d`, b.lastStackIndex))
	}
	if b.req.Name != "" {
		b.values = append(b.values, b.req.Name)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"name" = $%d`, b.lastStackIndex))
	}

	if b.req.Tel != "" {
		b.values = append(b.values, b.req.Tel)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"tel" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateUserBuilder) updateImages() error {
	// Retrieve existing images associated with the house model
	existingImages := make([]*entities.Image, 0)
	if err := b.db.Select(
		&existingImages,
		`SELECT "id", "filename", "url" FROM "user_images" WHERE "user_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing images: %v", err)
	}

	// Compare existing images with new images
	for _, existingImage := range existingImages {
		imageFound := false
		for _, newImage := range b.req.Images {
			if existingImage.FileName == newImage.FileName {
				imageFound = true
				break
			}
		}
		// If existing image not found in the new images, delete it
		if !imageFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "user_images" WHERE "id" = $1;`,
				existingImage.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing image: %v", err)
			}
			// Also delete the file from storage
			b.filesUsecases.DeleteFileOnStorage([]*files.DeleteFileReq{
				{Destination: fmt.Sprintf("images/users/%s", path.Base(existingImage.Url))},
			})
		}
	}

	// Insert new images
	for _, newImage := range b.req.Images {
		imageFound := false
		for _, existingImage := range existingImages {
			if newImage.FileName == existingImage.FileName {
				imageFound = true
				break
			}
		}
		// If new image not found in existing images, insert it
		if !imageFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`INSERT INTO "user_images" ("filename", "url", "user_id") VALUES ($1, $2, $3);`,
				newImage.FileName, newImage.Url, b.req.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new image: %v", err)
			}
		}
	}

	return nil
}

func (b *updateUserBuilder) updateRole() error {
	// Retrieve existing items associated with the house model
	existingItems := make([]*users.UserRole, 0)
	if err := b.db.Select(
		&existingItems,
		`SELECT "id", "user_id", "role_id" FROM "user_roles" WHERE "user_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing user_roles: %v", err)
	}

	// Compare existing items with new items
	for _, existingItem := range existingItems {
		itemFound := false
		for _, newItem := range b.req.UserRole {
			if existingItem.RoleId == newItem.RoleId {
				itemFound = true
				break
			}
		}
		// If existing items not found in the new items, delete it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "user_roles" WHERE "id" = $1;`,
				existingItem.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing item: %v", err)
			}
		}
	}

	// Insert new type items
	for _, newItem := range b.req.UserRole {
		itemFound := false
		for _, existingItem := range existingItems {
			if newItem.RoleId == existingItem.RoleId {
				itemFound = true
				break
			}
		}
		// If new type item not found in existing type items, insert it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`INSERT INTO "user_roles" ("role_id", "user_id") VALUES ($1, $2);`,
				newItem.RoleId, b.req.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new type item: %v", err)
			}
		}
	}
	return nil
}


func (b *updateUserBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateUserBuilder) updateUser() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update user failed: %v", err)
	}
	return nil
}

func (b *updateUserBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateUserBuilder) getValues() []any         { return b.values }
func (b *updateUserBuilder) getQuery() string         { return b.query }
func (b *updateUserBuilder) setQuery(query string)    { b.query = query }
func (b *updateUserBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateUserBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateUserEngineer(b IUpdateUserBuilder) *updateUserEngineer {
	return &updateUserEngineer{builder: b}
}

func (en *updateUserEngineer) sumQueryFields() {
	en.builder.updateQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
}

func (en *updateUserEngineer) UpdateUser() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	if err := en.builder.updateUser(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.updateImages(); err != nil {
			return err
		}
	}

	if err := en.builder.updateRole(); err != nil {
		return err
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}