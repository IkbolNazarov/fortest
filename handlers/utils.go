package handlers

import (
	"fortest/models"
	"net/http"
	"strconv"
)

func GeneratePaginationFromRequest(c *http.Request) models.Pagination {
	limit := 16
	page := 1
	query := c.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		}
	}
	return models.Pagination{
		Limit: limit,
		Page:  page,
	}
}

func TotalPageTasks(limit int64, length int64) (int64, error) {
	if limit == 0 {
		limit = 10
	}

	totalPage := length / limit
	if length%limit != 0 {
		totalPage++
	}
	return totalPage, nil
}
