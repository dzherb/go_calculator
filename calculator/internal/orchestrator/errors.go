package orchestrator

import "errors"

var errInvalidRequestBody = errors.New("invalid request body")
var errInvalidIdInUrl = errors.New("invalid id in url")

var errExpressionNotFound = errors.New("expression not found")

var errTaskNotFound = errors.New("task not found")
var errNoTasksToProcess = errors.New("no tasks to process")
