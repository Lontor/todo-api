package custom_errors

const (
	CodeNotFound = "NotFound"
	CodeDBError  = "DBError"
)

type RepositoryError struct {
	Code    string
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}

func NewRepositoryError(code, message string) *RepositoryError {
	return &RepositoryError{
		Code:    code,
		Message: message,
	}
}
