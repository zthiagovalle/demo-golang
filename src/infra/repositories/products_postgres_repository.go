package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
)

const (
	productsBaseQuery = `
		SELECT
			p.id,
			p.name,
			p.description,
			p.price_cents,
			p.status,
			p.created_at,
			p.updated_at
		FROM products p`

	insertProductQuery = `
		INSERT INTO products (id, name, description, price_cents, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`

	updateProductQuery = `
		UPDATE products
		   SET name = $2, description = $3, price_cents = $4, updated_at = $5
		 WHERE id = $1;`

	deleteProductQuery = `DELETE FROM products WHERE id = $1;`

	updateProductStatusQuery = `UPDATE products SET status = $2, updated_at = NOW() WHERE id = $1;`

	findProductByIDQuery = productsBaseQuery + ` WHERE p.id = $1;`

	countProductsQuery = `SELECT COUNT(*) FROM products;`

	findAllPaginatedProductsQuery = productsBaseQuery + `
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2;`
)

type ProductsPostgresRepository struct {
	db *pgxpool.Pool
}

func NewProductsPostgresRepository(db *pgxpool.Pool) *ProductsPostgresRepository {
	return &ProductsPostgresRepository{db: db}
}

func (r *ProductsPostgresRepository) Insert(ctx context.Context, p *models.Product) error {
	_, err := r.db.Exec(ctx, insertProductQuery,
		p.ID, p.Name, p.Description, p.PriceCents, string(p.Status), p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *ProductsPostgresRepository) Update(ctx context.Context, p *models.Product) error {
	_, err := r.db.Exec(ctx, updateProductQuery,
		p.ID, p.Name, p.Description, p.PriceCents, p.UpdatedAt)
	return err
}

func (r *ProductsPostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteProductQuery, id)
	return err
}

func (r *ProductsPostgresRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	row := r.db.QueryRow(ctx, findProductByIDQuery, id)

	var p models.Product
	var status string
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.PriceCents, &status, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Status = enums.ProductStatus(status)
	return &p, nil
}

func (r *ProductsPostgresRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, updateProductStatusQuery, id, status)
	return err
}

func (r *ProductsPostgresRepository) Paginate(ctx context.Context, params *models.ProductPageParams) (*models.ProductPage, error) {
	offset := (params.Page - 1) * params.PageSize

	var total int64
	if err := r.db.QueryRow(ctx, countProductsQuery).Scan(&total); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, findAllPaginatedProductsQuery, params.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var status string
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PriceCents, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Status = enums.ProductStatus(status)
		items = append(items, p)
	}

	return &models.ProductPage{
		Items:      items,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: total,
	}, nil
}
