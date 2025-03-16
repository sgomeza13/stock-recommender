package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sgomeza13/stock-recommender/api/models"
	"github.com/sgomeza13/stock-recommender/api/service"
	"github.com/sgomeza13/stock-recommender/utils"
)

type StockController struct {
	StockService *service.StockService
}

func NewStockController(stockService *service.StockService) *StockController {
	return &StockController{StockService: stockService}
}

// ✅ Handle all stocks request
func (sc *StockController) GetAllStocks(c *gin.Context) {
	stocks, err := sc.StockService.GetAllStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// ✅ Handle paginated stock request
func (sc *StockController) GetStocksPaginated(c *gin.Context) {
	// Get query params
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	// Call service
	stocks, err := sc.StockService.GetStocksPaginated(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, stocks)
}

func (sc *StockController) GetStockByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	stock, err := sc.StockService.GetStockByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if stock == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// CreateStock handles a single stock creation request
func (c *StockController) CreateStock(ctx *gin.Context) {
	var input map[string]string

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	stock, err := c.parseStockFromMap(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.StockService.CreateStock(stock); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Stock created successfully"})
}

// CreateStocks handles multiple stocks creation request
func (c *StockController) CreateStocks(ctx *gin.Context) {
	var rawStocks []map[string]interface{}

	if err := ctx.ShouldBindJSON(&rawStocks); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	stocks := make([]*models.Stock, 0, len(rawStocks))

	for i, rawStock := range rawStocks {
		// Convert the map[string]interface{} to map[string]string for consistent handling
		stringMap := make(map[string]string)
		for k, v := range rawStock {
			// Handle different value types appropriately
			switch val := v.(type) {
			case string:
				stringMap[k] = val
			case float64:
				// Convert numeric values to string
				stringMap[k] = fmt.Sprintf("%g", val)
			case int:
				stringMap[k] = fmt.Sprintf("%d", val)
			case nil:
				// Handle nil values as empty strings
				stringMap[k] = ""
			default:
				// For other types, try string conversion or fail
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Field '%s' in item %d has unsupported type: %T", k, i, v),
				})
				return
			}
		}

		stock, err := c.parseStockFromMap(stringMap)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Error in item %d: %s", i, err.Error()),
				"item":  rawStock, // Include the problematic item for debugging
			})
			return
		}

		stocks = append(stocks, stock)
	}

	if err := c.StockService.CreateStocks(stocks); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stocks"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Stocks created successfully"})
}

// parseStockFromMap transforms a map into a Stock model, handling validation and conversion
func (c *StockController) parseStockFromMap(input map[string]string) (*models.Stock, error) {
	// Check if required fields exist
	requiredFields := []string{"ticker", "target_from", "target_to", "company", "action", "brokerage", "rating_from", "rating_to", "time"}
	for _, field := range requiredFields {
		if value, exists := input[field]; !exists || value == "" {
			return nil, fmt.Errorf("missing required field: %s", field)
		}
	}

	// Parse target_from with the improved CleanDecimal function
	targetFrom, err := utils.CleanDecimal(input["target_from"])
	if err != nil {
		return nil, fmt.Errorf("invalid target_from value '%s': %w", input["target_from"], err)
	}

	// Parse target_to with the improved CleanDecimal function
	targetTo, err := utils.CleanDecimal(input["target_to"])
	if err != nil {
		return nil, fmt.Errorf("invalid target_to value '%s': %w", input["target_to"], err)
	}

	// Parse the time string into a time.Time object
	timeStr := input["time"]
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// If the standard RFC3339 format fails, try a more flexible approach
		parsedTime, err = parseTimeFlexibly(timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time format '%s': %w", timeStr, err)
		}
	}

	return &models.Stock{
		Ticker:     input["ticker"],
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Company:    input["company"],
		Action:     input["action"],
		Brokerage:  input["brokerage"],
		RatingFrom: input["rating_from"],
		RatingTo:   input["rating_to"],
		Time:       parsedTime,
	}, nil
}

// parseTimeFlexibly tries multiple common time formats to parse a time string
func parseTimeFlexibly(timeStr string) (time.Time, error) {
	// Try various common formats
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05", // ISO without timezone
		"2006-01-02 15:04:05", // Common SQL datetime format
		"2006-01-02",          // Simple date only
		"01/02/2006",          // US date format
		"02/01/2006",          // European date format
		"2006/01/02",          // Year first date format
	}

	var parseErr error
	for _, format := range formats {
		parsedTime, err := time.Parse(format, timeStr)
		if err == nil {
			return parsedTime, nil
		}
		parseErr = err
	}

	// If all formats fail, return the last error
	return time.Time{}, fmt.Errorf("could not parse time string with any known format: %w", parseErr)
}
func (sc *StockController) DeleteStockByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	if err := sc.StockService.DeleteStockByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock deleted successfully"})
}

func (sc *StockController) UpdateStockByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.StockService.UpdateStockByID(id, &stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
