package sugar

type BindParam interface {
	Init() //
}

// 后台请求列表基类,分页查询
type BindListParam struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1,max=100"`
}

func (p *BindListParam) Init() {

}
