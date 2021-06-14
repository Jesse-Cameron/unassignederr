package a

type naughtyError struct{}

func (ne naughtyError) Error() string {
	return "oh no"
}

func MakePointerError() error {
	// todo makes sure this works with non pointer errors
	// todo make sure this works with parenthesis declarations: https://stackoverflow.com/questions/35830676/what-is-this-parenthesis-enclosed-variable-declaration-syntax-in-go/35830718
	var badErr *naughtyError
	if badErr != nil {
		panic("should not execute")
	}
	// todo make a test case where err is overwritten and make sure that is valid

	// returns a `nil` error
	return badErr // want "uninitialized custom error returned \"badErr\""
}

func MakeError() error {
	var badErr naughtyError
	// returns a `nil` error
	return badErr // want "uninitialized custom error returned \"badErr\""
}
