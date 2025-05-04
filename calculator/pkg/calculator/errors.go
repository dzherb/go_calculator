package calc

import "errors"

var ErrTaskIsCompleted = errors.New("task is already completed")
var ErrTaskIsCanceled = errors.New("task was canceled")
