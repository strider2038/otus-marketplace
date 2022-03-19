package trading

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrItemNotFound  = errors.New("item not found")
	ErrDenied        = errors.New("operation denied")
	ErrCannotCancel  = errors.New("cannot cancel order")
	ErrItemIsOnSale  = errors.New("item is on sale already")
)

type UnexpectedStatusError struct {
	ActualStatus   string
	ExpectedStatus string
}

func newUnexpectedStatusError(actual, expected string) error {
	return errors.WithStack(UnexpectedStatusError{
		ActualStatus:   actual,
		ExpectedStatus: expected,
	})
}

func (err UnexpectedStatusError) Error() string {
	return fmt.Sprintf(`unexpected status "%s", expected is "%s"`, err.ActualStatus, err.ExpectedStatus)
}
