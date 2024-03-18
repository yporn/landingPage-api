package projectsPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/projects"
)

type IInsertProjectBuilder interface {
	initTransaction() error
	insertProject() error
	insertHouseTypeItem() error
	insertDescAreaItem() error
	insertComfortableItem() error
	insertAttachment() error
	commit() error
	getProjectId() string
}

type insertProjectBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *projects.Project
}

func InsertProjectBuilder(db *sqlx.DB, req *projects.Project) IInsertProjectBuilder {
	return &insertProjectBuilder{
		db:  db,
		req: req,
	}
}

type insertProjectEngineer struct {
	builder IInsertProjectBuilder
}

func (b *insertProjectBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertProjectBuilder) insertProject() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "projects" (
		"name",
		"index",
		"heading",
		"text",
		"location",
		"price",
		"status_project",
		"type_project",
		"description",
		"name_facebook",
		"link_facebook",
		"tel",
		"address",
		"link_location",
		"display"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING "id";
	`
	if err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.Name,
		b.req.Index,
		b.req.Heading,
		b.req.Text,
		b.req.Location,
		b.req.Price,
		b.req.StatusProject,
		b.req.TypeProject,
		b.req.Description,
		b.req.NameFacebook,
		b.req.LinkFacebook,
		b.req.Tel,
		b.req.Address,
		b.req.LinkLocation,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert project failed: %v", err)
	}
	return nil
}

func (b *insertProjectBuilder) insertHouseTypeItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "project_house_type_items" (
		"project_id",
		"name"
	)
	VALUES ($1, $2);
	`

	for _, houseTypeItem := range b.req.HouseTypeItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			houseTypeItem.Name,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert project house type item failed: %v", err)
		}
	}

	return nil
}

func (b *insertProjectBuilder) insertDescAreaItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "project_desc_area_items" (
		"project_id",
		"item",
		"amount",
		"unit"
	)
	VALUES ($1, $2, $3, $4);
	`

	for _, descAreaItem := range b.req.DescAreaItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			descAreaItem.ItemArea,
			descAreaItem.Amount,
			descAreaItem.Unit,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert project area items failed: %v", err)
		}
	}

	return nil
}

func (b *insertProjectBuilder) insertComfortableItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "project_comfortable_items" (
		"project_id",
		"item"
	)
	VALUES ($1, $2);
	`
	for _, comfortableItem := range b.req.ComfortableItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			comfortableItem.Item,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert project Facilities items failed: %v", err)
		}
	}

	return nil
}

func (b *insertProjectBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "project_images" (
		"filename",
		"url",
		"project_id"
	)
	VALUES`

	valueStack := make([]any, 0)
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

func (b *insertProjectBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (b *insertProjectBuilder) getProjectId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertProjectEngineer(b IInsertProjectBuilder) *insertProjectEngineer {
	return &insertProjectEngineer{builder: b}
}

func (en *insertProjectEngineer) InsertProject() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertProject(); err != nil {
		return "", err
	}
	if err := en.builder.insertHouseTypeItem(); err != nil {
		return "", err
	}
	if err := en.builder.insertDescAreaItem(); err != nil {
		return "", err
	}
	if err := en.builder.insertComfortableItem(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getProjectId(), nil
}
