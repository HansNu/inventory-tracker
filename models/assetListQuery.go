package models

// AssetListQuery is the normalised form of AssetListParams: strings parsed,
// defaults applied, sort direction resolved. The service produces it, the
// repository consumes it — so the repository never has to guess what "page 0"
// or an empty sortOrder means.
type AssetListQuery struct {
	Search        string
	Status        string
	AssetCategory string
	SortField     string
	SortOrder     string // "ASC" or "DESC"
	Limit         int
	Offset        int
}
