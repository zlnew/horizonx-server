package domain

import "strings"

type ListOptions struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Search     string `json:"search"`
	IsPaginate bool   `json:"is_paginate"`
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

type ListResult[T any] struct {
	Data []T   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

func CalculateMeta(total int64, page, limit int) *Meta {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	lastPage := int(total) / limit
	if int(total)%limit != 0 {
		lastPage++
	}
	if lastPage == 0 {
		lastPage = 1
	}

	return &Meta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		LastPage:    lastPage,
	}
}

func ContainsAny(s string, xs []string) bool {
	s = strings.ToLower(s)

	for _, x := range xs {
		x = strings.ToLower(x)

		if strings.HasSuffix(x, "*") {
			prefix := strings.TrimSuffix(x, "*")
			if strings.HasPrefix(s, prefix) {
				return true
			}
			continue
		}

		if strings.Contains(s, x) {
			return true
		}
	}

	return false
}
