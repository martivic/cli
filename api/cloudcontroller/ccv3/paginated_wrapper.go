package ccv3

// PaginatedWrapper represents the standard pagination format of a request that
// returns back more than one object.
type PaginatedWrapper struct {
	Pagination struct {
		Next Link `json:"next"`
	} `json:"pagination"`
	Resources interface{} `json:"resources"`
}

type Link struct {
	URL      string `json:"href"`
	Method   string `json:"method"`
	Metadata struct {
		Version string `json:"version"`
	} `json:"meta"`
}
