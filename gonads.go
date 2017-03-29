package gonads

// A Gonad can be used to chain set of concurrent operations. Each step in the
// chain contains a set of operations. If any error occurs the chain is broken
// and skips to the next error handler.
type Gonad struct {
	value chan interface{}
	err   chan error
}

// Do an operation concurrently.
func Do(f func() (interface{}, error)) Gonad {
	g := Gonad{
		value: make(chan interface{}),
		err:   make(chan error),
	}

	go func() {
		value, err := f()
		g.value <- value
		g.err <- err
	}()

	return g
}

// DoAll operations concurrently.
func DoAll(fs ...func() (interface{}, error)) Gonad {
	g := Gonad{
		value: make(chan interface{}),
		err:   make(chan error),
	}

	values := make(chan interface{}, len(fs))
	errs := make(chan error, len(fs))
	for _, f := range fs {
		go func(f func() (interface{}, error)) {
			value, err := f()
			values <- value
			errs <- err
		}(f)
	}

	go func() {
		value := make([]interface{}, len(fs))
		var err error
		for i := 0; i < len(fs); i++ {
			value[i] = <-values
		}
		for err == nil {
			err = <-errs
		}
		close(values)
		close(errs)

		g.value <- value
		g.err <- err
	}()

	return g
}

// Then do an operation concurrently, after all operations in the Gonad
// complete without error.
func (g Gonad) Then(f func(interface{}) (interface{}, error)) Gonad {
	nextG := Gonad{
		value: make(chan interface{}),
		err:   make(chan error),
	}

	go func() {
		value := <-g.value
		err := <-g.err
		close(g.value)
		close(g.err)

		if err == nil {
			value, err = f(value)
		}

		nextG.value <- value
		nextG.err <- err
	}()

	return nextG
}

// ThenAll operations are done concurrently, after all operations in the
// Gonad complete without error.
func (g Gonad) ThenAll(fs ...func(interface{}) (interface{}, error)) Gonad {
	nextG := Gonad{
		value: make(chan interface{}),
		err:   make(chan error),
	}

	go func() {
		value := <-g.value
		err := <-g.err
		close(g.value)
		close(g.err)

		if err != nil {
			nextG.value <- value
			nextG.err <- err
			return
		}

		values := make(chan interface{}, len(fs))
		errs := make(chan error, len(fs))
		for _, f := range fs {
			go func(f func(interface{}) (interface{}, error)) {
				value, err := f(value)
				values <- value
				errs <- err
			}(f)
		}

		nextValue := make([]interface{}, len(fs))
		var nextErr error
		for i := 0; i < len(fs); i++ {
			nextValue[i] = <-values
		}
		for nextErr == nil {
			nextErr = <-errs
		}
		close(values)
		close(errs)

		nextG.value <- nextValue
		nextG.err <- nextErr
	}()

	return nextG
}

// Else do an error handling operation concurrently, after any operation in
// the Gonad completes with an error.
func (g Gonad) Else(f func(err error)) Gonad {
	nextG := Gonad{
		value: make(chan interface{}),
		err:   make(chan error),
	}

	go func() {
		value := <-g.value
		err := <-g.err
		close(g.value)
		close(g.err)

		if err != nil {
			f(err)
		}

		nextG.value <- value
		nextG.err <- err
	}()

	return nextG
}

// Finally do an operation concurrently, after all operations in the Gonad
// have completed, regardless of any errors.
func (g Gonad) Finally(f func()) {
	go func() {
		<-g.value
		<-g.err
		close(g.value)
		close(g.err)

		f()
	}()
}
