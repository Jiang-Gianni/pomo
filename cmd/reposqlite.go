//go:build !inmemory
// +build !inmemory

package cmd

import (
	"github.com/Jiang-Gianni/pomo/pomodoro"
	"github.com/Jiang-Gianni/pomo/pomodoro/repository"
	"github.com/spf13/viper"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewSQLite3Repo(viper.GetString("db"))
}
