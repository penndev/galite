package sugar

import (
	"log"
	"time"

	"github.com/penndev/wga/config"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// gorm.Model
	bindModel any //查询绑定的表
	bindWhere func(*gorm.DB) *gorm.DB
	bindParam BindParam
}

func (m *Model) Bind(bindModel any, param ...any) {
	m.bindModel = bindModel
	if len(param) >= 1 {
		if item, ok := param[0].(func(*gorm.DB) *gorm.DB); ok {
			m.bindWhere = item
		}
	}
	if len(param) >= 2 {
		if item, ok := param[1].(BindParam); ok {
			m.bindParam = item //设置分页排序
		}
	}
}

func (m *Model) List(t *int64, l any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindWhere != nil {
		query = query.Scopes(m.bindWhere)
	}
	if err := query.Count(t).Error; err != nil {
		return err
	}
	if m.bindParam != nil { //数据库穿透攻击，上游处理
		//  count 可以被db缓存，实际查询不会，如果limit offset 超过 count 则进行何种操作?
		log.Println(m.bindParam.Order(), m.bindParam.Offset(), m.bindParam.Limit())
		query = query.Order(m.bindParam.Order()).Offset(m.bindParam.Offset()).Limit(m.bindParam.Limit())
	}
	if err := query.Find(l).Error; err != nil {
		return err
	}
	return nil
}

func (m *Model) Create(data any) error {
	return config.DB.Model(m.bindModel).Create(data).Error
}

func (m *Model) Update(data any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindWhere != nil {
		query = query.Scopes(m.bindWhere)
	}
	return query.Updates(data).Error
}

func (m *Model) Delete(data any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindWhere != nil {
		query = query.Scopes(m.bindWhere)
	}
	return query.Delete(data).Error
}
