package constants

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

func (s SortOrder) SQL() string {
	if s == SortOrderDesc {
		return "DESC"
	}
	return "ASC"
}
