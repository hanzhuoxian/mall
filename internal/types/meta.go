package types

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Objecter 是所有资源对象的通用访问接口，提供对主键、名称及时间字段的读写能力。
type Objecter interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
}

// Lister 是列表资源的通用接口，用于读写分页总数。
type Lister interface {
	GetTotalCount() int64
	SetTotalCount(count int64)
}

// Extend 是键值扩展字段，用于存储任意附加信息。
type Extend map[string]any

// String 将 Extend 序列化为 JSON 字符串。
func (e Extend) String() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// Merge 将 extendShadow（JSON 字符串）中的键合并到当前 Extend，
// 仅补充缺失的键，不覆盖已有值。
func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend
	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

// ObjectMeta 是所有数据库模型的公共元数据基类，嵌入后可获得主键、
// 唯一实例 ID、名称、扩展字段及软删除支持。
type ObjectMeta struct {
	// InstanceID 是系统生成的全局唯一实例标识符。
	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instance_id;type:varchar(32);not null"`
	// Name 是资源名称。
	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(64);not null" validate:"name"`
	// Extend 存储扩展键值对，不落库，通过 ExtendShadow 与数据库同步。
	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`
	// ExtendShadow 是 Extend 的 JSON 序列化字符串，持久化到数据库。
	ExtendShadow string `json:"-" gorm:"column:extend_shadow" validate:"omitempty"`
	// CreatedAt 记录资源创建时间，由数据库自动填充。
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	// UpdatedAt 记录资源最后更新时间，由数据库自动维护。
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	// DeletedAt 记录软删除时间，由系统在优雅删除时填充，只读。
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index:idx_deleted_at"`
	// ID 是自增主键。
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`
}

var _ Objecter = &ObjectMeta{}

// GetObjectMeta 返回 ObjectMeta 自身，实现 Objecter 接口的元数据访问方法。
func (meta *ObjectMeta) GetObjectMeta() Objecter { return meta }

// GetID 返回资源的自增主键。
func (meta *ObjectMeta) GetID() uint64 { return meta.ID }

// SetID 设置资源的自增主键。
func (meta *ObjectMeta) SetID(id uint64) { meta.ID = id }

// GetName 返回资源名称。
func (meta *ObjectMeta) GetName() string { return meta.Name }

// SetName 设置资源名称。
func (meta *ObjectMeta) SetName(name string) { meta.Name = name }

// GetCreatedAt 返回资源创建时间。
func (meta *ObjectMeta) GetCreatedAt() time.Time { return meta.CreatedAt }

// SetCreatedAt 设置资源创建时间。
func (meta *ObjectMeta) SetCreatedAt(createdAt time.Time) { meta.CreatedAt = createdAt }

// GetUpdatedAt 返回资源最后更新时间。
func (meta *ObjectMeta) GetUpdatedAt() time.Time { return meta.UpdatedAt }

// SetUpdatedAt 设置资源最后更新时间。
func (meta *ObjectMeta) SetUpdatedAt(updatedAt time.Time) { meta.UpdatedAt = updatedAt }

// BeforeCreate 在记录写入数据库前将 Extend 序列化到 ExtendShadow。
func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// BeforeUpdate 在记录更新到数据库前将 Extend 序列化到 ExtendShadow。
func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// AfterFind 在从数据库查询后将 ExtendShadow 反序列化回 Extend。
func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

// ListMeta 是列表资源的公共元数据，包含分页总数。
type ListMeta struct {
	// TotalCount 是满足查询条件的记录总数。
	TotalCount int64 `json:"totalCount,omitempty"`
}

var _ Lister = &ListMeta{}

// GetListMeta 返回 ListMeta 自身，实现 Lister 接口的元数据访问方法。
func (l *ListMeta) GetListMeta() Lister { return l }

// SetTotalCount 设置列表的记录总数。
func (l *ListMeta) SetTotalCount(c int64) { l.TotalCount = c }

// GetTotalCount 返回列表的记录总数。
func (l *ListMeta) GetTotalCount() int64 { return l.TotalCount }
