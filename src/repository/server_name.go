package repository

import (
	"context"

	"dxkite.cn/meownest/pkg/datasource"
	"dxkite.cn/meownest/src/entity"
	"gorm.io/gorm"
)

type ServerName interface {
	Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error)
	Get(ctx context.Context, id uint64) (*entity.ServerName, error)
	List(ctx context.Context, param *ListServerNameParam) ([]*entity.ServerName, error)
	Update(ctx context.Context, id uint64, ent *entity.ServerName) error
	Delete(ctx context.Context, id uint64) error
}

func NewServerName() ServerName {
	return &serverName{}
}

type serverName struct {
}

func (r *serverName) Create(ctx context.Context, serverName *entity.ServerName) (*entity.ServerName, error) {
	if err := r.dataSource(ctx).Create(&serverName).Error; err != nil {
		return nil, err
	}
	return serverName, nil
}

func (r *serverName) Get(ctx context.Context, id uint64) (*entity.ServerName, error) {
	var cert entity.ServerName
	if err := r.dataSource(ctx).Where("id = ?", id).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *serverName) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete(entity.ServerName{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverName) Update(ctx context.Context, id uint64, ent *entity.ServerName) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

type ListServerNameParam struct {
	Name          string
	Limit         int
	StartingAfter uint64
	EndingBefore  uint64
}

func (r *serverName) List(ctx context.Context, param *ListServerNameParam) ([]*entity.ServerName, error) {
	var items []*entity.ServerName
	db := r.dataSource(ctx).Model(entity.ServerName{})

	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}

	if param.StartingAfter != 0 {
		db = db.Where("id > ?", param.StartingAfter)
	}

	if param.EndingBefore != 0 {
		db = db.Where("id < ?", param.EndingBefore)
	}

	if param.Limit != 0 {
		db = db.Limit(param.Limit)
	}

	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *serverName) dataSource(ctx context.Context) *gorm.DB {
	return datasource.Get(ctx).DB()
}
