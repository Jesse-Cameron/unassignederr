package a

type naughtyError struct {
	msg string
}

func (ne naughtyError) Error() string {
	return ne.msg
}

func MakePointerError() error {
	// todo make sure this works with parenthesis declarations: https://stackoverflow.com/questions/35830676/what-is-this-parenthesis-enclosed-variable-declaration-syntax-in-go/35830718
	var badErrPointer *naughtyError
	if badErrPointer != nil {
		panic("should not execute")
	}
	// todo make a test case where err is overwritten and make sure that is valid

	// returns a `nil` error
	return badErrPointer // want "uninitialized custom error returned \"badErrPointer\""
}

func MakeError() error {
	var badErr naughtyError
	// returns a `nil` error
	return badErr // want "uninitialized custom error returned \"badErr\""
}

func ErrorAssigned() error {
	var barErrReassigned naughtyError

	if true {
		barErrReassigned = naughtyError{msg: "error"}
	}

	return barErrReassigned
}

func MultiErrorAssigned() error {
	var barErr, fooErr *naughtyError

	fooErr = &naughtyError{msg: "error"}
	if fooErr != nil {
		return fooErr
	}

	return barErr // want "uninitialized custom error returned \"barErr\""
}

func MakeErrorParenDecl() error {
	var (
		badErrParens *naughtyError
		otherType    int
	)

	if otherType != 0 {
		panic("here")
	}

	return badErrParens // want "uninitialized custom error returned \"badErrParens\""
}

func MakeErrorListDecl() error {
	var badErrOne, badErrTwo, badErrThree *naughtyError

	if badErrOne != badErrThree {
		panic("here")
	}

	return badErrTwo // want "uninitialized custom error returned \"badErrTwo\""
}
