package models

import "time"

const (
	TimeFormat              = time.RFC3339
	SessionExpiredAfter     = time.Minute * 60
	MaxCategoryLimitForPost = 5

	// SqlSortAsc = "DESC"
	// SqlSortDesc = "ASC"
	SqlLimitInfinity = -1
)
