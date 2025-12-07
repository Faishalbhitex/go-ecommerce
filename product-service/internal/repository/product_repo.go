package repository

import (
	"context"
	"database/sql"
	"errors"
	"product-service/internal/models"
	"strconv"
	"strings"
)

var ErrNotFound = errors.New("product not found")

type ProductRepository interface {
	Create(ctx context.Context, p *models.Product) error
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	List(ctx context.Context) ([]*models.Product, error)
	ListPaged(ctx context.Context, limit, offset int) ([]*models.Product, error)
	Search(ctx context.Context, q string) ([]*models.Product, error)
	Update(ctx context.Context, p *models.Product) (*models.Product, error)
	Patch(ctx context.Context, id int64, patch *models.ProductPatch) (*models.Product, error)
	Delete(ctx context.Context, id int64) error
}

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// // @REPO:CREATE-BEGIN
func (r *PostgresProductRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
	INSERT INTO products (name, description, price, qty, category, create_at, update_at)
	VALUES ($1, $2, $3, $4, $5, now(), now())
	RETURNING id, create_at, update_at
	`
	return r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price, p.Qty, p.Category).Scan(&p.ID, &p.CreateAt, &p.UpdateAt)
}

//// @REPO:CREATE-END

// // @REPO:READ-BEGIN
func (r *PostgresProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `
	SELECT id, name, description, price, qty, category, create_at, update_at
	FROM products
	WHERE id=$1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	p := &models.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Qty, &p.Category, &p.CreateAt, &p.UpdateAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresProductRepository) List(ctx context.Context) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, qty, category, create_at, update_at
		FROM products
		ORDER BY id DESC
		LIMIT 100
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Qty, &p.Category, &p.CreateAt, &p.UpdateAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *PostgresProductRepository) ListPaged(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, qty, category, create_at, update_at
		FROM products ORDER BY id DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Qty, &p.Category, &p.CreateAt, &p.UpdateAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PostgresProductRepository) Search(ctx context.Context, q string) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, qty, category, create_at, update_at
		FROM products
		WHERE name ILIKE $1 OR category ILIKE $1
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Qty, &p.Category, &p.CreateAt, &p.UpdateAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

//// @REPO:READ-END

// // @REPO:UPDATE-BEGIN
func (r *PostgresProductRepository) Update(ctx context.Context, p *models.Product) (*models.Product, error) {
	query := `
		UPDATE products
		SET name=$1, description=$2, price=$3, qty=$4, category=$5, update_at=now()
		WHERE id=$6
		RETURNING id, name, description, price, qty, category, create_at, update_at
	`
	var updated models.Product
	err := r.db.QueryRowContext(ctx, query,
		p.Name,
		p.Description,
		p.Price,
		p.Qty,
		p.Category,
		p.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.Price,
		&updated.Qty,
		&updated.Category,
		&updated.CreateAt,
		&updated.UpdateAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &updated, nil
}

func (r *PostgresProductRepository) Patch(ctx context.Context, id int64, patch *models.ProductPatch) (*models.Product, error) {
	parts := []string{}
	args := []interface{}{}
	i := 1

	if patch.Name != nil {
		parts = append(parts, "name=$"+strconv.Itoa(i))
		args = append(args, *patch.Name)
		i++
	}
	if patch.Description != nil {
		parts = append(parts, "description=$"+strconv.Itoa(i))
		args = append(args, *patch.Description)
		i++
	}
	if patch.Price != nil {
		parts = append(parts, "price=$"+strconv.Itoa(i))
		args = append(args, *patch.Price)
		i++
	}
	if patch.Qty != nil {
		parts = append(parts, "qty=$"+strconv.Itoa(i))
		args = append(args, *patch.Qty)
		i++
	}
	if patch.Category != nil {
		parts = append(parts, "category=$"+strconv.Itoa(i))
		args = append(args, *patch.Category)
		i++
	}

	if len(parts) == 0 {
		parts = append(parts, "update_at=now()")
	}

	query := `
		UPDATE products SET ` + strings.Join(parts, ", ") + `, update_at=now()
		WHERE id=$` + strconv.Itoa(i) + `
		RETURNING id, name, description, price, qty, category, create_at, update_at
	`

	args = append(args, id)

	var updated models.Product
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.Price,
		&updated.Qty,
		&updated.Category,
		&updated.CreateAt,
		&updated.UpdateAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &updated, nil
}

//// @REPO:UPDATE-END

// // @REPO:DELETE-BEGIN
func (r *PostgresProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id=$1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

//// @REPO:DELETE-END
