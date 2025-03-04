package orchestrator

import (
	"go_calculator/pkg/calculator"
	"sync"
)

type Storage[T any] interface {
	Put(value T)
	Get(id uint64) (T, bool)
	GetAll() []T
}

type expressionStorage struct {
	expressions map[uint64]*calculator.Expression
	mu          sync.RWMutex
}

func (s *expressionStorage) Put(expression *calculator.Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expressions[expression.Id] = expression
}

func (s *expressionStorage) Get(id uint64) (*calculator.Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exp, ok := s.expressions[id]
	return exp, ok
}

func (s *expressionStorage) GetAll() []*calculator.Expression {
	s.mu.RLock()
	defer s.mu.RUnlock()
	expressions := make([]*calculator.Expression, 0, len(s.expressions))
	for _, exp := range s.expressions {
		expressions = append(expressions, exp)
	}
	return expressions
}

var ExpressionStorageInstance = &expressionStorage{
	expressions: make(map[uint64]*calculator.Expression),
}

type taskStorage struct {
	tasks map[uint64]*calculator.Task
	mu    sync.RWMutex
}

func (s *taskStorage) Put(task *calculator.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.Id] = task
}

func (s *taskStorage) Get(id uint64) (*calculator.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

func (s *taskStorage) GetAll() []*calculator.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*calculator.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

var TaskStorageInstance = &taskStorage{
	tasks: make(map[uint64]*calculator.Task),
}
