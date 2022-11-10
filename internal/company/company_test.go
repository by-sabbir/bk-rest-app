package company

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCompanyService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("test get company by id", func(t *testing.T) {
		companyStoreMock := NewMockStore(mockCtrl)

		company := Company{
			ID:             "431A5348-5277-4E34-A71B-CE764398A261",
			Name:           "ACME Corp",
			Description:    "Lorem Ipsum Dolor",
			TotalEmployees: 100,
			IsRegistered:   true,
			Type:           CompanyType[0],
		}
		ctx := context.Background()
		companyStoreMock.
			EXPECT().
			GetCompany(ctx, company.ID).Return(company, nil)

		companyService := NewService(companyStoreMock)

		got, err := companyService.GetCompany(ctx, company.ID)

		assert.NoError(t, err)
		assert.Equal(t, company, got)
	})

	t.Run("test create company", func(t *testing.T) {
		companyStoreMock := NewMockStore(mockCtrl)

		company := Company{
			ID:             "431A5349-5277-4E34-A71B-CE764398A261",
			Name:           "ACME Corp",
			Description:    "Lorem Ipsum Dolor",
			TotalEmployees: 100,
			IsRegistered:   true,
			Type:           CompanyType[0],
		}
		ctx := context.Background()
		companyStoreMock.
			EXPECT().
			PostCompany(ctx, company).Return(company, nil)

		companyService := NewService(companyStoreMock)

		got, err := companyService.PostCompany(ctx, company)

		assert.NoError(t, err)
		assert.Equal(t, company, got)
	})

	t.Run("test delete company by id", func(t *testing.T) {
		companyStoreMock := NewMockStore(mockCtrl)

		company := Company{
			ID:             "431A5349-5277-4E34-A71B-CE764398A261",
			Name:           "ACME Corp",
			Description:    "Lorem Ipsum Dolor",
			TotalEmployees: 100,
			IsRegistered:   true,
			Type:           CompanyType[0],
		}
		ctx := context.Background()
		companyStoreMock.
			EXPECT().
			DeleteCompany(ctx, company.ID).Return(nil)

		companyService := NewService(companyStoreMock)

		err := companyService.DeleteCompany(ctx, company.ID)

		assert.NoError(t, err)
	})

	// t.Run("test partial update company by id", func(t *testing.T) {
	// 	companyStoreMock := NewMockStore(mockCtrl)

	// 	company := Company{
	// 		ID:             "431A5349-5277-4E34-A71B-CE764398A261",
	// 		Name:           "ACME Corp",
	// 		TotalEmployees: 100,
	// 		IsRegistered:   true,
	// 		Type:           CompanyType[0],
	// 	}
	// 	ctx := context.Background()
	// 	companyStoreMock.
	// 		EXPECT().
	// 		PartialUpdateCompany(ctx, company.ID, company).Return(company, nil)

	// 	companyService := NewService(companyStoreMock)

	// 	got, err := companyService.PartialUpdateCompany(ctx, company.ID, company)

	// 	assert.NoError(t, err)
	// 	assert.Equal(t, company, got)
	// })

}

func TestScanType(t *testing.T) {

	t.Run("test company type exist", func(t *testing.T) {
		companyWithCorrectType := &Company{
			ID:             "431A5349-5277-4E34-A71B-CE764398A261",
			Name:           "ACME Corp",
			Description:    "Lorem Ipsum Dolor",
			TotalEmployees: 100,
			IsRegistered:   true,
			Type:           CompanyType[0],
		}

		err := companyWithCorrectType.ScanType()
		assert.NoError(t, err)
	})
	t.Run("test company type does not exist", func(t *testing.T) {
		companyWithInCorrectType := &Company{
			ID:             "431A5349-5277-4E34-A71B-CE764398A261",
			Name:           "ACME Corp",
			Description:    "Lorem Ipsum Dolor",
			TotalEmployees: 100,
			IsRegistered:   true,
			Type:           "Does Not Exist",
		}

		err := companyWithInCorrectType.ScanType()
		assert.ErrorIs(t, err, ErrTypeNotFound)
	})
}
