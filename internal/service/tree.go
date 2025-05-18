package service

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateTree implements the TreeService.CreateTree method
func (s *service) CreateTree(ctx context.Context, estateID uuid.UUID, x, y, height int) (uuid.UUID, error) {
	// Validate estate exists
	width, length, err := s.repo.GetEstate(ctx, estateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errors.New("estate not found")
		}
		return uuid.Nil, err
	}

	// Validate coordinates are within estate bounds
	if x < 1 || x > width || y < 1 || y > length {
		return uuid.Nil, errors.New("tree coordinates outside estate boundaries")
	}

	// Validate height
	if height < 1 || height > 30 {
		return uuid.Nil, errors.New("invalid tree height")
	}

	// The database has a unique constraint on (estate_id, x, y) so if there's already a tree
	// at this location, the repository layer will return an error

	return s.repo.CreateTree(ctx, estateID, x, y, height)
}

// GetTreeStats implements the TreeService.GetTreeStats method
func (s *service) GetTreeStats(ctx context.Context, estateID uuid.UUID) (count, maxHeight, minHeight, medianHeight int, err error) {
	// Check if estate exists
	_, _, err = s.repo.GetEstate(ctx, estateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, 0, 0, errors.New("estate not found")
		}
		return 0, 0, 0, 0, err
	}

	// Get all trees for the estate
	trees, err := s.repo.GetTrees(ctx, estateID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// If no trees, return all zeros as per requirements
	count = len(trees)
	if count == 0 {
		return 0, 0, 0, 0, nil
	}

	// Calculate max and min height
	maxHeight = trees[0].Height
	minHeight = trees[0].Height
	
	// Collect heights for median calculation
	heights := make([]int, count)
	for i, tree := range trees {
		heights[i] = tree.Height
		
		if tree.Height > maxHeight {
			maxHeight = tree.Height
		}
		if tree.Height < minHeight {
			minHeight = tree.Height
		}
	}

	// Calculate median height
	sort.Ints(heights)
	if count%2 == 0 {
		// Even number of trees, average the middle two
		medianHeight = (heights[count/2-1] + heights[count/2]) / 2
	} else {
		// Odd number of trees, take the middle one
		medianHeight = heights[count/2]
	}

	return count, maxHeight, minHeight, medianHeight, nil
} 