package interestsPatterns

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/interests"
)


type IInsertInterestBuidler interface {
	initTransaction() error
	insertInterest() error
	commit() error
	getInterestId() string
	// insertAttachment() error
}

type insertInterestBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *interests.Interest
}

func InsertInterestBuilder(db *sqlx.DB, req *interests.Interest) IInsertInterestBuidler {
	return &insertInterestBuilder{
		db:  db,
		req: req,
	}
}

type insertInterestEngineer struct {
	builder IInsertInterestBuidler
}

func (b *insertInterestBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *insertInterestBuilder) insertInterest() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "interests" (
		"bank_name",
		"interest_rate",
		"note",
		"display",
		"filename",
		"url"
	)
	VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.BankName,
		b.req.InterestRate,
		b.req.Note,
		b.req.Display,	
		b.req.FileName,
		b.req.Url,	
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert interest failed: %v", err)
	}
	return nil
}

// func (b *insertInterestBuilder) insertAttachment() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
// 	defer cancel()

// 	query := `
// 	INSERT INTO "interest_images" (
// 		"filename",
// 		"url",
// 		"interest_id"
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
// 		ctx,
// 		query,
// 		valueStack...,
// 	); err != nil {
// 		b.tx.Rollback()
// 		return fmt.Errorf("insert images failed: %v", err)
// 	}
// 	return nil
// }

func (b *insertInterestBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b *insertInterestBuilder) getInterestId() string {
	return strconv.Itoa(b.req.Id)
}

func InsertInterestEngineer(b IInsertInterestBuidler) *insertInterestEngineer {
	return &insertInterestEngineer{builder: b}
}

func (en *insertInterestEngineer) InsertInterest() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertInterest(); err != nil {
		return "", err
	}

	// if err := en.builder.insertAttachment(); err != nil {
	// 	return "", err
	// }
	
	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getInterestId(), nil
}
