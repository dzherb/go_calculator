package calculator

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var taskIdSeries = atomic.Uint64{}

// Task - структура задачи для вычисления
type Task struct {
	Id   uint64
	Node *operatorNode
}

func newTask(node *operatorNode) *Task {
	return &Task{
		Id:   taskIdSeries.Add(1),
		Node: node,
	}
}

func (t *Task) Complete(result float64) {
	// Найти родителя и заменить текущий узел на numberNode
	if t.Node.parent != nil {
		if t.Node.parent.left == t.Node {
			t.Node.parent.left = &numberNode{value: result}
		} else if t.Node.parent.right == t.Node {
			t.Node.parent.right = &numberNode{value: result}
		}
	} else {
		// Это корневой узел, заменяем его содержимое
		*t.Node = operatorNode{left: &numberNode{value: result}, processed: true}
	}
}

// compute - выполняет вычисление задачи и обновляет AST
func compute(left, right float64, operator string) (float64, error) {
	var result float64
	switch operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = left / right
	default:
		return 0, fmt.Errorf("invalid or unsupported operator: %s", operator)
	}

	return result, nil
}

var ExpressionIdSeries = atomic.Uint64{}

type Expression struct {
	Id   uint64
	Root *operatorNode
	mu   sync.Mutex
}

func NewExpression(node *operatorNode) *Expression {
	return &Expression{
		Id:   ExpressionIdSeries.Add(1),
		Root: node,
	}
}

func (e *Expression) String() string {
	return fmt.Sprintf("( #%d %s )", e.Id, e.Root.String())
}

func (e *Expression) GetNextTask() (*Task, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, ok := e.Root.nextReadyForProcessingNode()
	if !ok {
		return nil, false
	}
	return newTask(node), true
}

func (e *Expression) IsEvaluated() bool {
	return e.Root.processed
}

func (e *Expression) GetResult() (float64, error) {
	if !e.IsEvaluated() {
		return 0, fmt.Errorf("expressions is not evaluated")
	}
	if resultNode, ok := e.Root.left.(*numberNode); ok {
		return resultNode.value, nil
	}
	return 0, errors.New("expression result node is not a number")
}

// simpleEvaluation - параллельный цикл вычислений
func simpleEvaluation(exp *Expression) error {
	var wg sync.WaitGroup
	errChan := make(chan error)
	isFinished := make(chan bool)
	defer close(errChan)
	defer close(isFinished)

	for {
		task, ok := exp.GetNextTask()
		if !ok {
			// Ждём, если задач временно нет, но AST еще не завершен
			time.Sleep(10 * time.Millisecond)
			task, ok = exp.GetNextTask()
			if !ok {
				break // Выход, если задач больше нет
			}
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			left := task.Node.left.(*numberNode).value
			right := task.Node.right.(*numberNode).value
			result, err := compute(left, right, task.Node.operator)
			if err != nil {
				errChan <- err
			}
			task.Complete(result)
		}()
	}

	go func() {
		wg.Wait()
		isFinished <- true
	}()

	select {
	case err := <-errChan:
		return err
	case <-isFinished:
		return nil
	}
}
