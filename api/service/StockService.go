package service

import (
	"github.com/sgomeza13/stock-recommender/api/models"
	"github.com/sgomeza13/stock-recommender/api/repository"
)

type StockService struct {
	Repository *repository.StockRepository
}

func NewStockService(stockRepo *repository.StockRepository) *StockService {
	return &StockService{
		Repository: repository.NewStockRepository(),
	}
}

func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	return s.Repository.GetAllStocks()
}

// Define a pagination response struct at the service level
type PaginatedStocksResponse struct {
	Stocks     []models.Stock `json:"stocks"`
	TotalCount int            `json:"totalCount"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// Updated service method with page-based pagination
func (s *StockService) GetStocksPaginated(page, pageSize int) (PaginatedStocksResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}

	// Call the repository with the updated pagination method
	paginatedStocks, err := s.Repository.GetStocksPaginated(page, pageSize)
	if err != nil {
		return PaginatedStocksResponse{}, err
	}

	// Map repository response to service response
	return PaginatedStocksResponse{
		Stocks:     paginatedStocks.Stocks,
		TotalCount: paginatedStocks.TotalCount,
		Page:       paginatedStocks.Page,
		PageSize:   paginatedStocks.PageSize,
		TotalPages: paginatedStocks.TotalPages,
	}, nil
}
func (s *StockService) GetStockByID(id int) (*models.Stock, error) {
	return s.Repository.GetStockByID(id)
}

func (s *StockService) CreateStock(stock *models.Stock) error {
	return s.Repository.CreateStock(stock)
}

func (s *StockService) CreateStocks(stocks []*models.Stock) error {
	return s.Repository.CreateStocks(stocks)
}

func (s *StockService) DeleteStockByID(id int) error {
	return s.Repository.DeleteStockByID(id)
}

func (s *StockService) UpdateStockByID(id int, stock *models.Stock) error {
	return s.Repository.UpdateStockByID(id, stock)
}
