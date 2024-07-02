package sugar

import (
	"github.com/penndev/wga/config"
	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	bindModel any //绑定的模型实例，
}

func (m *Model) Bind(bindModel any, bindParam BindParam, bindWhere any) {
	m.bindModel = bindModel
}

func (m Model) List(t *int64, l any) error {
	// var total int64
	// var list []interface{}
	query := config.DB.Model(m.bindModel).Where(m.bindModel)
	query.Count(t)
	query.Find(l)
	return nil
}
