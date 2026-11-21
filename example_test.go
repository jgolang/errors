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

func ExampleKindOf() {
	titleRequired := codes.New("PLN-001", "title is required").WithKind(codes.KindValidation)
	err := fmt.Errorf("saving task: %w", jerrors.NewC(titleRequired, "title is empty"))

	fmt.Println(jerrors.KindOf(err))
	fmt.Println(jerrors.KindOf(jerrors.NewC(codes.Unknown, "no kind")) == "")
	// Output:
	// validation
	// true
}

func ExampleError_Detail() {
	err := jerrors.With(stderrors.New("connection refused"), "loading board")

	var coded *jerrors.Error
	if stderrors.As(err, &coded) {
		fmt.Println(coded.Detail())
	}
	// Output:
	// connection refused
}
