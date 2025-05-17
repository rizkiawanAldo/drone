package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for database operations
type Repository interface {
	// Estate methods
	CreateEstate(ctx context.Context, width, length int) (uuid.UUID, error)
	GetEstate(ctx context.Context, id uuid.UUID) (width, length int, err error)
	
	// Tree methods
	CreateTree(ctx context.Context, estateID uuid.UUID, x, y, height int) (uuid.UUID, error)
	GetTrees(ctx context.Context, estateID uuid.UUID) ([]Tree, error)
}

// Tree represents a tree in the database
type Tree struct {
	ID     uuid.UUID
	EstateID uuid.UUID
	X      int
	Y      int
	Height int
}

// Stats represents tree statistics for an estate
type Stats struct {
	Count       int
	MaxHeight   int
	MinHeight   int
	MedianHeight int
}

// repository implements the Repository interface
type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new repository with the given database connection
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

// CreateEstate creates a new estate in the database
func (r *repository) CreateEstate(ctx context.Context, width, length int) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO estates (width, length) VALUES ($1, $2) RETURNING id",
		width, length).Scan(&id)
	return id, err
}

// GetEstate retrieves an estate from the database by ID
func (r *repository) GetEstate(ctx context.Context, id uuid.UUID) (width, length int, err error) {
	err = r.db.QueryRow(ctx,
		"SELECT width, length FROM estates WHERE id = $1",
		id).Scan(&width, &length)
	return
}

// CreateTree creates a new tree in the database
func (r *repository) CreateTree(ctx context.Context, estateID uuid.UUID, x, y, height int) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		"INSERT INTO trees (estate_id, x, y, height) VALUES ($1, $2, $3, $4) RETURNING id",
		estateID, x, y, height).Scan(&id)
	return id, err
}

// GetTrees retrieves all trees for an estate from the database
func (r *repository) GetTrees(ctx context.Context, estateID uuid.UUID) ([]Tree, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = $1",
		estateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trees []Tree
	for rows.Next() {
		var tree Tree
		if err := rows.Scan(&tree.ID, &tree.EstateID, &tree.X, &tree.Y, &tree.Height); err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trees, nil
} 