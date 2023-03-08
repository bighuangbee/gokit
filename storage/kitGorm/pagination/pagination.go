package pagination

import "errors"

const (
	defaultLimit     = 10
	maxLimit         = 1000
	defaultOrderType = "desc"
	defaultOrderBy   = "id"
)

type IPageArg interface {
	GetOffset() int32
	GetLimit() int32
	GetOrderBy() string
	GetOrderType() string

	GetPageNo() int32
	GetPageSize() int32
}

type IPagination interface {
	GetOffset() int
	GetLimit() int
	GetOrderBy() string
	GetOrderType() string
	IsValidOrderField() bool

	GetPageNo() int
	GetPageSize() int
}

type Pagination struct {
	Offset    int
	Limit     int
	OrderBy   string
	OrderType string
	WhileList map[string]int8

	PageNo    int
	PageSzie  int
}

func (p *Pagination) GetOffset() int {
	return p.Offset
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetPageNo() int {
	return p.PageNo
}

func (p *Pagination) GetPageSize() int {
	return p.PageSzie
}


func (p *Pagination) GetOrderBy() string {
	return p.OrderBy
}

func (p *Pagination) GetOrderType() string {
	return p.OrderType
}

func (p *Pagination) IsValidOrderField() bool {
	_, ok := p.WhileList[p.OrderBy]
	return ok
}

func (p *Pagination) SetWhiteList(fields []string) {
	for _, v := range fields {
		p.WhileList[v] = 1
	}
}

func New(page IPageArg, opts ...option) (IPagination, error) {
	p := &Pagination{
		Offset:    int(page.GetOffset()),
		Limit:     int(page.GetLimit()),
		OrderBy:   page.GetOrderBy(),
		OrderType: page.GetOrderType(),
		WhileList: map[string]int8{
			"id":        1,
			"create_at": 1,
		},
	}
	return NewPage(p, opts...)
}

func NewPage(p *Pagination, opts ...option) (IPagination, error) {
	if p.Limit < 1 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	if p.OrderBy == "" {
		p.OrderBy = defaultOrderBy
	}
	if p.OrderType == "" {
		p.OrderType = defaultOrderType
	}
	for _, fn := range opts {
		fn(p)
	}

	if !p.IsValidOrderField() {
		return nil, errors.New("无效的排序字段")
	}

	return p, nil
}

type option func(p *Pagination) *Pagination

func SetWhiteList(fields []string) func(p *Pagination) *Pagination {
	return func(p *Pagination) *Pagination {
		p.SetWhiteList(fields)
		return p
	}
}
