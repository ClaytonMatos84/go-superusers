package pkg

import (
	"net/http"
	"strconv"
)

type PaginationData struct {
	Page         int `json:"page"`
	StartItems   int `json:"start_items"`
	EndItems     int `json:"end_items"`
	TotalPages   int `json:"total_pages"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
}

func Pagination(queryPage string, queryItemsPerPage string, w http.ResponseWriter, totalItems int) (PaginationData, bool) {
	page := 1
	itemsPerPage := 100
	var err error

	if queryPage != "" {
		page, err = strconv.Atoi(queryPage)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page on request", http.StatusBadRequest)
			return PaginationData{}, true
		}
	}

	if queryItemsPerPage != "" {
		itemsPerPage, err = strconv.Atoi(queryItemsPerPage)
		if err != nil || itemsPerPage < 1 {
			http.Error(w, "Invalid items on request", http.StatusBadRequest)
			return PaginationData{}, true
		}
	}

	start := (page - 1) * itemsPerPage
	if start >= totalItems {
		http.Error(w, "Page number out of range", http.StatusBadRequest)
		return PaginationData{}, true
	}

	end := min(start+itemsPerPage, totalItems)
	paginationData := PaginationData{
		Page:         page,
		StartItems:   start,
		EndItems:     end,
		TotalPages:   (totalItems + itemsPerPage - 1) / itemsPerPage,
		ItemsPerPage: itemsPerPage,
		TotalItems:   totalItems,
	}

	return paginationData, false
}
