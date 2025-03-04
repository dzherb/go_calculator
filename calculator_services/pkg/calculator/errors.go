package calculator

import "errors"

var TaskIsCompletedError = errors.New("task is already completed")
var TaskIsCanceledError = errors.New("task was canceled")
