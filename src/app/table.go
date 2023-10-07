package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kallepan/redisql/constant"
	"github.com/redis/go-redis/v9"
)

type Table struct {
	/**
	 * Table represents a table in a database.
	 **/
	Name string
}

type TableImpl interface {
	/**
	 * TableImpl represents a table in a database.
	 **/
	AddRow(ctx context.Context, rdb *redis.Client, id string, key string, value interface{}) error
	DeleteRow(ctx context.Context, rdb *redis.Client, id string, key string) error
	AddRows(ctx context.Context, rdb *redis.Client, id string, data map[string]interface{}) error
	GetRow(ctx context.Context, rdb *redis.Client, id string) (map[string]interface{}, error)
	GetTable(ctx context.Context, rdb *redis.Client) ([]map[string]interface{}, error)
}

func TableImplInit(table string) TableImpl {
	/**
	 * TableImplInit initializes a TableImpl.
	 */
	return &Table{
		Name: table,
	}
}

func (t Table) AddRow(ctx context.Context, rdb *redis.Client, id string, key string, value interface{}) error {
	/**
	 * AddRow adds a row to a table.
	 */
	identifier := fmt.Sprintf("%s:%s", t.Name, id)

	// Check if row already exists.
	if _, err := rdb.HGet(ctx, identifier, key).Result(); err != redis.Nil && err != nil {
		// Row does not exist.
		slog.Error(err.Error())
		return constant.InternalError.GetError()
	} else if err != redis.Nil {
		// Row already exists.
		return constant.Conflict.GetError()
	}

	// Add row to the table.
	if err := rdb.HSet(ctx, identifier, key, value).Err(); err != nil {
		slog.Error(err.Error())
		return constant.InternalError.GetError()
	}

	return constant.Success.GetError()
}

func (t Table) DeleteRow(ctx context.Context, rdb *redis.Client, id string, key string) error {
	/**
	 * DeleteRow deletes a row from a table.
	 */
	identifier := fmt.Sprintf("%s:%s", t.Name, id)

	// Check if row exists.
	if _, err := rdb.HGet(ctx, identifier, key).Result(); err != nil {
		return constant.NotFound.GetError()
	}

	// Delete row from the table.
	if err := rdb.HDel(ctx, identifier, key).Err(); err != nil {
		slog.Error(err.Error())
		return constant.InternalError.GetError()
	}

	return constant.Success.GetError()
}

func (t Table) AddRows(ctx context.Context, rdb *redis.Client, id string, data map[string]interface{}) error {
	/**
	 * AddRows adds multiple rows to a table.
	 */
	identifier := fmt.Sprintf("%s:%s", t.Name, id)

	if _, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for k, v := range data {
			pipe.HSet(ctx, identifier, k, v)
		}

		_, err := pipe.Exec(ctx)
		return err
	}); err != nil {
		slog.Error(err.Error())
		return constant.InternalError.GetError()
	}

	return constant.Success.GetError()
}

func (t Table) GetRow(ctx context.Context, rdb *redis.Client, id string) (map[string]interface{}, error) {
	/**
	 * GetRow gets all rows from a table.
	 */
	identifier := fmt.Sprintf("%s:%s", t.Name, id)

	// Get all rows from the table.
	columns, err := rdb.HGetAll(ctx, identifier).Result()
	if err != nil {
		slog.Error(err.Error())
		return nil, constant.InternalError.GetError()
	}

	// Check if all rows are empty.
	if len(columns) == 0 {
		return nil, constant.NotFound.GetError()
	}

	// convert rows to map[string]interface{}
	row := make(map[string]interface{})
	for k, v := range columns {
		row[k] = v
	}

	return row, constant.Success.GetError()
}

func (t Table) GetTable(ctx context.Context, rdb *redis.Client) ([]map[string]interface{}, error) {
	/**
	 * GetRows gets all rows from a table.
	 */
	identifier := fmt.Sprintf("%s:*", t.Name)

	// Get all rows from the table.
	rowKeys, _, err := rdb.Scan(ctx, 0, identifier, 0).Result()
	if err != nil {
		slog.Error(err.Error())
		return nil, constant.InternalError.GetError()
	}
	if len(rowKeys) == 0 {
		return nil, constant.NotFound.GetError()
	}

	// Get all rows from the table.
	results := make([]map[string]interface{}, len(rowKeys))
	if _, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, rowKey := range rowKeys {
			columns, err := pipe.HGetAll(ctx, fmt.Sprintf("%s:%s", t.Name, rowKey)).Result()
			if err != nil {
				return fmt.Errorf("failed to get rows: %w", err)
			}

			// convert rows to map[string]interface{}
			row := make(map[string]interface{})
			for k, v := range columns {
				row[k] = v
			}

			results = append(results, row)
		}
		_, err := pipe.Exec(ctx)
		return err
	}); err != nil {
		slog.Error(err.Error())
		return nil, constant.InternalError.GetError()
	}

	return results, constant.Success.GetError()
}
