package service

import (
	"context"
	"fmt"
	"lab1-rsoi/internal/dto"
	"lab1-rsoi/internal/entity"
	"lab1-rsoi/internal/repo"
)

type PersonServiceIface interface {
	Create(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error)
	List(ctx context.Context) ([]dto.PersonResponse, error)
	Get(ctx context.Context, id uint64) (*dto.PersonResponse, error)
	Update(ctx context.Context, id uint64, req dto.PersonResponse) (*dto.PersonResponse, error)
	Delete(ctx context.Context, id uint64) error
}
type PersonService struct {
	repo repo.PersonRepository
}

func New(repo repo.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) Create(ctx context.Context, req *dto.CreatePersonRequest) (uint64, error) {
	//if strings.TrimSpace(req.Name) == "" {
	//	return nil, ErrEmptyName
	//}
	//if req.Age != nil && *req.Age <= 0 {
	//	return nil, ErrInvalidAge
	//}
	//if req.Address != nil && len(*req.Address) > 255 {
	//	return nil, ErrFieldTooLong
	//}
	//if req.Work != nil && len(*req.Work) > 255 {
	//	return nil, ErrFieldTooLong
	//}

	p := &entity.Person{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Work:    req.Work,
	}

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PersonService) List(ctx context.Context) ([]dto.PersonResponse, error) {
	persons, err := s.repo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.PersonResponse, len(persons))
	for i, p := range persons {
		resp[i] = dto.PersonResponse{
			ID:      p.ID,
			Name:    &p.Name,
			Age:     p.Age,
			Address: p.Address,
			Work:    p.Work,
		}
	}

	return resp, nil
}

func (s *PersonService) Get(ctx context.Context, id uint64) (*dto.PersonResponse, error) {
	p, err := s.repo.Fetch(ctx, id)
	if err != nil {
		return nil, err
	}
	//if p == nil {
	//	return nil, nil
	//}

	return &dto.PersonResponse{
		ID:      p.ID,
		Name:    &p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}, nil
}

func (s *PersonService) Update(ctx context.Context, id uint64, req dto.PersonResponse) (*dto.PersonResponse, error) {
	p, err := s.repo.Fetch(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Age != nil {
		p.Age = req.Age
	}
	if req.Address != nil {
		p.Address = req.Address
	}
	if req.Work != nil {
		p.Work = req.Work
	}
	fmt.Println(p.Name, *p.Age, *p.Address, *p.Work)

	err = s.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	return &dto.PersonResponse{
		ID:      p.ID,
		Name:    &p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}, nil
}

func (s *PersonService) Delete(ctx context.Context, id uint64) error {
	_, err := s.repo.Fetch(ctx, id)
	if err != nil {
		return err
	}
	//if p == nil {
	//	return nil
	//}
	return s.repo.Delete(ctx, id)
}
