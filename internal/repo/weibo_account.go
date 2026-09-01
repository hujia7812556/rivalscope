package repo

import (
	"gorm.io/gorm"

	"rivalscope/internal/dto"
	"rivalscope/internal/model"
)

// WeiboAccountRepo 微博账号数据访问。
type WeiboAccountRepo struct {
	db *gorm.DB
}

// NewWeiboAccountRepo 创建微博账号 repo。
func NewWeiboAccountRepo(db *gorm.DB) *WeiboAccountRepo {
	return &WeiboAccountRepo{db: db}
}

// List 查询全部有效账号(按 id 升序)。
func (r *WeiboAccountRepo) List() ([]model.WeiboAccount, error) {
	var list []model.WeiboAccount
	err := r.db.Where("status = ?", model.WeiboAccountStatusValid).
		Order("id asc").Find(&list).Error
	return list, err
}

// GetByID 按主键查询;未找到返回 404 业务错误。
func (r *WeiboAccountRepo) GetByID(id int) (*model.WeiboAccount, error) {
	var a model.WeiboAccount
	err := r.db.First(&a, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, dto.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetByUid 按 uid 查询;未找到返回 (nil, nil)。
func (r *WeiboAccountRepo) GetByUid(uid string) (*model.WeiboAccount, error) {
	var a model.WeiboAccount
	err := r.db.Where("uid = ?", uid).First(&a).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Create 新增账号。
func (r *WeiboAccountRepo) Create(a *model.WeiboAccount) error {
	return r.db.Create(a).Error
}

// Update 保存账号(全量字段)。
func (r *WeiboAccountRepo) Update(a *model.WeiboAccount) error {
	return r.db.Save(a).Error
}

// Delete 按主键物理删除账号。
func (r *WeiboAccountRepo) Delete(id int) error {
	return r.db.Delete(&model.WeiboAccount{}, id).Error
}
