package calculator

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var taskIdSeries = atomic.Uint64{}

type Task struct {
	Id          uint64
	node        *operatorNode
	IsCompleted bool
	IsCancelled bool
	mu          sync.Mutex
}

func newTask(node *operatorNode) *Task {
	return &Task{
		Id:   taskIdSeries.Add(1),
		node: node,
	}
}

func (t *Task) GetArguments() (float64, float64) {
	return t.node.left.(*numberNode).value, t.node.right.(*numberNode).value
}

func (t *Task) GetOperator() string {
	return t.node.operator
}

func (t *Task) Complete(result float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.IsCompleted {
		return errors.New("task is already completed")
	}
	if t.IsCancelled {
		return errors.New("task was cancelled")
	}

	// Найти родителя и заменить текущий узел на numberNode
	if t.node.parent != nil {
		if t.node.parent.left == t.node {
			t.node.parent.left = &numberNode{value: result}
		} else if t.node.parent.right == t.node {
			t.node.parent.right = &numberNode{value: result}
		}
	} else {
		// Это корневой узел, заменяем его содержимое
		*t.node = operatorNode{left: &numberNode{value: result}, processed: true}
	}
	t.IsCompleted = true
	return nil
}

func (t *Task) Cancel() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.IsCompleted {
		return errors.New("task is completed and cannot be cancelled")
	}

	t.IsCancelled = true
	t.node.processed = false
	return nil
}

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
	Id           uint64
	Root         *operatorNode
	IsProcessing bool
	mu           sync.RWMutex
}

func NewExpression(expression string) (*Expression, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return nil, err
	}

	// Переводим токены в обратную польскую нотацию (RPN)
	rpnOrganizedTokens := shuntingYard(tokens)
	// Составлем абстрактное синтактическое дерево
	ast := buildAST(rpnOrganizedTokens)

	var root *operatorNode

	switch n := ast.(type) {
	case *operatorNode:
		root = n
	case *numberNode:
		root = &operatorNode{left: &numberNode{value: n.value}, processed: true}
	}

	return &Expression{
		Id:   ExpressionIdSeries.Add(1),
		Root: root,
	}, nil
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
	e.IsProcessing = true
	return newTask(node), true
}

func (e *Expression) IsEvaluated() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.Root.processed
}

func (e *Expression) GetResult() (float64, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

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
			left, right := task.GetArguments()
			result, err := compute(left, right, task.node.operator)
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
