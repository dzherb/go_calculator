package orchestrator

import "errors"

var invalidRequestBodyError = errors.New("invalid request body")
var invalidIdInUrlError = errors.New("invalid id in url")

var expressionNotFoundError = errors.New("expression not found")

var taskNotFoundError = errors.New("task not found")
var noTasksToProcessError = errors.New("no tasks to process")
