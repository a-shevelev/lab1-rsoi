package repo

import (
	"context"
	"fmt"
	"lab1-rsoi/internal/entity"
	"lab1-rsoi/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type PersonRepository interface {
	Create(context.Context, *entity.Person) (uint64, error)
	Fetch(context.Context, uint64) (entity.Person, error)
	FetchAll(context.Context) ([]entity.Person, error)
	Update(context.Context, entity.Person) error
	Delete(context.Context, uint64) error
}
type personRepo struct {
	conn postgres.Connection
}

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func New(client postgres.Client) PersonRepository {
	return &personRepo{conn: client.Conn()}
}

func (r *personRepo) Create(ctx context.Context, p *entity.Person) (uint64, error) {
	query, args, err := qb.Insert("person").
		Columns("name", "age", "address", "work").
		Values(p.Name, p.Age, p.Address, p.Work).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[uint64])
}

func (r *personRepo) FetchAll(ctx context.Context) ([]entity.Person, error) {
	query := qb.Select("id", `"name"`, "age", "address", `"work"`).From("person")

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons, err := pgx.CollectRows[entity.Person](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *personRepo) Fetch(ctx context.Context, id uint64) (entity.Person, error) {
	query := qb.Select("id", `"name"`, "age", "address", "work").
		From("person").
		Where("id = ?", id)
	sql, args, err := query.ToSql()
	if err != nil {
		log.WithError(err).Error("failed to build query")
		return entity.Person{}, nil
	}
	//fmt.Printf(sql, args)

	rows, err := r.conn.Query(ctx, sql, args...)
	return pgx.CollectOneRow[entity.Person](rows, pgx.RowToStructByNameLax)

}

func (r *personRepo) Update(ctx context.Context, p entity.Person) error {
	query := qb.Update("person").
		Set("address", p.Address).
		Set("age", p.Age).
		Set("work", p.Work).
		Where("id = ?", p.ID)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	fmt.Println(p.Name, *p.Age, *p.Address, *p.Work)
	fmt.Println(query, args)
	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *personRepo) Delete(ctx context.Context, id uint64) error {
	query := qb.Delete("person").Where("id = ?", id)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
