package utils

import (
	"net/http"

	"github.com/Lontor/todo-api/pkg/custom_errors"
)

func RepositoryErrorToHTTPError(err error) error {
	if err == nil {
		return nil
	}

	if repoErr, ok := err.(*custom_errors.RepositoryError); ok {
		switch repoErr.Code {
		case custom_errors.CodeNotFound:
			return custom_errors.NewHTTPError(http.StatusNotFound, repoErr.Message)
		case custom_errors.CodeDBError:
			return custom_errors.NewHTTPError(http.StatusInternalServerError, repoErr.Message)
		default:
			return custom_errors.NewHTTPError(http.StatusInternalServerError, repoErr.Message)
		}
	}
	return err
}
