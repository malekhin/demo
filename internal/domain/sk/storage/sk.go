package storage

import (
	"context"
	"demo/internal/common"
	serviceModel "demo/internal/domain/sk/model"
	"demo/internal/domain/sk/storage/model"
	"fmt"

	"demo/internal/common/wrapper/database"
)

type Sk struct {
	Db database.Db
}

func NewSk(db database.Database) *Sk {
	return &Sk{
		Db: db.Db(),
	}
}

func (s *Sk) List(ctx context.Context, filter serviceModel.SkFilter) ([]model.Sk, int, error) {
	sql := `
		select id, name, is_active from sk order by id asc
	`
	params := map[string]interface{}{}

	result := make([]model.Sk, 0)
	err := s.Db.NamedRows(ctx, &result, sql+fmt.Sprintf(" offset %d limit %d", filter.Offset, filter.Limit), params)
	if err != nil {
		return nil, 0, common.Wrap(err, "NamedRows")
	}

	var count int
	err = s.Db.NamedRow(ctx, &count, `select count(*) from (`+sql+`) c`, params)
	if err != nil {
		return nil, 0, common.Wrap(err, "NamedRow")
	}

	return result, count, nil
}

func (s *Sk) Add(ctx context.Context, sk model.Sk) (*model.Sk, error) {
	var res model.Sk
	err := s.Db.NamedRow(ctx, &res, `
		insert into sk
		    (id, name, is_active, created_at, updated_at)
		values 
			(:id, :name, :is_active, now(), now())
		returning id, name, is_active, created_at, updated_at
	`, sk)
	if err != nil {
		return nil, common.Wrap(err, "NamedRow")
	}

	return &res, nil
}

func (s *Sk) Edit(ctx context.Context, sk model.Sk) (*model.Sk, error) {
	var res model.Sk
	err := s.Db.NamedRow(ctx, &res, `
		update sk 
		set 
			name = :name, 
			is_active = :is_active, 
			updated_at = now()
		where id = :id
		returning id, name, is_active, created_at, updated_at
	`, sk)
	if err != nil {
		return nil, common.Wrap(err, "NamedRow")
	}

	return &res, nil
}

func (s *Sk) IsDuplicate(ctx context.Context, sk model.Sk) (bool, error) {
	var exists bool
	err := s.Db.Row(ctx, &exists, `select count(*) from sk where id = ? or name = ?`, sk.Id, sk.Name)
	if err != nil {
		return false, common.Wrap(err, "NamedRow")
	}

	return exists, nil
}
