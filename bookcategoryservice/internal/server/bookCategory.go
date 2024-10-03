package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/models"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/category"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bookCategoryService interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.BookCategory, error)
	GetAllCategories(ctx context.Context) ([]models.BookCategory, error)
}

type bookCategoryGRPCServer struct {
	pb.UnimplementedBookCategoryServiceServer // Embed to have forward compatible implementations.
	service                                   bookCategoryService
}

// NewBookCategoryGRPCServer creates a new instance of BookCategoryGRPCServer.
func NewBookCategoryGRPCServer(service bookCategoryService) *bookCategoryGRPCServer {
	return &bookCategoryGRPCServer{service: service}
}

// GetCategories retrieves all categories from the database.
func (s *bookCategoryGRPCServer) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.CategoryListResponse, error) {
	categories, err := s.service.GetAllCategories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve category: %v", err)
	}

	var categoryList []*pb.CategoryResponse
	for _, category := range categories {
		categoryList = append(categoryList, &pb.CategoryResponse{
			Id:   category.ID.String(),
			Name: category.Name,
		})
	}

	return &pb.CategoryListResponse{Categories: categoryList}, nil
}

// GetCategoryByID retrieves a category by its ID from the database.
func (s *bookCategoryGRPCServer) GetCategoryByID(ctx context.Context, req *pb.GetCategoryByIDRequest) (*pb.CategoryResponse, error) {
	// Parse the UUID from the request
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID format: %v", err)
	}

	// Get the category from the service layer
	category, err := s.service.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve category: %v", err)
	}

	// If no category is found, return a NotFound error
	if category == nil {
		return nil, status.Error(codes.NotFound, "category not found")
	}

	// Return the category in the response
	return &pb.CategoryResponse{
		Id:   category.ID.String(),
		Name: category.Name,
	}, nil
}
