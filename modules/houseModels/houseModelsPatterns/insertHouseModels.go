package houseModelsPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/houseModels"
)

type IInsertHouseModelBuilder interface {
	initTransaction() error
	insertHouseModel() error
	insertHouseModelTypeItem() error
	insertAttachmentHouseModel() error
	insertHouseModelPlan() error
	insertHouseModelPlanItem() error
	insertAttachmentHouseModelPlan() error
	commit() error
	rollback() error
	getHouseModelId() string
	
}

func (b *insertHouseModelBuilder) rollback() error {
	if b.tx != nil {
		return b.tx.Rollback()
	}
	return nil
}

type insertHouseModelBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *houseModels.HouseModel
}

func InsertHouseModelBuilder(db *sqlx.DB, req *houseModels.HouseModel, planItems []*houseModels.HouseModelPlanItem) IInsertHouseModelBuilder {
	return &insertHouseModelBuilder{
		db:  db,
		req: req,
	}
}

type insertHouseModelEngineer struct {
	builder IInsertHouseModelBuilder
}

func (b *insertHouseModelBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertHouseModelBuilder) insertHouseModel() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "house_models" (
		"project_id",
		"name",
		"description",
		"link_video",
		"link_virtual_tour",
		"display",
		"index"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "id";
	`

	if err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.ProjectId,
		b.req.Name,
		b.req.Description,
		b.req.LinkVideo,
		b.req.LinkVirtualTour,
		b.req.Display,
		b.req.Index,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert house model failed: %v", err)
	}
	return nil
}

func (b *insertHouseModelBuilder) insertHouseModelTypeItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "house_model_type_items" (
		"house_model_id",
		"room_type",
		"amount"
	)
	VALUES ($1, $2, $3);
	`

	for _, typeItem := range b.req.TypeItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			typeItem.RoomType,
			typeItem.Amount,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert house model type item failed: %v", err)
		}
	}

	return nil
}

func (b *insertHouseModelBuilder) insertAttachmentHouseModel() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "house_model_images" (
		"filename",
		"url",
		"house_model_id"
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
		return fmt.Errorf("insert house model images failed: %v", err)
	}
	return nil
}

func (b *insertHouseModelBuilder) insertHouseModelPlan() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
    INSERT INTO "house_model_plans" (
        "house_model_id",
        "floor",
        "size"
    )
    VALUES ($1, $2, $3)
    RETURNING "id";
    `

	// Use QueryRowContext to execute the INSERT statement and retrieve the inserted ID
	for _, housePlan := range b.req.HousePlan {
		if err := b.tx.QueryRowContext(
			ctx,
			query,
			b.req.Id,
			housePlan.Floor,
			housePlan.Size,
		).Scan(&housePlan.Id); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert house model plans failed: %v", err)
		}
	}

	return nil
}

func (b *insertHouseModelBuilder) insertHouseModelPlanItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
    INSERT INTO "house_model_plan_items" (
        "house_model_plan_id",
        "room_type",
        "amount"
    )
    VALUES ($1, $2, $3);
    `

	for _, housePlan := range b.req.HousePlan {
		if housePlan.PlanItem == nil {
			continue
		}
		for _, planItem := range housePlan.PlanItem {
			if _, err := b.tx.ExecContext(
				ctx,
				query,
				housePlan.Id, // Use the plan item's house_model_plan_id
				planItem.RoomType,
				planItem.Amount,
			); err != nil {
				b.tx.Rollback()
				return fmt.Errorf("insert house model plan items failed: %v", err)
			}
		}
	}

	return nil
}

func (b *insertHouseModelBuilder) insertAttachmentHouseModelPlan() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "house_model_plan_images" (
		"filename",
		"url",
		"house_model_plan_id"
	)
	VALUES`

	valueStack := make([]interface{}, 0)
	var index int
	for _, housePlan := range b.req.HousePlan {
		if housePlan.Images == nil {
			continue
		}
		for _, image := range housePlan.Images {
			valueStack = append(valueStack,
				image.FileName,
				image.Url,
				housePlan.Id,
			)

			if index != 0 {
				query += ","
			}
			query += fmt.Sprintf("($%d, $%d, $%d)", index+1, index+2, index+3)
			index += 3
		}
	}

	if len(valueStack) == 0 {
		// No images to insert
		return nil
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert house model plan images failed: %v", err)
	}
	return nil
}

func (b *insertHouseModelBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertHouseModelBuilder) getHouseModelId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertProjectEngineer(b IInsertHouseModelBuilder) *insertHouseModelEngineer {
	return &insertHouseModelEngineer{builder: b}
}

func (en *insertHouseModelEngineer) InsertHouseModel() (string, error) {
	// houseModelReq := en.builder.GetHouseModelRequest()
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertHouseModel(); err != nil {
		return "", err
	}
	if err := en.builder.insertHouseModelTypeItem(); err != nil {
		return "", err
	}
	if err := en.builder.insertAttachmentHouseModel(); err != nil {
		return "", err
	}
	if err := en.builder.insertHouseModelPlan(); err != nil {
		return "", err
	}
	if err := en.builder.insertHouseModelPlanItem(); err != nil {
		en.builder.rollback()
		return "", err
	}

	if err := en.builder.insertAttachmentHouseModelPlan(); err != nil {
		en.builder.rollback()
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getHouseModelId(), nil
}
