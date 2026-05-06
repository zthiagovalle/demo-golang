package repositories_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/infra/repositories"
)

var productsDatasets = []string{"clear-data.sql", "products-insert.sql"}

func newProductsRepository() *repositories.ProductsPostgresRepository {
	return repositories.NewProductsPostgresRepository(testPool)
}

func TestProductsPostgresRepository_Insert(t *testing.T) {
	loadDatasets(t, "clear-data.sql")
	repo := newProductsRepository()

	now := time.Now().UTC().Truncate(time.Second)
	p := &models.Product{
		ID:          uuid.NewString(),
		Name:        "Espresso",
		Description: "Doppio",
		PriceCents:  3990,
		Status:      enums.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	require.NoError(t, repo.Insert(testCtx, p))

	got, err := repo.FindByID(testCtx, p.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, p.Name, got.Name)
	assert.Equal(t, p.Description, got.Description)
	assert.Equal(t, p.PriceCents, got.PriceCents)
	assert.Equal(t, p.Status, got.Status)
}

func TestProductsPostgresRepository_FindByID_Found(t *testing.T) {
	loadDatasets(t, productsDatasets...)
	repo := newProductsRepository()

	got, err := repo.FindByID(testCtx, "11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Cafe Especial", got.Name)
	assert.Equal(t, int64(4990), got.PriceCents)
	assert.Equal(t, enums.ProductStatusActive, got.Status)
}

func TestProductsPostgresRepository_FindByID_NotFound(t *testing.T) {
	loadDatasets(t, "clear-data.sql")
	repo := newProductsRepository()

	got, err := repo.FindByID(testCtx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestProductsPostgresRepository_Update(t *testing.T) {
	loadDatasets(t, productsDatasets...)
	repo := newProductsRepository()

	id := "11111111-1111-1111-1111-111111111111"
	now := time.Now().UTC().Truncate(time.Second)
	updated := &models.Product{
		ID:          id,
		Name:        "Cafe Especial Reserva",
		Description: "Lote unico",
		PriceCents:  6990,
		UpdatedAt:   now,
	}
	require.NoError(t, repo.Update(testCtx, updated))

	got, err := repo.FindByID(testCtx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Cafe Especial Reserva", got.Name)
	assert.Equal(t, "Lote unico", got.Description)
	assert.Equal(t, int64(6990), got.PriceCents)
}

func TestProductsPostgresRepository_Delete(t *testing.T) {
	loadDatasets(t, productsDatasets...)
	repo := newProductsRepository()

	id := "22222222-2222-2222-2222-222222222222"
	require.NoError(t, repo.Delete(testCtx, id))

	got, err := repo.FindByID(testCtx, id)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestProductsPostgresRepository_UpdateStatus(t *testing.T) {
	loadDatasets(t, productsDatasets...)
	repo := newProductsRepository()

	id := "11111111-1111-1111-1111-111111111111"
	require.NoError(t, repo.UpdateStatus(testCtx, id, string(enums.ProductStatusInactive)))

	got, err := repo.FindByID(testCtx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, enums.ProductStatusInactive, got.Status)
}

func TestProductsPostgresRepository_Paginate(t *testing.T) {
	loadDatasets(t, productsDatasets...)
	repo := newProductsRepository()

	t.Run("first page returns ordered by created_at desc", func(t *testing.T) {
		page, err := repo.Paginate(testCtx, &models.ProductPageParams{Page: 1, PageSize: 2})
		require.NoError(t, err)
		require.NotNil(t, page)
		assert.Equal(t, int64(3), page.TotalItems)
		require.Len(t, page.Items, 2)
		assert.Equal(t, "Mate", page.Items[0].Name)      // newest first
		assert.Equal(t, "Cha Verde", page.Items[1].Name)
	})

	t.Run("second page returns remaining item", func(t *testing.T) {
		page, err := repo.Paginate(testCtx, &models.ProductPageParams{Page: 2, PageSize: 2})
		require.NoError(t, err)
		require.NotNil(t, page)
		assert.Equal(t, int64(3), page.TotalItems)
		require.Len(t, page.Items, 1)
		assert.Equal(t, "Cafe Especial", page.Items[0].Name)
	})

	t.Run("page beyond total returns empty slice", func(t *testing.T) {
		page, err := repo.Paginate(testCtx, &models.ProductPageParams{Page: 99, PageSize: 10})
		require.NoError(t, err)
		require.NotNil(t, page)
		assert.Equal(t, int64(3), page.TotalItems)
		assert.Empty(t, page.Items)
	})
}
