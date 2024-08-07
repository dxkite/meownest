package {{ .ModuleName }}

import (
	"context"

	"{{ .PackageName }}/pkg/database"
	"gorm.io/gorm"
)

type {{ .Name }}Repository interface {
	Create(ctx context.Context, {{ .LowerCamelName }} *{{ .Name }}) (*{{ .Name }}, error)
	Get(ctx context.Context, id uint64) (*{{ .Name }}, error)
	Update(ctx context.Context, id uint64, ent *{{ .Name }}) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error)
	BatchGet(ctx context.Context, ids []uint64) ([]*{{ .Name }}, error)
}

func New{{ .Name }}Repository() {{ .Name }}Repository {
	return new({{ .LowerCamelName }}Repository)
}

type {{ .LowerCamelName }}Repository struct {
}

func (r *{{ .LowerCamelName }}Repository) Get(ctx context.Context, id uint64) (*{{ .Name }}, error) {
	var item {{ .Name }}
	if err := r.dataSource(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *{{ .LowerCamelName }}Repository) BatchGet(ctx context.Context, ids []uint64) ([]*{{ .Name }}, error) {
	var items []*{{ .Name }}
	if err := r.dataSource(ctx).Where("id in ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type List{{ .Name }}Param struct {
	// TODO external condition

	// pagination
	Page         int
	PerPage      int
	IncludeTotal bool
}

type List{{ .Name }}Result struct {
	Data  []*{{ .Name }}
	Total int64
}


func (r *{{ .LowerCamelName }}Repository) List(ctx context.Context, param *List{{ .Name }}Param) (*List{{ .Name }}Result, error) {
	var items []*{{ .Name }}
	db := r.dataSource(ctx)

	// condition
	condition := func(db *gorm.DB) *gorm.DB {
		return db
	}

	// pagination
	query := db.Scopes(condition)
	if param.Page > 0 && param.PerPage > 0 {
		query.Offset((param.Page - 1) * param.PerPage).Limit(param.PerPage)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	rst := &List{{ .Name }}Result{}
	rst.Data = items

	if param.IncludeTotal {
		if err := db.Model({{ .Name }}{}).Scopes(condition).Count(&rst.Total).Error; err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func (r *{{ .LowerCamelName }}Repository) Create(ctx context.Context, {{ .LowerCamelName }} *{{ .Name }}) (*{{ .Name }}, error) {
	if err := r.dataSource(ctx).Create(&{{ .LowerCamelName }}).Error; err != nil {
		return nil, err
	}
	return {{ .LowerCamelName }}, nil
}

func (r *{{ .LowerCamelName }}Repository) Update(ctx context.Context, id uint64, ent *{{ .Name }}) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Updates(&ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .LowerCamelName }}Repository) Delete(ctx context.Context, id uint64) error {
	if err := r.dataSource(ctx).Where("id = ?", id).Delete({{ .Name }}{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .LowerCamelName }}Repository) dataSource(ctx context.Context) *gorm.DB {
	return database.Get(ctx).Engine().(*gorm.DB)
}
