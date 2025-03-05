package calculation

import "errors"

var (
	ErrMismatchedBracket          = errors.New("mismatched bracket")
	ErrInvalidSymbols             = errors.New("invalid symbols")
	ErrInvalidOperationsPlacement = errors.New("invalid operations placement")
	ErrZeroDivision               = errors.New("division by zero")
	ErrInvalidExpression          = errors.New("invalid expression")
	ErrCalculation 				  = errors.Join(
										ErrInvalidExpression,
										ErrInvalidOperationsPlacement,
										ErrInvalidSymbols,
										ErrMismatchedBracket,
										ErrZeroDivision,
									)
)
