package service

import (
	"context"

	"dxkite.cn/meownest/pkg/identity"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/dto"
	"dxkite.cn/meownest/src/entity"
	"dxkite.cn/meownest/src/repository"
)

type Certificate interface {
	Create(ctx context.Context, create *CreateCertificateParam) (*dto.Certificate, error)
	Update(ctx context.Context, param *UpdateCertificateParam) (*dto.Certificate, error)
	Get(ctx context.Context, param *GetCertificateParam) (*dto.Certificate, error)
	Delete(ctx context.Context, param *DeleteCertificateParam) error
	List(ctx context.Context, param *ListCertificateParam) (*ListCertificateResult, error)
}

type CreateCertificateParam struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Key         string `json:"key" form:"key" binding:"required"`
	Certificate string `json:"certificate" form:"key" binding:"required"`
}

func NewCertificate(r repository.Certificate) Certificate {
	return &certificate{r: r}
}

type certificate struct {
	r repository.Certificate
}

func (s *certificate) Create(ctx context.Context, param *CreateCertificateParam) (*dto.Certificate, error) {
	entity, err := entity.NewCertificateWithCertificateKey(param.Certificate, param.Key)
	if err != nil {
		return nil, err
	}

	resp, err := s.r.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return dto.NewCertificate(resp), nil
}

type GetCertificateParam struct {
	Id     string   `json:"id" uri:"id" binding:"required"`
	Expand []string `json:"expand" form:"expand"`
}

func (s *certificate) Get(ctx context.Context, param *GetCertificateParam) (*dto.Certificate, error) {
	ent, err := s.r.Get(ctx, identity.Parse(constant.CertificatePrefix, param.Id))
	if err != nil {
		return nil, err
	}
	obj := dto.NewCertificate(ent)
	return obj, nil
}

type DeleteCertificateParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
}

func (s *certificate) Delete(ctx context.Context, param *DeleteCertificateParam) error {
	err := s.r.Delete(ctx, identity.Parse(constant.CertificatePrefix, param.Id))
	if err != nil {
		return err
	}
	return nil
}

type ListCertificateParam struct {
	Name          string   `form:"name"`
	Limit         int      `form:"limit" binding:"max=1000"`
	StartingAfter string   `form:"starting_after"`
	EndingBefore  string   `form:"ending_before"`
	Expand        []string `json:"expand" form:"expand"`
}

type ListCertificateResult struct {
	HasMore bool               `json:"has_more"`
	Data    []*dto.Certificate `json:"data"`
}

func (s *certificate) List(ctx context.Context, param *ListCertificateParam) (*ListCertificateResult, error) {
	if param.Limit == 0 {
		param.Limit = 10
	}

	entities, err := s.r.List(ctx, &repository.ListCertificateParam{
		Name:          param.Name,
		Limit:         param.Limit,
		StartingAfter: identity.Parse(constant.CertificatePrefix, param.StartingAfter),
		EndingBefore:  identity.Parse(constant.CertificatePrefix, param.EndingBefore),
	})
	if err != nil {
		return nil, err
	}

	n := len(entities)

	items := make([]*dto.Certificate, n)

	for i, v := range entities {
		items[i] = dto.NewCertificate(v)
	}

	rst := &ListCertificateResult{}
	rst.Data = items
	rst.HasMore = n == param.Limit
	return rst, nil
}

type UpdateCertificateParam struct {
	Id string `json:"id" uri:"id" binding:"required"`
	CreateCertificateParam
}

func (s *certificate) Update(ctx context.Context, param *UpdateCertificateParam) (*dto.Certificate, error) {
	id := identity.Parse(constant.CertificatePrefix, param.Id)
	err := s.r.Update(ctx, id, &entity.Certificate{
		Name: param.Name,
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, &GetCertificateParam{Id: param.Id})
}
