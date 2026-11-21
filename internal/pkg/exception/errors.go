package exception

type ServerError struct {
	Message string
}

func (e *ServerError) Error() string {
	return e.Message
}

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}

type AuthorizationError struct {
	Message string
}

func (e *AuthorizationError) Error() string {
	return e.Message
}
