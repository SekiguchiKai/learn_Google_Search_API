package search_store

import "google.golang.org/appengine/search"

type SearchOptionsWrapper struct {
	SearchOptions search.SearchOptions
}

func NewQuery(field, signal, value string) string{
	return field + signal + value
}
