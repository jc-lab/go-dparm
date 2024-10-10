package tcg

import (
	"fmt"
)

type TcgError struct {
	Status MethodStatus
}

func (e *TcgError) Error() string {
	return fmt.Sprintf("tcg error: %d", e.Status)
}
