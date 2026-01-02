package must

import (
	"fmt"
	"os"
)

type Controller[T any] interface {
	Fallback(err error) T
}

type ErrorHandler interface {
	Handle(err error)
}

type errorHandler[T any] struct {
	c Controller[T]
}

func (h errorHandler[T]) Handle(err error) {
	h.c.Fallback(err)
}

func HandlerOf[T any](c Controller[T]) ErrorHandler {
	return errorHandler[T]{c}
}

func Do[T any](c Controller[T]) func(func() (T, error)) T {
	return func(f func() (T, error)) T {
		v, err := f()
		if err != nil {
			v = c.Fallback(err)
		}
		return v
	}
}

func Handle(h ErrorHandler) func(func() error) {
	return func(f func() error) {
		if err := f(); err != nil {
			h.Handle(err)
		}
	}
}

func Have[T any](c Controller[T]) func(T, error) T {
	return func(t T, err error) T {
		return Do(c)(func() (T, error) {
			return t, err
		})
	}
}

func ExitHandler(code int) ErrorHandler {
	return HandlerOf(exitController[struct{}]{code: code})
}

func ExitController[T any](code int) Controller[T] {
	return exitController[T]{code: code}
}

type exitController[T any] struct {
	code     int
	fallback T
}

func (c exitController[T]) Fallback(err error) T {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(c.code)
	return c.fallback
}
