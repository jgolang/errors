package errors_test

import (
	stderrors "errors"
	"fmt"

	jerrors "github.com/jgolang/errors"
	"github.com/jgolang/errors/codes"
)

var ErrTitleRequired = jerrors.Define(codes.New("PLN-001", "title is required"), "title is required")

func validateTitle(title string) error {
	if title == "" {
		return jerrors.Wrap(ErrTitleRequired)
	}
	return nil
}

func ExampleDefine() {
	err := validateTitle("")

	fmt.Println(stderrors.Is(err, ErrTitleRequired))
	fmt.Println(jerrors.CodeOf(err).Str())
	// Output:
	// true
	// PLN-001
}
