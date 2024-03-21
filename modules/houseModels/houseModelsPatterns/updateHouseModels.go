package houseModelsPatterns

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels"
)

type IUpdateHouseModelBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	updateImagesHouseModel() error
	updateTypeItem() error
	updateHousePlan() error
	updateHousePlanItems() error
	updateImagesHousePlan() error
	closeQuery()
	updateHouseModel() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateHouseModelBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *houseModels.HouseModel
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateHouseModelBuilder(db *sqlx.DB, req *houseModels.HouseModel, filesUsecases filesUsecases.IFilesUsecase) IUpdateHouseModelBuilder {
	return &updateHouseModelBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

func (b *updateHouseModelBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateHouseModelBuilder) initQuery() {
	b.query += `
	UPDATE "house_models" SET `
}

func (b *updateHouseModelBuilder) updateQuery() {
	setStatements := []string{}

	if b.req.Name != "" {
		b.values = append(b.values, b.req.Name)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"name" = $%d`, b.lastStackIndex))
	}
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}
	if b.req.LinkVideo != "" {
		b.values = append(b.values, b.req.LinkVideo)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"link_video" = $%d`, b.lastStackIndex))
	}
	if b.req.LinkVirtualTour != "" {
		b.values = append(b.values, b.req.LinkVirtualTour)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"link_virtual_tour" = $%d`, b.lastStackIndex))
	}
	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}
	if b.req.Index != 0 {
		b.values = append(b.values, b.req.Index)
		b.lastStackIndex = len(b.values)
		setStatements = append(setStatements, fmt.Sprintf(`"index" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

// func (b *updateHouseModelBuilder) updateImagesHouseModel() error {
// 	getData := `
// 	SELECT
// 		"id",
// 		"filename",
// 		"url"
// 	FROM "house_model_images"
// 	WHERE "house_model_id" = $1;`

// 	images := make([]*entities.Image, 0)
// 	if err := b.db.Select(
// 		&images,
// 		getData,
// 		b.req.Id,
// 	); err != nil {
// 		return fmt.Errorf("failed to retrieve images: %v", err)
// 	}

// 	deleteData := `
// 	DELETE FROM "house_model_images"
// 	WHERE "house_model_id" = $1;`

// 	if len(images) > 0 {
// 		deleteFileReq := make([]*files.DeleteFileReq, 0)
// 		for _, img := range images {
// 			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
// 				Destination: fmt.Sprintf("images/house_models/%s", img.FileName),
// 			})
// 		}
// 		b.filesUsecases.DeleteFileOnStorage(deleteFileReq)
// 	}

// 	if _, err := b.tx.ExecContext(
// 		context.Background(),
// 		deleteData,
// 		b.req.Id,
// 	); err != nil {
// 		b.tx.Rollback()
// 		return fmt.Errorf("failed to delete images: %v", err)
// 	}

// 	query := `
// 	INSERT INTO "house_model_images" (
// 		"filename",
// 		"url",
// 		"house_model_id"
// 	)
// 	VALUES`

// 	valueStack := make([]any, 0)
// 	var index int
// 	for i := range b.req.Images {
// 		valueStack = append(valueStack,
// 			b.req.Images[i].FileName,
// 			b.req.Images[i].Url,
// 			b.req.Id,
// 		)

// 		if i != len(b.req.Images)-1 {
// 			query += fmt.Sprintf(`
// 			($%d, $%d, $%d),`, index+1, index+2, index+3)
// 		} else {
// 			query += fmt.Sprintf(`
// 			($%d, $%d, $%d);`, index+1, index+2, index+3)
// 		}
// 		index += 3
// 	}

// 	if _, err := b.tx.ExecContext(
// 		context.Background(),
// 		query,
// 		valueStack...,
// 	); err != nil {
// 		b.tx.Rollback()
// 		return fmt.Errorf("failed to update images: %v", err)
// 	}
// 	return nil
// }

func (b *updateHouseModelBuilder) updateImagesHouseModel() error {
	// Retrieve existing images associated with the house model
	existingImages := make([]*entities.Image, 0)
	if err := b.db.Select(
		&existingImages,
		`SELECT "id", "filename", "url" FROM "house_model_images" WHERE "house_model_id" = $1;`,
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
				`DELETE FROM "house_model_images" WHERE "id" = $1;`,
				existingImage.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing image: %v", err)
			}
			// Also delete the file from storage
			b.filesUsecases.DeleteFileOnStorage([]*files.DeleteFileReq{
				{Destination: fmt.Sprintf("images/house_models/%s", path.Base(existingImage.Url))},
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
				`INSERT INTO "house_model_images" ("filename", "url", "house_model_id") VALUES ($1, $2, $3);`,
				newImage.FileName, newImage.Url, b.req.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new image: %v", err)
			}
		}
	}

	return nil
}

func (b *updateHouseModelBuilder) updateTypeItem() error {
	// Retrieve existing items associated with the house model
	existingItems := make([]*houseModels.HouseModelTypeItem, 0)
	if err := b.db.Select(
		&existingItems,
		`SELECT "id", "room_type", "amount" FROM "house_model_type_items" WHERE "house_model_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing house_model_type_items: %v", err)
	}

	// Compare existing items with new items
	for _, existingItem := range existingItems {
		itemFound := false
		for _, newItem := range b.req.TypeItem {
			if existingItem.RoomType == newItem.RoomType && existingItem.Amount == newItem.Amount {
				itemFound = true
				break
			}
		}
		// If existing items not found in the new items, delete it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "house_model_type_items" WHERE "id" = $1;`,
				existingItem.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing item: %v", err)
			}
		}
	}

	// Insert new type items
	for _, newItem := range b.req.TypeItem {
		itemFound := false
		for _, existingItem := range existingItems {
			if newItem.RoomType == existingItem.RoomType && newItem.Amount == existingItem.Amount {
				itemFound = true
				break
			}
		}
		// If new type item not found in existing type items, insert it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`INSERT INTO "house_model_type_items" ("room_type", "amount", "house_model_id") VALUES ($1, $2, $3);`,
				newItem.RoomType, newItem.Amount, b.req.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new type item: %v", err)
			}
		}
	}
	return nil
}

func (b *updateHouseModelBuilder) updateHousePlan() error {
	// Retrieve existing plans associated with the house model
	existingPlans := make([]*houseModels.HouseModelPlan, 0)
	if err := b.db.Select(
		&existingPlans,
		`SELECT "id", "floor", "size" FROM "house_model_plans" WHERE "house_model_id" = $1;`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing house_model_plans items: %v", err)
	}

	// Compare existing plans with new plans
	for _, existingPlan := range existingPlans {
		planFound := false
		for _, newPlan := range b.req.HousePlan {
			if existingPlan.Id == newPlan.Id {
				planFound = true
				break
			}
		}
		// If existing plan not found in the new plans, delete it
		if !planFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "house_model_plans" WHERE "id" = $1;`,
				existingPlan.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing house_model_plans item: %v", err)
			}
		}
	}

	// Insert new plans
	for _, housePlan := range b.req.HousePlan {
		planFound := false
		for _, existingPlan := range existingPlans {
			if existingPlan.Id == housePlan.Id {
				planFound = true
				break
			}
		}
		// If new plan not found in existing plans, insert it
		if !planFound {
			// Insert the new plan into the house_model_plans table
			if err := b.tx.QueryRowContext(
				context.Background(),
				`INSERT INTO "house_model_plans" ("floor", "size", "house_model_id") VALUES ($1, $2, $3) RETURNING "id";`,
				housePlan.Floor, housePlan.Size, b.req.Id,
			).Scan(&housePlan.Id); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to insert new house_model_plans item: %v", err)
			}
		}
	}
	return nil
}

func (b *updateHouseModelBuilder) updateHousePlanItems() error {
	// Retrieve existing plan items associated with the house model
	existingPlanItems := make([]*houseModels.HouseModelPlanItem, 0)
	if err := b.db.Select(
		&existingPlanItems,
		`SELECT "id", "house_model_plan_id", "room_type", "amount" 
		FROM "house_model_plan_items" WHERE "house_model_plan_id" 
		IN (SELECT "id" FROM "house_model_plans" WHERE "house_model_id" = $1);`,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("failed to retrieve existing house_model_plan_items: %v", err)
	}

	// Compare existing plan items with new plan items
	for _, existingItem := range existingPlanItems {
		itemFound := false
		for _, housePlan := range b.req.HousePlan {
			for _, newItem := range housePlan.PlanItem {
				if existingItem.Id == newItem.Id {
					itemFound = true
					break
				}
			}
		}
		// If existing item not found in the new items, delete it
		if !itemFound {
			if _, err := b.tx.ExecContext(
				context.Background(),
				`DELETE FROM "house_model_plan_items" WHERE "id" = $1;`,
				existingItem.Id,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("failed to delete existing house_model_plan_item: %v", err)
			}
		}
	}

	// Insert new plan items
	for _, housePlan := range b.req.HousePlan {
		if housePlan.PlanItem == nil {
			continue
		}
		for _, planItem := range housePlan.PlanItem {
			itemFound := false
			for _, existingItem := range existingPlanItems {
				if existingItem.Id == planItem.Id {
					itemFound = true
					break
				}
			}
			// If new item not found in existing items, insert it
			if !itemFound {
				// Insert the new plan item into the house_model_plan_items table
				if _, err := b.tx.ExecContext(
					context.Background(),
					`INSERT INTO "house_model_plan_items" ("house_model_plan_id", "room_type", "amount") VALUES ($1, $2, $3);`,
					housePlan.Id, planItem.RoomType, planItem.Amount,
				); err != nil {
					b.tx.Rollback()
					return fmt.Errorf("failed to insert new house_model_plan_item: %v", err)
				}
			}
		}
	}

	return nil
}

func (b *updateHouseModelBuilder) updateImagesHousePlan() error {
	// Loop through each house plan
	for _, housePlan := range b.req.HousePlan {
		// If the plan does not have any images, skip to the next plan
		if housePlan.Images == nil {
			continue
		}

		// Retrieve existing images associated with the current house plan
		existingImages := make([]*entities.Image, 0)
		if err := b.db.Select(
			&existingImages,
			`SELECT "id", "filename", "url" FROM "house_model_plan_images" WHERE "house_model_plan_id" = $1;`,
			housePlan.Id,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("failed to retrieve existing images for house plan %d: %v", housePlan.Id, err)
		}

		// Compare existing images with new images
		for _, existingImage := range existingImages {
			imageFound := false
			for _, newImage := range housePlan.Images {
				if existingImage.FileName == newImage.FileName {
					imageFound = true
					break
				}
			}
			// If existing image not found in the new images, delete it
			if !imageFound {
				if _, err := b.tx.ExecContext(
					context.Background(),
					`DELETE FROM "house_model_plan_images" WHERE "id" = $1;`,
					existingImage.Id,
				); err != nil {
					b.tx.Rollback()
					return fmt.Errorf("failed to delete existing image: %v", err)
				}
				// Also delete the file from storage
				b.filesUsecases.DeleteFileOnStorage([]*files.DeleteFileReq{
					{Destination: fmt.Sprintf("images/house_model_plans/%s", path.Base(existingImage.Url))},
				})
			}
		}

		// Insert new images
		for _, newImage := range housePlan.Images {
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
					`INSERT INTO "house_model_plan_images" ("filename", "url", "house_model_plan_id") VALUES ($1, $2, $3);`,
					newImage.FileName, newImage.Url, housePlan.Id,
				); err != nil {
					b.tx.Rollback()
					return fmt.Errorf("failed to insert new image for house plan %d: %v", housePlan.Id, err)
				}
			}
		}
	}

	return nil
}


func (b *updateHouseModelBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateHouseModelBuilder) updateHouseModel() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update house model failed: %v", err)
	}
	return nil
}

func (b *updateHouseModelBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateHouseModelBuilder) getValues() []any         { return b.values }
func (b *updateHouseModelBuilder) getQuery() string         { return b.query }
func (b *updateHouseModelBuilder) setQuery(query string)    { b.query = query }
func (b *updateHouseModelBuilder) getImagesLen() int        { return len(b.req.Images) }

func (b *updateHouseModelBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

type updateHouseModelEngineer struct {
	builder IUpdateHouseModelBuilder
}

func UpdateHouseModelEngineer(b IUpdateHouseModelBuilder) *updateHouseModelEngineer {
	return &updateHouseModelEngineer{builder: b}
}

func (en *updateHouseModelEngineer) sumQueryFields() {
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

func (en *updateHouseModelEngineer) UpdateHouseModel() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.updateImagesHouseModel(); err != nil {
			return err
		}
	}

	if err := en.builder.updateTypeItem(); err != nil {
		return err
	}

	if err := en.builder.updateHousePlan(); err != nil {
		return err
	}

	if err := en.builder.updateHousePlanItems(); err != nil {
		return err
	}

	if err := en.builder.updateImagesHousePlan(); err != nil {
		return err
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}
