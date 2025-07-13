package errctx

import (
	"errors"
	"fmt"
)

type ErrCtx struct {
	err      error
	metadata map[string]any
}

func Errorf(msg string, e error) ErrCtx {
	var err ErrCtx
	if errors.As(e, &err) {
		err.err = fmt.Errorf(msg, err.err)

		return err
	}

	return ErrCtx{
		metadata: make(map[string]any),
		err:      fmt.Errorf(msg, e),
	}
}

func Is(err, target error) bool {
	var errCtx ErrCtx
	if errors.As(err, &errCtx) {
		return Is(errCtx.err, target)
	}

	return errors.Is(err, target)
}

func As(err error, target any) bool {
	var errCtx ErrCtx
	if errors.As(err, &errCtx) {
		if _, ok := target.(*ErrCtx); ok {
			return true
		}

		return As(errCtx.err, target)
	}

	return errors.As(err, target)
}

func NewFromErr(err error) ErrCtx {
	return ErrCtx{
		metadata: make(map[string]any),
		err:      err,
	}
}

func New(msg string) ErrCtx {
	return ErrCtx{
		metadata: make(map[string]any),
		err:      errors.New(msg),
	}
}

func ErrToCtx(err error) ErrCtx {
	var errCustom ErrCtx
	if errors.As(err, &errCustom) {
		return errCustom
	}

	return ErrCtx{
		metadata: make(map[string]any),
		err:      err,
	}
}

func (e ErrCtx) Value(key string) any {
	return e.metadata[key]
}

func (e ErrCtx) Values() map[string]any {
	return e.metadata
}

func (e ErrCtx) Error() string {
	return e.err.Error()
}

func (e ErrCtx) Join(err error) ErrCtx {
	e.err = errors.Join(err, e.err)

	return e
}

func (e ErrCtx) With(field string, value any) ErrCtx {
	e.metadata[field] = value

	return e
}
