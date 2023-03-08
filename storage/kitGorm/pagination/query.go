package pagination

import (
	"fmt"

	"gorm.io/gorm"
)

//分页查询
//如果 page==nil 则默认 limit 10 offset 0 orderBy id desc
func PageQuery(db *gorm.DB, page IPagination) *gorm.DB {
	var (
		pOrderType = defaultOrderType
		pOrderBy   = defaultOrderBy
		pLimit     = defaultLimit
		pOffset    = 0
		pNo		   = 0
		pSize	   = 0
	)
	if page != nil {
		pOrderType = page.GetOrderType()
		pOrderBy = page.GetOrderBy()
		pLimit = page.GetLimit()
		pOffset = page.GetOffset()

		pNo = page.GetPageNo()
		pSize = page.GetPageSize()
		if pNo > 0 && pSize > 0{
			pLimit = pSize
			if pSize > 0 {
				pOffset = (pNo - 1) * (pSize)
			}
		}
	}
	db.Order(fmt.Sprintf("%v %v", pOrderBy, pOrderType)).
		Limit(pLimit).
		Offset(pOffset)
	return db
}
