package setup

import (
	"database/sql"

	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/repository"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/server"
	"github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/internal/service"
	pb "github.com/sir-shalahuddin/grpc-learn/bookcategoryservice/proto/category"
	"google.golang.org/grpc"
)

func GRPCServer(grpc *grpc.Server, db *sql.DB) {
	// Initialize repositories, services, and servers
	categoryRepo := repository.NewBookCategoryRepository(db)
	categoryService := service.NewBookCategoryService(categoryRepo)
	categoryServer := server.NewBookCategoryGRPCServer(categoryService)

	// Register AuthService routes
	pb.RegisterBookCategoryServiceServer(grpc, categoryServer)

}
