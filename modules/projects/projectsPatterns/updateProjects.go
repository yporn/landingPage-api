package projectsPatterns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/projects"
)

type IUpdateProjectBuilder interface {
	initTransaction() error
	initQuery()
	updateQuery()
	insertHouseTypeItem() error
	getOldHouseTypeItem() []*projects.ProjectHouseTypeItem
	deleteOldHouseTypeItem() error
	insertDescAreaItem() error
	getOldDescAreaItem() []*projects.ProjectDescAreaItem
	deleteOldDescAreaItem() error
	insertFacilityItem() error
	getOldFacilityItem() []*projects.ProjectFacilityItem
	deleteOldFacilityItem() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProject() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	getHouseTypeItemLen() int
	getDescAreaItemLen() int
	getFacilityItemLen() int
	commit() error
}

type updateProjectBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *projects.Project
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
}

func UpdateProjectBuilder(db *sqlx.DB, req *projects.Project, filesUsecases filesUsecases.IFilesUsecase) IUpdateProjectBuilder {
	return &updateProjectBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryFields:   make([]string, 0),
		values:        make([]any, 0),
	}
}

type updateProjectEngineer struct {
	builder IUpdateProjectBuilder
}

func (b *updateProjectBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateProjectBuilder) initQuery() {
	b.query += `
	UPDATE "projects" SET`
}

func (b *updateProjectBuilder) updateQuery() {
	setStatements := []string{}
	if b.req.Name != "" {
		b.values = append(b.values, b.req.Name)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"name" = $%d`, b.lastStackIndex))
	}

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

	if b.req.Text != "" {
		b.values = append(b.values, b.req.Text)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"text" = $%d`, b.lastStackIndex))
	}

	if b.req.Location != "" {
		b.values = append(b.values, b.req.Location)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"location" = $%d`, b.lastStackIndex))
	}

	if b.req.Price != 0 {
		b.values = append(b.values, b.req.Price)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"price" = $%d`, b.lastStackIndex))
	}

	if b.req.StatusProject != "" {
		b.values = append(b.values, b.req.StatusProject)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"status_project" = $%d`, b.lastStackIndex))
	}

	if b.req.TypeProject != "" {
		b.values = append(b.values, b.req.TypeProject)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"type_project" = $%d`, b.lastStackIndex))
	}

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"description" = $%d`, b.lastStackIndex))
	}

	if b.req.NameFacebook != "" {
		b.values = append(b.values, b.req.NameFacebook)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"name_facebook" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkFacebook != "" {
		b.values = append(b.values, b.req.LinkFacebook)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_facebook" = $%d`, b.lastStackIndex))
	}

	if b.req.Tel != "" {
		b.values = append(b.values, b.req.Tel)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"tel" = $%d`, b.lastStackIndex))
	}

	if b.req.Address != "" {
		b.values = append(b.values, b.req.Address)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"address" = $%d`, b.lastStackIndex))
	}

	if b.req.LinkLocation != "" {
		b.values = append(b.values, b.req.LinkLocation)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"link_location" = $%d`, b.lastStackIndex))
	}

	if b.req.Display != "" {
		b.values = append(b.values, b.req.Display)
		b.lastStackIndex = len(b.values)

		setStatements = append(setStatements, fmt.Sprintf(`"display" = $%d`, b.lastStackIndex))
	}

	b.query += strings.Join(setStatements, ", ")
}

func (b *updateProjectBuilder) insertHouseTypeItem() error {
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

func (b *updateProjectBuilder) getOldHouseTypeItem() []*projects.ProjectHouseTypeItem {
	query := `
	SELECT
		"id",
		"name"
	FROM "project_house_type_items"
	WHERE "project_id" = $1;
	`
	houseTypeItem := make([]*projects.ProjectHouseTypeItem, 0)
	if err := b.db.Select(
		&houseTypeItem,
		query,
		b.req.Id,
	); err != nil {
		return make([]*projects.ProjectHouseTypeItem, 0)
	}
	return houseTypeItem
}

func (b *updateProjectBuilder) deleteOldHouseTypeItem() error {
	query := `
	DELETE FROM "project_house_type_items"
	WHERE "project_id" = $1;
	`
	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete house type items failed: %v", err)
	}
	return nil
}

func (b *updateProjectBuilder) insertDescAreaItem() error {
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
			return fmt.Errorf("insert project desc area item failed: %v", err)
		}
	}

	return nil
}

func (b *updateProjectBuilder) getOldDescAreaItem() []*projects.ProjectDescAreaItem {
	query := `
	SELECT
		"id",
		"item",
		"amount",
		"unit"
	FROM "project_desc_area_items"
	WHERE "project_id" = $1;
	`
	descAreaItem := make([]*projects.ProjectDescAreaItem, 0)
	if err := b.db.Select(
		&descAreaItem,
		query,
		b.req.Id,
	); err != nil {
		return make([]*projects.ProjectDescAreaItem, 0)
	}
	return descAreaItem
}

func (b *updateProjectBuilder) deleteOldDescAreaItem() error {
	query := `
	DELETE FROM "project_desc_area_items"
	WHERE "project_id" = $1;
	`
	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete desc area items failed: %v", err)
	}
	return nil
}

func (b *updateProjectBuilder) insertFacilityItem() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "project_facility_items" (
		"project_id",
		"item"
	)
	VALUES ($1, $2);
	`

	for _, facilityItem := range b.req.FacilityItem {
		if _, err := b.tx.ExecContext(
			ctx,
			query,
			b.req.Id,
			facilityItem.Item,
			
		); err != nil {
			b.tx.Rollback()
			return fmt.Errorf("insert facilities failed: %v", err)
		}
	}

	return nil
}

func (b *updateProjectBuilder) getOldFacilityItem() []*projects.ProjectFacilityItem {
	query := `
	SELECT
		"id",
		"item"
	FROM "project_facility_items"
	WHERE "project_id" = $1;
	`
	facilityItem := make([]*projects.ProjectFacilityItem, 0)
	if err := b.db.Select(
		&facilityItem,
		query,
		b.req.Id,
	); err != nil {
		return make([]*projects.ProjectFacilityItem, 0)
	}
	return facilityItem
}

func (b *updateProjectBuilder) deleteOldFacilityItem() error {
	query := `
	DELETE FROM "project_facility_items"
	WHERE "project_id" = $1;
	`
	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete facilities items failed: %v", err)
	}
	return nil
}

func (b *updateProjectBuilder) insertImages() error {
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
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update images failed: %v", err)
	}
	return nil
}

func (b *updateProjectBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "project_images"
	WHERE "project_id" = $1;`

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

func (b *updateProjectBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "project_images"
	WHERE "project_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/projects/%s", img.FileName),
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

func (b *updateProjectBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateProjectBuilder) updateProject() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update project failed: %v", err)
	}
	return nil
}

func (b *updateProjectBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateProjectBuilder) getValues() []any         { return b.values }
func (b *updateProjectBuilder) getQuery() string         { return b.query }
func (b *updateProjectBuilder) setQuery(query string)    { b.query = query }
func (b *updateProjectBuilder) getImagesLen() int        { return len(b.req.Images) }
func (b *updateProjectBuilder) getHouseTypeItemLen() int { return len(b.req.HouseTypeItem) }
func (b *updateProjectBuilder) getDescAreaItemLen() int { return len(b.req.DescAreaItem) }
func (b *updateProjectBuilder) getFacilityItemLen() int { return len(b.req.FacilityItem) }
func (b *updateProjectBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateProjectEngineer(b IUpdateProjectBuilder) *updateProjectEngineer {
	return &updateProjectEngineer{builder: b}
}

func (en *updateProjectEngineer) sumQueryFields() {
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

func (en *updateProjectEngineer) UpdateProject() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	if err := en.builder.updateProject(); err != nil {
		return err
	}

	if en.builder.getHouseTypeItemLen() > 0 {
		if err := en.builder.deleteOldHouseTypeItem(); err != nil {
			return err
		}
		if err := en.builder.insertHouseTypeItem(); err != nil {
			return err
		}
	}

	if en.builder.getDescAreaItemLen() > 0 {
		if err := en.builder.deleteOldDescAreaItem(); err != nil {
			return err
		}
		if err := en.builder.insertDescAreaItem(); err != nil {
			return err
		}
	}

	if en.builder.getFacilityItemLen() > 0 {
		if err := en.builder.deleteOldFacilityItem(); err != nil {
			return err
		}
		if err := en.builder.insertFacilityItem(); err != nil {
			return err
		}
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
