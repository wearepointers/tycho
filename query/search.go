package query

type Search string

func ParseSearch(raw string) *Search {
	search := Search(raw)
	return &search
}
