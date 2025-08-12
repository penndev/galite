package orm

import (
	"time"

	"gorm.io/gorm"
)

var defaultDB *gorm.DB

func SetDB(db *gorm.DB) {
	if db == nil {
		panic("db is nil")
	}
	defaultDB = db
}

type ModelInterface interface {
	DB() *gorm.DB
	Model() *gorm.DB
	//
	Bind(bindModel ModelInterface, param ...any) ModelInterface
	BindGorm() *gorm.DB

	// 重载语法糖
	Save()
	Find(dest any, conds ...any) (tx *gorm.DB)
	First(dest any, conds ...any) (tx *gorm.DB)

	// 拓展语法糖

	// List 方法用于查询数据列表。
	// 参数:
	// - total: [*int64] 用于存储查询结果的总数。
	// - data: [any] 用于存储查询结果的数据集。
	// 返回:
	// - error: 如果查询过程中发生错误，则返回错误信息。
	// 如果查询成功，则返回 nil。
	List(total *int64, data any) error
}

type ModelBase struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// gorm.ModelBase
	bindModel  ModelInterface //查询绑定的表实例
	bindScopes func(*gorm.DB) *gorm.DB
	bindParam  BindParam
}

// DB 方法返回 gorm.DB 实例。该实例未关联模型
func (m *ModelBase) DB() *gorm.DB {
	return defaultDB
}

// Model 方法返回绑定的模型的 gorm.DB 实
// 如果没有绑定模型，则会抛出 panic。
func (m *ModelBase) Model() *gorm.DB {
	if m.bindModel == nil {
		panic("bindModel is nil, please use Bind() method to bind model")
	}
	return m.DB().Model(m.bindModel)
}

// Bind 方法绑定模型和可选参数。
//
// 参数:
// -bindModel:
//   - m.bindModel [*Model] 要绑定的模型对象, 用于反射。如果传入了结构体参数则会使用Where条件进行查询。
//
// -param: [可选] 动态的参数
//   - args[1] m.bindScopes [`func(*gorm.DB) *gorm.DB`]: gorm自定义条件注入。
//   - args[2] m.bindParam  [`BindParam`]: 用于分页和排序处理封装。
//
// 返回:
// - [*gorm.DB] 经过条件绑定的原生gorm。
func (m *ModelBase) Bind(bindModel ModelInterface, param ...any) ModelInterface {
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
	return m
}

// 返回经过条件绑定处理的原生Gorm
func (m *ModelBase) BindGorm() *gorm.DB {
	query := m.DB().Where(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	if m.bindParam != nil {
		query = query.Order(m.bindParam.Order()).Offset(m.bindParam.Offset()).Limit(m.bindParam.Limit())
	}
	return query
}

// 重载Gorm语法糖
func (m *ModelBase) Save() {
	m.BindGorm().Save(m.bindModel)
}

func (m *ModelBase) Find(dest any, conds ...any) (tx *gorm.DB) {
	return m.BindGorm().Find(dest, conds...)
}

func (m *ModelBase) First(dest any, conds ...any) (tx *gorm.DB) {
	return m.BindGorm().First(dest, conds...)
}

func (m *ModelBase) List(total *int64, data any) error {
	query := m.Model().Where(m.bindModel)
	if m.bindScopes != nil {
		query = query.Scopes(m.bindScopes)
	}
	if err := query.Count(total).Error; err != nil {
		return err
	}
	if m.bindParam != nil { //数据库穿透攻击，上游处理
		// count 可以被db缓存，实际查询不会，如果limit offset 超过 count 则进行何种操作?
		// log.Println(m.bindParam.Order(), m.bindParam.Offset(), m.bindParam.Limit())
		query = query.Order(m.bindParam.Order()).Offset(m.bindParam.Offset()).Limit(m.bindParam.Limit())
	}
	if err := query.Find(data).Error; err != nil {
		return err
	}
	return nil
}
