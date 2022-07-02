package utils

type DBError struct {
	Msg  string
	Code int
}

func (e *DBError) Error() string {
	return e.Msg
}

func (e *DBError) StatusCode() int {
	return e.Code
}
