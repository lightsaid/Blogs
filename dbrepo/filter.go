package dbrepo

import (
	"math"
	"strings"
)

// Filters 过滤参数
type Filters struct {
	Page           int      // 第几页
	PageSize       int      // 每页多少条
	SortFields     []string // 排序的字段，必须存在于安全的字段SortSafeFields里
	SortSafeFields []string // 允许排序的字段, 必须预先设置好, 设置规则：id、-id，对应排序规则：id:ASC、-id:DESC
}

// Metadata 查询列表的基础（元）数据
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// sortField 排序字段设置, 返回格式：xx ASC, yy DESC
// 如果没有排序字段会返回空串，因此使用时注意空串情况
func (f Filters) sortField() string {
	var sortStrs []string
	for _, field := range f.SortFields {
		for _, safeVal := range f.SortSafeFields {
			if field == safeVal {
				if strings.HasPrefix(field, "-") {
					sortStr := strings.Trim(field, "-") + " DESC"
					sortStrs = append(sortStrs, sortStr)
				} else {
					sortStr := field + " DESC"
					sortStrs = append(sortStrs, sortStr)
				}
			}
		}
	}

	return strings.Join(sortStrs, ",")
}

// defaultSort 提供一个默认排序值, updated_at DESC, 表里必须要有 updated_at 字段
func (f Filters) defaultSort() string {
	return " updated_at DESC "
}

func (f Filters) limit() int {
	if f.Page < 1 {
		return 1
	}
	return f.PageSize
}

func (f Filters) offset() int {
	if f.PageSize <= 0 || f.PageSize > 100 {
		return 10
	}
	return (f.Page - 1) * f.PageSize
}

// calculateMetadata 计算基础元数据
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
