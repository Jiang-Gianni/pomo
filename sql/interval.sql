create table if not exists interval (
    id integer not null,
    start_time datetime not null,
    planned_duration integer not null default 0,
    actual_duration integer not null default 0,
    category text not null,
    state integer not null default 1,
    primary key(id)
);

-- name: GetIntervals :many
select * from interval;

-- name: CreateInterval :one
insert into interval(start_time, planned_duration, actual_duration, category, state) values(?, ?, ?, ?, ?)
returning *;

-- name: UpdateInterval :exec
update interval set start_time=?, actual_duration=?, state=? where id=?;

-- name: ByID :one
select * from interval where id=?;

-- name: Count :one
select count(*) from interval;

-- name: Last :one
select * from interval order by id desc limit 1;

-- name: Breaks :many
select * from interval where category like '%Break' order by id desc limit ?;


-- The strftime() routine returns the date formatted according to the format string specified as the first argument.
-- name: CategorySummary :one
select sum(actual_duration) from interval
where category like ? and strftime('%y-%m-%d', start_time, 'localtime') = strftime('%y-%m-%d', sqlc.arg(day), 'localtime')