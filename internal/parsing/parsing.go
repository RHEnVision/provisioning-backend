package parsing

func IsHTTPStatus2xx(status int) bool {
	return status >= 200 && status <= 299
}

func IsHTTPNotFound(status int) bool {
	return status == 404
}
