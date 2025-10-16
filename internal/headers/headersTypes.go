package headers

type Headers map[string]string

type errString string

func (e errString) Error() string {
	return string(e)
}

const InvaidHeaderSpacingError = errString("invalid header spacing")

const InvalidHeaderNameError = errString("invalid header name")
