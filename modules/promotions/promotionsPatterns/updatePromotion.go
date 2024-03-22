package promotionsPatterns

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/promotions"
)

type IUpdatePromotionBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	updateFreeItem() error
	updateHouseModel() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updatePromotion() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updatePromotionBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *promotions.Promotion
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdatePromotionBuilder(db *sqlx.DB, req *promotions.Promotion, filesUsecases filesUsecases.IFilesUsecase) IUpdatePromotionBuilder {
	return &updatePromotionBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updatePromotionEngineer struct {
	builder IUpdatePromotionBuilder
}

func (b *updatePromotionBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updatePromotionBuilder) initQuery() {
	b.query += `
	UPDATE "promotions" SET`
}

func (b *updatePromotionBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Index != 0 {
		b.values = append(b.values, b.req.Index)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"index" = $%d`, b.lastStackIndex))
	}

	if b.req.Heading != "" {
		b.values = append(b.values, b.req.Heading)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"heading" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.StartDate != "" {
		b.values = append(b.values, b.req.StartDate)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"start_date" = $%d`, b.lastStackIndex))
	}

	if b.req.EndDate != "" {
		b.values = append(b.values, b.req.EndDate)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"end_date" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updatePromotionBuilder) updateFreeItem() error {
	// Retrieve existing items associated with the house model
	existingItems := make([]*promotions.PromotionFreeItem, 0)
	if err := b.db.Select(
		&existingItems,
		`SELECT "id", "promotion_id", "description" FROM "promotion_free_items" WHERE "promotion_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing promotion_free_items: %v", err)
	}

	// Compare existing items with new items
	for _, existingItem := range existingItems {
		itemFound := false
		for _, newItem := range b.req.FreeItem {
			if existingItem.Description == newItem.Description && existingItem.PromotionId == newItem.PromotionId {
				itemFound = true
				break
			}
		}
		// If existing items not found in the new items, delete it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "promotion_free_items" WHERE "id" = $1;`,
				existingItem.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing item: %v", err)
			}
		}
	}

	// Insert new type items
	for _, newItem := range b.req.FreeItem {
		itemFound := false
		for _, existingItem := range existingItems {
			if newItem.Description == existingItem.Description && newItem.PromotionId == existingItem.PromotionId {
				itemFound = true
				break
			}
		}
		// If new type item not found in existing type items, insert it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`INSERT INTO "promotion_free_items" ("description", "promotion_id") VALUES ($1, $2);`,
				newItem.Description, b.req.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new free item: %v", err)
			}
		}
	}
	return nil
}

func (b *updatePromotionBuilder) updateHouseModel() error {
	// Retrieve existing items associated with the house model
	existingItems := make([]*promotions.PromotionHouseModel, 0)
	if err := b.db.Select(
		&existingItems,
		`SELECT "id", "promotion_id", "house_model_id" FROM "promotion_house_models" WHERE "promotion_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing promotion_house_models: %v", err)
	}

	// Compare existing items with new items
	for _, existingItem := range existingItems {
		itemFound := false
		for _, newItem := range b.req.HouseModel {
			if existingItem.HouseModelId == newItem.HouseModelId {
				itemFound = true
				break
			}
		}
		// If existing items not found in the new items, delete it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "promotion_house_models" WHERE "id" = $1;`,
				existingItem.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing item: %v", err)
			}
		}
	}

	// Insert new type items
	for _, newItem := range b.req.HouseModel {
		itemFound := false
		for _, existingItem := range existingItems {
			if newItem.HouseModelId == existingItem.HouseModelId {
				itemFound = true
				break
			}
		}
		// If new type item not found in existing type items, insert it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`INSERT INTO "promotion_house_models" ("promotion_id", "house_model_id") VALUES ($1, $2);`,
				b.req.Id, newItem.HouseModelId,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new promotion house model item: %v", err)
			}
		}
	}
	return nil
}

func (b *updatePromotionBuilder) insertImages() error {
	query := `
	INSERT INTO "promotion_images" (
		"filename",
		"url",
		"promotion_id"
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
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}
func (b *updatePromotionBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "promotion_images"
	WHERE "promotion_id" = $1;`

	images := make([]*entities.Image, 0)
	if err := b.db.Select(
		&images,
		query,
		b.req.Id,
	); err != nil {
		return make([]*entities.Image, 0)
	}
	return images
}
func (b *updatePromotionBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "promotion_images"
	WHERE "promotion_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/promotions/%s", img.FileName),
			})
		}
		b.filesUsecases.DeleteFileOnStorage(deleteFileReq)
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete images failed: %v", err)
	}
	return nil
}
func (b *updatePromotionBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}
func (b *updatePromotionBuilder) updatePromotion() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update promotion failed: %v", err)
	}
	return nil
}
func (b *updatePromotionBuilder) getQueryFields() []string { return b.queryFields }
func (b *updatePromotionBuilder) getValues() []any         { return b.values }
func (b *updatePromotionBuilder) getQuery() string         { return b.query }
func (b *updatePromotionBuilder) setQuery(query string)    { b.query = query }
func (b *updatePromotionBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updatePromotionBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdatePromotionEngineer(b IUpdatePromotionBuilder) *updatePromotionEngineer {
	return &updatePromotionEngineer{builder: b}
}

func (en *updatePromotionEngineer) sumQueryFields() {
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
func (en *updatePromotionEngineer) UpdatePromotion() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())


	if err := en.builder.updatePromotion(); err != nil {
		return err
	}

	if err := en.builder.updateFreeItem(); err != nil {
		return err
	}

	if err := en.builder.updateHouseModel(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}
