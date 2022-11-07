package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/by-sabbir/company-microservice-rest/internal/company"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestCompanyDB(t *testing.T) {
	t.Setenv("DB_HOST", "0.0.0.0")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USERNAME", "bktest")
	t.Setenv("DB_PASSWORD", "hello")
	t.Setenv("DB_NAME", "postgres")
	t.Setenv("SSL_MODE", "disable")
	rand.Seed(time.Now().UnixNano())

	db, err := NewDatabase()
	assert.NoError(t, err)

	t.Run("test database migration", func(t *testing.T) {
		err := db.MigrateDB()
		assert.NoError(t, err)
	})

	t.Run("test no change at migration", func(t *testing.T) {
		err := db.MigrateDB()
		assert.NoError(t, err)
	})

	initCmp := company.Company{
		Name:           RandStringBytes(6),
		Description:    "Lorem Ipsum Dolor Sit",
		TotalEmployees: 120,
		IsRegistered:   true,
		Type:           company.CompanyType[1],
	}

	var id string
	t.Run("test create company", func(t *testing.T) {
		cmp, err := db.PostCompany(context.Background(), initCmp)
		assert.NoError(t, err)
		assert.NotEmpty(t, cmp.ID)
		id = cmp.ID
	})

	t.Run("test get company", func(t *testing.T) {
		gotCmt, err := db.GetCompany(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, initCmp.TotalEmployees, gotCmt.TotalEmployees)
	})

	t.Run("test get company with wrong id", func(t *testing.T) {
		_, err := db.GetCompany(context.Background(), uuid.NewString())
		assert.Error(t, err)
	})

	t.Run("test partial update", func(t *testing.T) {
		updatedComment, errUpdate := db.PartialUpdateCompany(context.Background(), id, company.Company{
			Name:           initCmp.Name,
			Description:    initCmp.Description,
			TotalEmployees: 150,
			IsRegistered:   false,
			Type:           initCmp.Type,
		})
		assert.NoError(t, errUpdate)
		assert.Equal(t, id, updatedComment.ID)
		assert.NotEqual(t, initCmp.TotalEmployees, updatedComment.TotalEmployees)
	})

	t.Run("test missing required value", func(t *testing.T) {
		newCmp := company.Company{
			Name:           RandStringBytes(6),
			Description:    "Lorem Ipsum Dolor Sit",
			TotalEmployees: 120,
			IsRegistered:   true,
		}
		_, err := db.PostCompany(context.Background(), newCmp)
		assert.Error(t, err)
		newCmp = company.Company{
			Name:         RandStringBytes(6),
			Description:  "Lorem Ipsum Dolor Sit",
			Type:         company.CompanyType[1],
			IsRegistered: true,
		}
		_, err = db.PostCompany(context.Background(), newCmp)
		assert.Error(t, err)
	})

	t.Run("test duplicate entry fails", func(t *testing.T) {
		_, err := db.PostCompany(context.Background(), initCmp)
		assert.Error(t, err)
	})

	t.Run("test delete entry", func(t *testing.T) {
		err := db.DeleteCompany(context.Background(), id)
		assert.NoError(t, err)
	})
}
