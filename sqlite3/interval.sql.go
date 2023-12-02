// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: interval.sql

package sqlite3

import (
	"context"
	"database/sql"
	"time"
)

const breaks = `-- name: Breaks :many
select id, start_time, planned_duration, actual_duration, category, state from interval where category like '%Break' order by id desc limit ?
`

func (q *Queries) Breaks(ctx context.Context, limit int64) ([]Interval, error) {
	rows, err := q.db.QueryContext(ctx, breaks, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Interval
	for rows.Next() {
		var i Interval
		if err := rows.Scan(
			&i.ID,
			&i.StartTime,
			&i.PlannedDuration,
			&i.ActualDuration,
			&i.Category,
			&i.State,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const byID = `-- name: ByID :one
select id, start_time, planned_duration, actual_duration, category, state from interval where id=?
`

func (q *Queries) ByID(ctx context.Context, id int64) (Interval, error) {
	row := q.db.QueryRowContext(ctx, byID, id)
	var i Interval
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.PlannedDuration,
		&i.ActualDuration,
		&i.Category,
		&i.State,
	)
	return i, err
}

const categorySummary = `-- name: CategorySummary :one
select sum(actual_duration) from interval
where category like ? and strftime('%y-%m-%d', start_time, 'localtime') = strftime('%y-%m-%d', ?2, 'localtime')
`

type CategorySummaryParams struct {
	Category string
	Day      interface{}
}

// The strftime() routine returns the date formatted according to the format string specified as the first argument.
func (q *Queries) CategorySummary(ctx context.Context, arg CategorySummaryParams) (sql.NullFloat64, error) {
	row := q.db.QueryRowContext(ctx, categorySummary, arg.Category, arg.Day)
	var sum sql.NullFloat64
	err := row.Scan(&sum)
	return sum, err
}

const count = `-- name: Count :one
select count(*) from interval
`

func (q *Queries) Count(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, count)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createInterval = `-- name: CreateInterval :one
insert into interval(start_time, planned_duration, actual_duration, category, state) values(?, ?, ?, ?, ?)
returning id, start_time, planned_duration, actual_duration, category, state
`

type CreateIntervalParams struct {
	StartTime       time.Time
	PlannedDuration int64
	ActualDuration  int64
	Category        string
	State           int64
}

func (q *Queries) CreateInterval(ctx context.Context, arg CreateIntervalParams) (Interval, error) {
	row := q.db.QueryRowContext(ctx, createInterval,
		arg.StartTime,
		arg.PlannedDuration,
		arg.ActualDuration,
		arg.Category,
		arg.State,
	)
	var i Interval
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.PlannedDuration,
		&i.ActualDuration,
		&i.Category,
		&i.State,
	)
	return i, err
}

const getIntervals = `-- name: GetIntervals :many
select id, start_time, planned_duration, actual_duration, category, state from interval
`

func (q *Queries) GetIntervals(ctx context.Context) ([]Interval, error) {
	rows, err := q.db.QueryContext(ctx, getIntervals)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Interval
	for rows.Next() {
		var i Interval
		if err := rows.Scan(
			&i.ID,
			&i.StartTime,
			&i.PlannedDuration,
			&i.ActualDuration,
			&i.Category,
			&i.State,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const last = `-- name: Last :one
select id, start_time, planned_duration, actual_duration, category, state from interval order by id desc limit 1
`

func (q *Queries) Last(ctx context.Context) (Interval, error) {
	row := q.db.QueryRowContext(ctx, last)
	var i Interval
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.PlannedDuration,
		&i.ActualDuration,
		&i.Category,
		&i.State,
	)
	return i, err
}

const updateInterval = `-- name: UpdateInterval :exec
update interval set start_time=?, actual_duration=?, state=? where id=?
`

type UpdateIntervalParams struct {
	StartTime      time.Time
	ActualDuration int64
	State          int64
	ID             int64
}

func (q *Queries) UpdateInterval(ctx context.Context, arg UpdateIntervalParams) error {
	_, err := q.db.ExecContext(ctx, updateInterval,
		arg.StartTime,
		arg.ActualDuration,
		arg.State,
		arg.ID,
	)
	return err
}
