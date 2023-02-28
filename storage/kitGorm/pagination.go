package kitGorm

type PageQuery struct {
	Page  	int `form:"page" json:"page"`
	Length 	int `form:"length" json:"length"`
}

var pageMin = int(1)
var LengthDefault = int(10)
var LengthMax = int(1000)

func (page *PageQuery)Offset()int{
	if page.Page <= 0 {
		page.Page = pageMin
	}
	return (page.Page-1)*page.Limit()
}

func (page *PageQuery)Limit()int{
	if page.Length <= 0 {
		return LengthDefault
	}
	if page.Length > LengthMax {

		return LengthMax
	}
	return page.Length
}

func (page *PageQuery)Pagination()int{
	if page.Page <= 0 {
		page.Page = pageMin
	}
	return page.Page
}

