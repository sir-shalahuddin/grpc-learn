package repository

import (
	"context"
	"fmt"

	pb "github.com/sir-shalahuddin/grpc-learn/bookservice/proto/categoryservice" // Adjust the import path as needed
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type categoryRepository struct {
	client pb.BookCategoryServiceClient
}

// NewCategoryRepository creates a new instance of categoryRepository.
func NewCategoryRepository(client pb.BookCategoryServiceClient) *categoryRepository {
	return &categoryRepository{client: client}
}

// GetCategories retrieves all categories using the gRPC client.
func (r *categoryRepository) GetCategories(ctx context.Context) ([]*pb.CategoryResponse, error) {
	req := &pb.GetCategoriesRequest{}

	resp, err := r.client.GetCategories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return resp.Categories, nil
}

// GetCategoryByID retrieves a category by ID using the gRPC client.
func (r *categoryRepository) GetCategoryByID(ctx context.Context, id string) (*pb.CategoryResponse, error) {
	req := &pb.GetCategoryByIDRequest{
		Id: id,
	}

	// Call the gRPC client method to get the category
	resp, err := r.client.GetCategoryByID(ctx, req)
	if err != nil {
		// Check if the error is a NotFound error
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	// Return the response from the gRPC server
	return resp, nil
}
