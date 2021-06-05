package sequence

import "fmt"

type ErrSequenceNotRegistered struct {
	error
	SequenceKey string
}

func (e *ErrSequenceNotRegistered) Error() string {
	return fmt.Sprintf("sequence %s not registered", e.SequenceKey)
}


