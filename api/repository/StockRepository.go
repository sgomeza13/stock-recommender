package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/sgomeza13/stock-recommender/api/models"
	"github.com/sgomeza13/stock-recommender/config"
)

type StockRepository struct {
	DB *pgx.Conn
}

func NewStockRepository() *StockRepository {
	return &StockRepository{
		DB: config.GetDB(),
	}
}

// GetAllStocks retrieves all stocks from the database
func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time FROM stock")
	if err != nil {
		log.Println("Error fetching stocks:", err)
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.TargetFrom, &stock.TargetTo,
			&stock.Company, &stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

type PaginatedStocks struct {
	Stocks     []models.Stock
	TotalCount int
	Page       int
	PageSize   int
	TotalPages int
}

func (r *StockRepository) GetStocksPaginated(page, pageSize int) (PaginatedStocks, error) {
	// Calculate offset from page number
	offset := (page - 1) * pageSize

	// First get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM stock`
	err := r.DB.QueryRow(context.Background(), countQuery).Scan(&totalCount)
	if err != nil {
		return PaginatedStocks{}, err
	}

	// Then get paginated data
	query := `SELECT id, ticker, target_from, target_to, company, action, brokerage,
              rating_from, rating_to, time
              FROM stock
              ORDER BY id
              LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		return PaginatedStocks{}, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.ID, &stock.Ticker, &stock.TargetFrom, &stock.TargetTo,
			&stock.Company, &stock.Action, &stock.Brokerage, &stock.RatingFrom,
			&stock.RatingTo, &stock.Time)
		if err != nil {
			return PaginatedStocks{}, err
		}
		stocks = append(stocks, stock)
	}

	// Calculate total pages
	totalPages := totalCount / pageSize
	if totalCount%pageSize > 0 {
		totalPages++
	}

	return PaginatedStocks{
		Stocks:     stocks,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetStockByID retrieves a stock by its ID
func (r *StockRepository) GetStockByID(id int) (*models.Stock, error) {
	var stock models.Stock
	err := r.DB.QueryRow(context.Background(), "SELECT id, ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time FROM stock WHERE id = $1", id).Scan(
		&stock.ID, &stock.Ticker, &stock.TargetFrom, &stock.TargetTo,
		&stock.Company, &stock.Action, &stock.Brokerage,
		&stock.RatingFrom, &stock.RatingTo, &stock.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// CreateStock creates a new stock in the database
func (r *StockRepository) CreateStock(stock *models.Stock) error {
	_, err := r.DB.Exec(context.Background(), "INSERT INTO stock (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		stock.Ticker, stock.TargetFrom, stock.TargetTo,
		stock.Company, stock.Action, stock.Brokerage,
		stock.RatingFrom, stock.RatingTo, stock.Time,
	)
	return err
}

// CreateStocks creates stocks in bulk in the database
func (r *StockRepository) CreateStocks(stocks []*models.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	query := "INSERT INTO stock (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time) VALUES "
	args := []interface{}{}
	argIndex := 1

	for _, stock := range stocks {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			argIndex, argIndex+1, argIndex+2, argIndex+3, argIndex+4,
			argIndex+5, argIndex+6, argIndex+7, argIndex+8)
		args = append(args, stock.Ticker, stock.TargetFrom, stock.TargetTo,
			stock.Company, stock.Action, stock.Brokerage,
			stock.RatingFrom, stock.RatingTo, stock.Time)
		argIndex += 9
	}

	// Remove last comma
	query = query[:len(query)-1]

	_, err := r.DB.Exec(context.Background(), query, args...)
	return err
}

// DeleteStockByID deletes a stock by its ID
func (r *StockRepository) DeleteStockByID(id int) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM stock WHERE id = $1", id)
	return err
}

// UpdateStockByID updates a stock by its ID
func (r *StockRepository) UpdateStockByID(id int, stock *models.Stock) error {
	_, err := r.DB.Exec(context.Background(), "UPDATE stock SET ticker=$1, target_from=$2, target_to=$3, company=$4, action=$5, brokerage=$6, rating_from=$7, rating_to=$8, time=$9 WHERE id=$10",
		stock.Ticker, stock.TargetFrom, stock.TargetTo,
		stock.Company, stock.Action, stock.Brokerage,
		stock.RatingFrom, stock.RatingTo, stock.Time, id,
	)
	return err
}
