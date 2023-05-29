package paginate

import (
	"errors"
	"math"
	"strconv"

	"raptor/models"
	"github.com/labstack/echo/v4"
)

// PaginateInfo ...
type PaginateInfo struct {
	Limit,
	Page,
	Offset,
	Pages,
	Total,
	Prev,
	Next int
}

// Paginate ..
func Paginate(limit, page, total int) *PaginateInfo {
	info := &PaginateInfo{
		Limit:  limit,
		Page:   page,
		Total:  total,
		Offset: page * limit,
	}
	tmpLimit := float64(limit)
	tmpTotal := float64(total)

	info.Pages = int(math.Ceil(tmpTotal / tmpLimit))

	if (page + 1) < info.Pages {
		info.Next = page + 1
	}

	if page > 0 {
		info.Prev = page - 1
	}

	return info
}

// ParsePaginationParams info from the context of http request
func ParsePaginationParams(ctx echo.Context) (int, int, error) {
	var err error
	var page = 0
	var limit = 10

	// Case of url without queryParams
	if ctx.QueryParam("limit") == "" && ctx.QueryParam("page") == "" {
		return limit, page, nil
	}

	limit, err = strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil {
		return 0, 0, models.InvalidSyntaxErr{Err: "Invalid Syntax"}
	}
	page, err = strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return 0, 0, models.InvalidSyntaxErr{Err: "Invalid Syntax"}
	}
	if limit < 0 {
		err = errors.New("limit must be positive")
	}
	if page < 0 {
		err = errors.New("page must be positive or zero")
	}
	if limit == 0 {
		limit = 10
	}

	// Maximum products per request
	if limit > 100 {
		limit = 100
	}
	return limit, page, err
}
