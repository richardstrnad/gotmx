package inmemory

import "fmt"
import "github.com/richardstrnad/gotmx/pkg/gotmx"

type InMemoryDataStore struct {
	tasks map[int]gotmx.Task
}

func NewInMemoryDataStore() *InMemoryDataStore {
	tasks := make(map[int]gotmx.Task)
	tasks[1] = gotmx.Task{
		ID:     1,
		UserID: 1,
		Title:  "First Task",
		Body:   "This is the first task",
	}
	tasks[2] = gotmx.Task{
		ID:     2,
		UserID: 1,
		Title:  "Second Task",
		Body:   "This is the second task",
	}
	tasks[3] = gotmx.Task{
		ID:     3,
		UserID: 1,
		Title:  "Third Task",
		Body:   "This is the third task",
	}
	return &InMemoryDataStore{
		tasks: tasks,
	}
}

func (s *InMemoryDataStore) GetTask(id int) (gotmx.Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return gotmx.Task{}, fmt.Errorf("task with id %d not found", id)
	}
	task.NextID = id + 1
	task.PrevID = id - 1
	task.Target = "#tasks"
	return task, nil
}
