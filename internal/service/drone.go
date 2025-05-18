package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"drone/internal/repository"
)

// CalculateDronePath implements the DroneService.CalculateDronePath method
func (s *service) CalculateDronePath(ctx context.Context, estateID uuid.UUID) (distance int, err error) {
	// First check if estate exists
	width, length, err := s.repo.GetEstate(ctx, estateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("estate not found")
		}
		return 0, err
	}

	// Get all trees for the estate
	trees, err := s.repo.GetTrees(ctx, estateID)
	if err != nil {
		return 0, err
	}

	// Create a map for quick tree lookup by coordinates
	treeMap := make(map[string]repository.Tree)
	for _, tree := range trees {
		key := coordKey(tree.X, tree.Y)
		treeMap[key] = tree
	}

	// Calculate the drone path
	return calculateDroneTravelDistance(width, length, treeMap), nil
}

// CalculateDronePathWithRest implements the DroneService.CalculateDronePathWithRest method
func (s *service) CalculateDronePathWithRest(ctx context.Context, estateID uuid.UUID, maxDistance int) (distance int, restX, restY int, err error) {
	// First check if estate exists
	width, length, err := s.repo.GetEstate(ctx, estateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, 0, errors.New("estate not found")
		}
		return 0, 0, 0, err
	}

	// Get all trees for the estate
	trees, err := s.repo.GetTrees(ctx, estateID)
	if err != nil {
		return 0, 0, 0, err
	}

	// Create a map for quick tree lookup by coordinates
	treeMap := make(map[string]repository.Tree)
	for _, tree := range trees {
		key := coordKey(tree.X, tree.Y)
		treeMap[key] = tree
	}

	// Calculate drone path with rest point
	totalDistance, restPos := calculateDronePathWithRest(width, length, treeMap, maxDistance)
	return totalDistance, restPos.x, restPos.y, nil
}

// Helper function to create a key for the tree map
func coordKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// Position represents a 3D position (x, y, z)
type position struct {
	x, y, z int
}

// calculateDroneTravelDistance calculates the total distance the drone travels
func calculateDroneTravelDistance(width, length int, treeMap map[string]repository.Tree) int {
	totalDistance := 0
	
	// Start at ground level at the southwestern-most plot (1,1)
	currentPos := position{x: 1, y: 1, z: 0}
	
	// The drone travels in a zigzag pattern from south to north
	for y := 1; y <= length; y++ {
		// For even-numbered rows, go from west to east
		if y%2 == 1 {
			for x := 1; x <= width; x++ {
				totalDistance += visitPlot(x, y, &currentPos, treeMap)
			}
		} else { // For odd-numbered rows, go from east to west
			for x := width; x >= 1; x-- {
				totalDistance += visitPlot(x, y, &currentPos, treeMap)
			}
		}
	}
	
	// Return to ground level at the last plot
	totalDistance += currentPos.z
	
	return totalDistance
}

// calculateDronePathWithRest calculates the drone path and determines the rest position
func calculateDronePathWithRest(width, length int, treeMap map[string]repository.Tree, maxDistance int) (int, position) {
	totalDistance := 0
	restPos := position{x: 0, y: 0, z: 0}
	
	// Start at ground level at the southwestern-most plot (1,1)
	currentPos := position{x: 1, y: 1, z: 0}
	
	// The drone travels in a zigzag pattern from south to north
	outerLoop:
	for y := 1; y <= length; y++ {
		// For even-numbered rows, go from west to east
		if y%2 == 1 {
			for x := 1; x <= width; x++ {
				distance := visitPlot(x, y, &currentPos, treeMap)
				totalDistance += distance
				
				// Check if we've reached max distance
				if totalDistance >= maxDistance {
					// If we go over the max distance exactly at the plot, the rest position is the current plot
					restPos = position{x: currentPos.x, y: currentPos.y, z: 0} // z=0 because we land on the ground
					break outerLoop
				}
			}
		} else { // For odd-numbered rows, go from east to west
			for x := width; x >= 1; x-- {
				distance := visitPlot(x, y, &currentPos, treeMap)
				totalDistance += distance
				
				// Check if we've reached max distance
				if totalDistance >= maxDistance {
					// If we go over the max distance exactly at the plot, the rest position is the current plot
					restPos = position{x: currentPos.x, y: currentPos.y, z: 0} // z=0 because we land on the ground
					break outerLoop
				}
			}
		}
	}
	
	// If we've completed the entire path without reaching max distance
	if totalDistance < maxDistance {
		// The drone rests at the final plot
		restPos = position{x: currentPos.x, y: currentPos.y, z: 0}
	}
	
	// For both cases, we need to descend to ground level for the rest, but we've already
	// accounted for this in the totalDistance calculation for the rest position
	
	return totalDistance, restPos
}

// visitPlot calculates the distance to visit a single plot
func visitPlot(x, y int, currentPos *position, treeMap map[string]repository.Tree) int {
	distance := 0
	
	// First move horizontally to the plot
	distance += abs(x - currentPos.x) + abs(y - currentPos.y)
	
	// Update current position horizontally
	currentPos.x = x
	currentPos.y = y
	
	// Determine target height
	targetHeight := 1 // 1m above an empty plot
	key := coordKey(x, y)
	if tree, exists := treeMap[key]; exists {
		targetHeight = tree.Height + 1 // 1m above the tree
	}
	
	// Move vertically to the target height
	distance += abs(targetHeight - currentPos.z)
	
	// Update current position vertically
	currentPos.z = targetHeight
	
	return distance
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
} 