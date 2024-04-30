package helpers

func Ptr[T any](x T) *T {
	return &x
}

func Deref[T any](x *T) T {
	if x == nil {
		x = new(T)
	}

	return *x
}

func Nil[T any](x T) *T {
	return nil
}

func Any[T any](x T) any {
	return x
}
