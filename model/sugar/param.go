package sugar

type BindParam interface {
	Offset() int
	Limit() int
	Order() any
}

// 请求列表基类,分页查询
type BindListParam struct {
	FormPage  int    `form:"page" binding:"required,min=1"`
	FormLimit int    `form:"limit" binding:"required,min=1,max=100"`
	FormOrder string `form:"order"`
}

// 给列表分页用
func (p *BindListParam) Offset() int {
	return (p.FormPage - 1) * p.FormLimit
}

func (p *BindListParam) Limit() int {
	return p.FormLimit
}

func (p *BindListParam) Order() any {
	switch {
	case p.FormOrder == "id+":
		return "id ASC"
	}
	return "id DESC"
}
