package constant

import "errors"

type Response int

const (
	Success Response = iota + 1
	InvalidRequest
	Conflict
	NotFound
	InternalError
)

var (
	InvalidRequestError = errors.New("invalid request")
	ConflictError       = errors.New("conflict")
	NotFoundError       = errors.New("not found")
	InternalErrorError  = errors.New("internal error")
)

func (r Response) GetError() error {
	switch r {
	case Success:
		return nil
	case InvalidRequest:
		return InvalidRequestError
	case Conflict:
		return ConflictError
	case NotFound:
		return NotFoundError
	case InternalError:
		return InternalErrorError
	default:
		return nil
	}
}
