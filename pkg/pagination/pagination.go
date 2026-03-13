package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params contains pagination parameters
type Params struct {
	Limit  int
	Offset int
	Page   int
}

// DefaultParams returns default pagination parameters
func DefaultParams() Params {
	return Params{
		Limit:  10,
		Offset: 0,
		Page:   1,
	}
}

// FromContext extracts pagination from gin context
func FromContext(c *gin.Context) Params {
	params := DefaultParams()

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil && val > 0 {
			params.Limit = val
		}
	}

	if page := c.Query("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil && val > 0 {
			params.Page = val
			params.Offset = (params.Page - 1) * params.Limit
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil && val >= 0 {
			params.Offset = val
		}
	}

	return params
}

// Meta contains pagination metadata
type Meta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Page   int   `json:"page"`
	Pages  int   `json:"pages"`
}

// NewMeta creates pagination metadata
func NewMeta(total int64, params Params) *Meta {
	pages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		pages++
	}

	return &Meta{
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
		Page:   params.Page,
		Pages:  pages,
	}
}
