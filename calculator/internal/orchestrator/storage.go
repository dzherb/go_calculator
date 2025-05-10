package orchestrator

import (
	"iter"
	"sync"

	"github.com/dzherb/go_calculator/pkg/calculator"
)

type Storage[T any] interface {
	Put(value T)
	Get(id uint64) (T, bool)
	Delete(id uint64)
	All() iter.Seq[T]
}

type exprStorage struct {
	expressions map[uint64]*calc.Expression
	mu          sync.RWMutex
}

func (s *exprStorage) Put(expression *calc.Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expressions[expression.Id] = expression
}

func (s *exprStorage) Get(id uint64) (*calc.Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exp, ok := s.expressions[id]

	return exp, ok
}

func (s *exprStorage) All() iter.Seq[*calc.Expression] {
	return func(yield func(*calc.Expression) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for _, exp := range s.expressions {
			if !yield(exp) {
				return
			}
		}
	}
}

func (s *exprStorage) Delete(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.expressions, id)
}

var ExpressionStorageInstance = &exprStorage{
	expressions: make(map[uint64]*calc.Expression),
}

type taskStorage struct {
	tasks map[uint64]*calc.Task
	mu    sync.RWMutex
}

func (s *taskStorage) Put(task *calc.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.Id] = task
}

func (s *taskStorage) Get(id uint64) (*calc.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]

	return task, ok
}

func (s *taskStorage) All() iter.Seq[*calc.Task] {
	return func(yield func(*calc.Task) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for _, exp := range s.tasks {
			if !yield(exp) {
				return
			}
		}
	}
}

func (s *taskStorage) Delete(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, id)
}

var TaskStorageInstance = &taskStorage{
	tasks: make(map[uint64]*calc.Task),
}
