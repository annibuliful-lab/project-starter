package error_utils

import "fmt"

type GraphqlError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e GraphqlError) Error() string {
	if e.Code == "" {
		e.Code = "Error"
	}

	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

func (e GraphqlError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}
