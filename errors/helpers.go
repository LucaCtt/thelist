package errors

func E(args ...interface{}) error {
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Message = arg
		case ErrorCode:
			e.Code = arg
		case ErrorSeverity:
			e.Severity = arg
		case error:
			e.Err = arg
		default:
			panic("invalid arg")
		}
	}

	return e
}

func Code(err error) ErrorCode {
	e, ok := err.(*Error)
	if !ok {
		return CodeUnexpected
	}

	if e.Code != 0 {
		return e.Code
	}

	return Code(e.Err)
}
