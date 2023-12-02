//go:build inmemory
// +build inmemory

package cmd

import (
	"github.com/Jiang-Gianni/pomo/pomodoro"
	"github.com/Jiang-Gianni/pomo/pomodoro/repository"
	_ "github.com/mattn/go-sqlite3"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
