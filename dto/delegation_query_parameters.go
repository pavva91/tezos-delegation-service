package dto

import "time"

type ListDelegationsQueryParameters struct {
	Year time.Time `form:"year" time_format:"2006"`
}
