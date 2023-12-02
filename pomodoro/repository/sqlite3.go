//go:build !inmemory
// +build !inmemory

package repository

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/Jiang-Gianni/pomo/pomodoro"
	"github.com/Jiang-Gianni/pomo/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableInterval = `create table if not exists interval (
		id integer,
		start_time datetime not null,
		planned_duration integer default 0,
		actual_duration integer default 0,
		category text not null,
		state integer default 1,
		primary key(id)
	);`
)

type dbRepo struct {
	db    *sql.DB
	query *sqlite3.Queries
	sync.RWMutex
}

func NewSQLite3Repo(dbfile string) (*dbRepo, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &dbRepo{
		db:    db,
		query: sqlite3.New(db),
	}, nil
}

func (r *dbRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := sqlite3.CreateIntervalParams{
		StartTime:       i.StartTime,
		PlannedDuration: int64(i.PlannedDuration),
		ActualDuration:  int64(i.ActualDuration),
		Category:        i.Category,
		State:           int64(i.State),
	}

	interval, err := r.query.CreateInterval(ctx, params)
	if err != nil {
		return -1, err
	}

	return interval.ID, nil
}

func (r *dbRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := sqlite3.UpdateIntervalParams{
		ID:             i.ID,
		StartTime:      i.StartTime,
		ActualDuration: int64(i.ActualDuration),
		State:          int64(i.State),
	}

	return r.query.UpdateInterval(ctx, params)
}

func (r *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i, err := r.query.ByID(ctx, id)

	return pomodoro.Interval{
		ID:              i.ID,
		StartTime:       i.StartTime,
		PlannedDuration: time.Duration(i.PlannedDuration),
		ActualDuration:  time.Duration(i.ActualDuration),
		Category:        i.Category,
		State:           int(i.State),
	}, err
}

func (r *dbRepo) Count() (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count, err := r.query.Count(ctx)
	if err != nil {
		return -1, err
	}

	return int(count), nil
}

func (r *dbRepo) Last() (pomodoro.Interval, error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count, err := r.Count()
	if err != nil {
		return pomodoro.Interval{}, err
	}

	if count == 0 {
		return pomodoro.Interval{}, pomodoro.ErrNoIntervals
	}

	i, err := r.query.Last(ctx)

	return pomodoro.Interval{
		ID:              i.ID,
		StartTime:       i.StartTime,
		PlannedDuration: time.Duration(i.PlannedDuration),
		ActualDuration:  time.Duration(i.ActualDuration),
		Category:        i.Category,
		State:           int(i.State),
	}, err
}

func (r *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	intervals, err := r.query.Breaks(ctx, int64(n))
	out := make([]pomodoro.Interval, len(intervals))

	for i := range intervals {
		out[i] = pomodoro.Interval{
			ID:              intervals[i].ID,
			StartTime:       intervals[i].StartTime,
			PlannedDuration: time.Duration(intervals[i].PlannedDuration),
			ActualDuration:  time.Duration(intervals[i].ActualDuration),
			Category:        intervals[i].Category,
			State:           int(intervals[i].State),
		}
	}

	return out, err
}

func (r *dbRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := sqlite3.CategorySummaryParams{
		Category: filter,
		Day:      day,
	}
	nf, err := r.query.CategorySummary(ctx, params)

	var d time.Duration

	if nf.Valid {
		d = time.Duration(nf.Float64)
	}

	return d, err
}
