package calculation

var (
	ValidTestSet = []struct {
		Name            string
		Expression      string
		Expected_answer []string
	}{
		{
			Name: "Valid sum expression",
			Expression: "2+2",
			Expected_answer: []string{"2", "2", "+"},
		},
		{
			Name: "Valid sub expression",
			Expression: "2-2",
			Expected_answer: []string{"2", "2", "-"},
		},
		{
			Name: "Valid mul expression",
			Expression: "2*6",
			Expected_answer: []string{"2", "6", "*"},
		},
		{
			Name: "Valid div expression",
			Expression: "12/3",
			Expected_answer: []string{"12", "3", "/"},
		},
		{
			Name: "Valid expression with rational numbers",
			Expression: "2.2/1.1",
			Expected_answer: []string{"2.2", "1.1", "/"},
		},
		{
			Name: "Valid expression with spaces",
			Expression: "2 + 6",
			Expected_answer: []string{"2", "6", "+"},
		},
		{
			Name: "Valid expression with brackets",
			Expression: "2*(2+2)",
			Expected_answer: []string{"2", "2", "2", "+", "*"},
		},
		{
			Name: "Valid expression with negative nums",
			Expression: "-2*(-4+2)",
			Expected_answer: []string{"-2", "-4", "2", "+", "*"},
		},
	}
	InvalidTestSet = []struct {
		Name 			string
		Expression 		string
		Expected_error  error
	}{
		{
			Name: "Invalid brackets 1",
			Expression: "2+2+2)",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 2",
			Expression: "2+2+2(()",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 3",
			Expression: "2+2+2(",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 4",
			Expression: "2+((2+2",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 5",
			Expression: "(2+2+2",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 6",
			Expression: ")2+2+2",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid brackets 7",
			Expression: "2+2+2)",
			Expected_error: ErrMismatchedBracket,
		},
		{
			Name: "Invalid symbols 1",
			Expression: "a",
			Expected_error: ErrInvalidSymbols,
		},
		{
			Name: "Invalid symbols 2",
			Expression: "2+O",
			Expected_error: ErrInvalidSymbols,
		},
		{
			Name: "Invalid symbols 3",
			Expression: "3^3",
			Expected_error: ErrInvalidSymbols,
		},
		{
			Name: "Invalid symbols 4",
			Expression: "2|2",
			Expected_error: ErrInvalidSymbols,
		},
		{
			Name: "Invalid operations placement 1",
			Expression: "2++",
			Expected_error: ErrInvalidOperationsPlacement,
		},
		{
			Name: "Invalid operations placement 2",
			Expression: "2++2-",
			Expected_error: ErrInvalidOperationsPlacement,
		},
		{
			Name: "Invalid operations placement 3",
			Expression: "2*(*2+2)",
			Expected_error: ErrInvalidOperationsPlacement,
		},
		{
			Name: "Invalid expression 1",
			Expression: "2-2-",
			Expected_error: ErrInvalidExpression,
		},
		{
			Name: "Invalid expression 2",
			Expression: "2-",
			Expected_error: ErrInvalidExpression,
		},
		{
			Name: "Invalid expression 3",
			Expression: "3+3/",
			Expected_error: ErrInvalidExpression,
		},
		{
			Name: "Invalid expression 4",
			Expression: "3+3/",
			Expected_error: ErrInvalidExpression,
		},
	}
)