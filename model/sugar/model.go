package sugar

import (
	"time"

	"github.com/penndev/wga/config"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// gorm.Model
	bindModel  any //查询绑定的表
	bindScopes func(*gorm.DB) *gorm.DB
	bindParam  BindParam
}

// Bind 方法绑定模型和可选参数。
//
// 参数:
//   - m.bindModel [*Model]: 要绑定的模型对象, 用于反射。
//   - m.bindScopes [`func(*gorm.DB) *gorm.DB`]: 函数参数gorm条件注入。
//   - m.bindParam  [`BindParam`]: 用于分页和排序处理封装。
func (m *Model) Bind(bindModel any, param ...any) *gorm.DB {
	m.bindModel = bindModel
	if len(param) >= 1 {
		if item, ok := param[0].(func(*gorm.DB) *gorm.DB); ok {
			m.bindScopes = item
		}
	}
	if len(param) >= 2 {
		if item, ok := param[1].(BindParam); ok {
			m.bindParam = item //设置分页排序
		}
	}
	return m.Gorm()
}

// 返回经过条件绑定的原生gorm
func (m *Model) Gorm() *gorm.DB {
	query := config.DB.Model(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	if m.bindParam != nil {
		query = query.Order(m.bindParam.Order()).Offset(m.bindParam.Offset()).Limit(m.bindParam.Limit())
	}
	return query
}

// total 数据总量
// data 结果集
func (m *Model) List(total *int64, data any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	if err := query.Count(total).Error; err != nil {
		return err
	}
	if m.bindParam != nil { //数据库穿透攻击，上游处理
		//  count 可以被db缓存，实际查询不会，如果limit offset 超过 count 则进行何种操作?
		// log.Println(m.bindParam.Order(), m.bindParam.Offset(), m.bindParam.Limit())
		query = query.Order(m.bindParam.Order()).Offset(m.bindParam.Offset()).Limit(m.bindParam.Limit())
	}
	if err := query.Find(data).Error; err != nil {
		return err
	}
	return nil
}

func (m *Model) Create(data any) error {
	return config.DB.Model(m.bindModel).Create(data).Error
}

func (m *Model) Save(data any) error {
	return config.DB.Model(m.bindModel).Save(data).Error
}

func (m *Model) Update(data any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	return query.Updates(data).Error
}

func (m *Model) Delete(data any) error {
	query := config.DB.Model(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	return query.Delete(data).Error
}
