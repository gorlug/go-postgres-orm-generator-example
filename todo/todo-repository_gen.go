package todo

// GENERATED FILE
// DO NOT EDIT

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-postgres-generator-example/logger"
	"time"
)

type TodoRepository struct {
	connPool *pgxpool.Pool
	dialect  goqu.DialectWrapper
}

func NewTodoRepository(connPool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		connPool: connPool,
		dialect:  goqu.Dialect("postgres"),
	}
}

func (r *TodoRepository) Create(todo Todo) (int, error) {
	sql, args, err := r.dialect.Insert("todo").
		Prepared(true).
		Rows(goqu.Record{

			"updated_at": time.Now(),
			"name":       todo.Name,
			"checked":    todo.Checked,
			"state":      todo.State,
			"user_id":    todo.UserId,
		}).
		Returning("id").
		ToSQL()
	if err != nil {
		logger.Error("error creating create Todo sql: %v", err)
		return -1, err
	}

	rows, err := r.connPool.Query(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error creating Todo: %v", err)
		return -1, err
	}
	defer rows.Close()
	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			logger.Error("error scanning User: %v", err)
			return -1, err
		}
	} else {
		logger.Error("Todo already exists")
		return -1, TodoAlreadyExistsError{Todo: todo}
	}

	return id, nil

}

type TodoAlreadyExistsError struct {
	Todo Todo
}

func (e TodoAlreadyExistsError) Error() string {
	return fmt.Sprint("Todo ", e.Todo, " already exists")
}

func (r *TodoRepository) GetById(id int) (Todo, error) {
	logger.Debug("Getting Todo by id ", id)
	sql, args, _ := r.dialect.From("todo").
		Prepared(true).
		Select(
			"id",
			"created_at",
			"updated_at",
			"name",
			"checked",
			"state",
			"user_id",
		).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	rows, err := r.connPool.Query(context.Background(), sql, args...)
	if err != nil {
		logger.Error("Failed to get Todo: ", err)
	}
	defer rows.Close()
	item := Todo{}
	for rows.Next() {
		err = rows.Scan(
			&item.Id,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Name,
			&item.Checked,
			&item.State,
			&item.UserId,
		)
		if err != nil {
			logger.Error("Failed to scan Todo: ", err)
			return item, err
		}
	}
	return item, nil
}

func (r *TodoRepository) Update(todo Todo) error {
	sql, args, err := r.dialect.Update("todo").
		Prepared(true).
		Set(goqu.Record{

			"updated_at": time.Now(),
			"name":       todo.Name,
			"checked":    todo.Checked,
			"state":      todo.State,
			"user_id":    todo.UserId,
		}).
		Where(goqu.Ex{"id": todo.Id}).
		ToSQL()
	if err != nil {
		logger.Error("error creating update Todo sql: %v", err)
		return err
	}

	_, err = r.connPool.Exec(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error updating Todo: %v", err)
		return err
	}

	return nil
}

func (r *TodoRepository) Delete(id int) error {
	sql, args, err := r.dialect.Delete("todo").
		Prepared(true).
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		logger.Error("error creating delete Todo sql: %v", err)
		return err
	}

	_, err = r.connPool.Exec(context.Background(), sql, args...)
	if err != nil {
		logger.Error("error deleting Todo: %v", err)
		return err
	}

	return nil
}

func jsonToString(jsonData any) string {
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
