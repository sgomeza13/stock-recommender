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

func (s *StockService) GetStocksPaginated(limit, offset int) ([]models.Stock, error) {
	return s.Repository.GetStocksPaginated(limit, offset)
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
