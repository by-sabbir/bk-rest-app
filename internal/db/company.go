package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/by-sabbir/company-microservice-rest/internal/company"
	"github.com/google/uuid"
)

type CompanyRow struct {
	ID             string
	Name           string
	Description    sql.NullString
	TotalEmployees int
	IsRegistered   bool
	Type           string
}

func convertCompany(c *CompanyRow) company.Company {
	return company.Company{
		ID:             c.ID,
		Name:           c.Name,
		Description:    c.Description.String,
		TotalEmployees: c.TotalEmployees,
		IsRegistered:   c.IsRegistered,
		Type:           c.Type,
	}
}

func (d *DataBase) GetCompany(ctx context.Context, uuid string) (company.Company, error) {
	var cmpRow CompanyRow
	row := d.Client.QueryRowContext(
		ctx,
		`select id, name, description, total_employees, is_registered, type
		from company
		where id::text=$1`,
		uuid,
	)
	err := row.Scan(&cmpRow.ID, &cmpRow.Name, &cmpRow.Description,
		&cmpRow.TotalEmployees, &cmpRow.IsRegistered, &cmpRow.Type)
	if err != nil {
		return company.Company{}, fmt.Errorf("error fetching company from uuid: %+v", err)
	}
	return convertCompany(&cmpRow), nil
}

func (d *DataBase) PostCompany(ctx context.Context, cmp company.Company) (company.Company, error) {
	cmp.ID = uuid.NewString()

	if err := cmp.ScanType(); err != nil {
		return company.Company{}, err
	}

	postRow := CompanyRow{
		ID:             cmp.ID,
		Name:           cmp.Name,
		Description:    sql.NullString{String: cmp.Description, Valid: true},
		TotalEmployees: cmp.TotalEmployees,
		IsRegistered:   cmp.IsRegistered,
		Type:           cmp.Type,
	}

	if err := postRow.checkNull(); err != nil {
		return company.Company{}, err
	}

	qs := `insert into company
	(id, name, description, total_employees, is_registered, type)
	values
	($1, $2, $3, $4, $5, $6);`
	row, err := d.Client.QueryContext(
		ctx,
		qs,
		postRow.ID, postRow.Name, postRow.Description,
		postRow.TotalEmployees, postRow.IsRegistered, postRow.Type,
	)

	if err != nil {
		return company.Company{}, fmt.Errorf("error posting Company: %+v", err)
	}
	if err := row.Close(); err != nil {
		return company.Company{}, fmt.Errorf("could not close the row, %+v", err)
	}

	return convertCompany(&postRow), nil
}

func (d *DataBase) DeleteCompany(ctx context.Context, uuid string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`delete from company
		where id::text=$1`,
		uuid,
	)

	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) PartialUpdateCompany(
	ctx context.Context, id string, cmp company.Company,
) (company.Company, error) {

	currentRow, err := d.GetCompany(ctx, id)
	if err != nil {
		return company.Company{}, err
	}

	cmpRow := &CompanyRow{
		ID:             id,
		Description:    sql.NullString{String: currentRow.Description, Valid: true},
		Name:           currentRow.Name,
		TotalEmployees: currentRow.TotalEmployees,
		IsRegistered:   currentRow.IsRegistered,
		Type:           currentRow.Type,
	}

	// null check
	if (cmp.Name != currentRow.Name) && (cmp.Name != "") {
		cmpRow.Name = cmp.Name
	}
	if (cmp.Type != currentRow.Type) && cmp.Type != "" {
		if err := cmp.ScanType(); err != nil {
			return company.Company{}, err
		}
		cmpRow.Type = cmp.Type
	}
	if (cmp.TotalEmployees != currentRow.TotalEmployees) && (cmp.TotalEmployees > 0) {
		cmpRow.TotalEmployees = cmp.TotalEmployees
	}
	if cmp.IsRegistered != currentRow.IsRegistered {
		cmpRow.IsRegistered = cmp.IsRegistered
	}
	if (cmp.Description != currentRow.Description) && (cmp.Description != "") {
		cmpRow.Description = sql.NullString{String: cmp.Description, Valid: true}
	}

	if err := cmpRow.checkNull(); err != nil {
		return company.Company{}, err
	}
	row, err := d.Client.QueryContext(
		ctx,
		`update company set
		name=$1, description=$2, total_employees=$3,
		is_registered=$4, type=$5
		where id=$6`,
		cmpRow.Name, cmpRow.Description, cmpRow.TotalEmployees,
		cmpRow.IsRegistered, cmpRow.Type, cmpRow.ID,
	)
	if err != nil {
		return company.Company{}, fmt.Errorf("error updating Company: %w", err)
	}
	if err := row.Close(); err != nil {
		return company.Company{}, fmt.Errorf("error updating row: %w", err)
	}
	return convertCompany(cmpRow), nil
}

func (c *CompanyRow) checkNull() error {
	if c.ID == "" {
		return errors.New("id cannot be null")
	}
	if c.TotalEmployees <= 0 {
		return errors.New("total_employees must be greater than 0")
	}
	if c.Type == "" {
		return errors.New("type cannot be null or empty string")
	}
	if c.Name == "" {
		return errors.New("name cannot be null or empty string")
	}

	return nil
}
