package promotionsPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/promotions"
)

type IInsertPromotionBuilder interface {
	initTransaction() error
	insertPromotion() error
	insertPromotionFreeItem() error
	insertPromotionHouseModel() error
	insertAttachment() error
	commit() error
	getPromotionId() string
}

type insertPromotionBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *promotions.Promotion
}

func InsertPromotionBuilder(db *sqlx.DB, req *promotions.Promotion) IInsertPromotionBuilder {
	return &insertPromotionBuilder{
		db:  db,
		req: req,
	}
}

func (b *insertPromotionBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertPromotionBuilder) insertPromotion() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "promotions" (
		"index",
		"heading",
		"description",
		"start_date",
		"end_date",
		"display"
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Index,
		b.req.Heading,
		b.req.Description,
		b.req.StartDate,
		b.req.EndDate,
		b.req.Display,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert promotion failed: %v", err)
	}
	return nil
}

func (b *insertPromotionBuilder) insertPromotionFreeItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	query := `
	INSERT INTO "promotion_free_items" (
		"promotion_id",
		"description"
	)
	VALUES ($1, $2);`

	for _, freeItem := range b.req.FreeItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			freeItem.Description,
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert promotion_free_items failed: %v", err)
		}
	}
	return nil
}

func (b *insertPromotionBuilder) insertPromotionHouseModel() error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()
    query := `
        INSERT INTO "promotion_house_models" (
            "promotion_id",
            "house_model_id"
        )
        VALUES `
		values := make([]any, 0)
		lastIndex := 0
		for i, houseModel := range b.req.HouseModel {
			values = append(
				values,
				b.req.Id,
				houseModel.HouseModelId,
			)
	
			if i != len(b.req.HouseModel)-1 {
				query += fmt.Sprintf(`
				($%d, $%d),`, lastIndex+1, lastIndex+2)
			} else {
				query += fmt.Sprintf(`
				($%d, $%d);`, lastIndex+1, lastIndex+2)
			}
	
			lastIndex += 2
		}
	
		if _, err := b.tx.ExecContext(ctx, query, values...); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert promotion_house_models failed: %v", err)
		}
    return nil
}

func (b *insertPromotionBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

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
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert promotion images failed: %v", err)
	}
	return nil
}

func (b *insertPromotionBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertPromotionBuilder) getPromotionId() string {
	return strconv.Itoa(b.req.Id)
}

type insertPromotionEngineer struct {
	builder IInsertPromotionBuilder
}

func InsertPromotionEngineer(b IInsertPromotionBuilder) *insertPromotionEngineer {
	return &insertPromotionEngineer{builder: b}
}

func (en *insertPromotionEngineer) InsertPromotion() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertPromotion(); err != nil {
		return "", err
	}
	if err := en.builder.insertPromotionFreeItem(); err != nil {
		return "", err
	}
	if err := en.builder.insertPromotionHouseModel(); err != nil {
		return "", err
	}
	
	if err := en.builder.insertAttachment(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getPromotionId(), nil
}
