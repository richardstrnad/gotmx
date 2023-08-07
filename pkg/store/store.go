package store

import "github.com/richardstrnad/gotmx/pkg/gotmx"

type Store interface {
	GetTask(id int) (gotmx.Task, error)
}
