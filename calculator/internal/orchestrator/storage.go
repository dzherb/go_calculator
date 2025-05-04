package orchestrator

import (
	"sync"

	"github.com/dzherb/go_calculator/pkg/calculator"
)

type Storage[T any] interface {
	Put(value T)
	Get(id uint64) (T, bool)
	GetAll() []T
}

type expressionStorage struct {
	expressions map[uint64]*calc.Expression
	mu          sync.RWMutex
}

func (s *expressionStorage) Put(expression *calc.Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expressions[expression.Id] = expression
}

func (s *expressionStorage) Get(id uint64) (*calc.Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exp, ok := s.expressions[id]

	return exp, ok
}

func (s *expressionStorage) GetAll() []*calc.Expression {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expressions := make([]*calc.Expression, 0, len(s.expressions))

	for _, exp := range s.expressions {
		expressions = append(expressions, exp)
	}

	return expressions
}

var ExpressionStorageInstance = &expressionStorage{
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

func (s *taskStorage) GetAll() []*calc.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*calc.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

var TaskStorageInstance = &taskStorage{
	tasks: make(map[uint64]*calc.Task),
}
