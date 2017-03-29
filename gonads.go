package gonads

// A Gonad can be used to chain set of concurrent operations. Each step in the
// chain contains a set of operations. If any error occurs the chain is broken
// and skips to the next error handler.
type Gonad struct {
	done chan struct{}
	err  chan error
}

// Do an operation concurrently.
func Do(f func() error) *Gonad {
	g := &Gonad{
		done: make(chan struct{}),
		err:  make(chan error),
	}

	go func() {
		err := f()
		g.done <- struct{}{}
		g.err <- err
	}()

	return g
}

// DoAll operations concurrently.
func DoAll(fs ...func() error) *Gonad {
	g := &Gonad{
		done: make(chan struct{}),
		err:  make(chan error),
	}

	dones := make(chan struct{}, len(fs))
	errs := make(chan error, len(fs))
	for _, f := range fs {
		go func(f func() error) {
			err := f()
			dones <- struct{}{}
			errs <- err
		}(f)
	}

	go func() {
		var err error
		for i := 0; i < len(fs); i++ {
			<-dones
		}
		for i := 0; err == nil && i < len(fs); i++ {
			err = <-errs
		}

		g.done <- struct{}{}
		g.err <- err
	}()

	return g
}

// Then do an operation concurrently, after all operations in the Gonad
// complete without error.
func (g *Gonad) Then(f func() error) *Gonad {
	nextG := &Gonad{
		done: make(chan struct{}),
		err:  make(chan error),
	}

	go func() {
		<-g.done
		err := <-g.err

		if err == nil {
			err = f()
		}

		nextG.done <- struct{}{}
		nextG.err <- err
	}()

	return nextG
}

// ThenAll operations are done concurrently, after all operations in the
// Gonad complete without error.
func (g *Gonad) ThenAll(fs ...func() error) *Gonad {
	nextG := &Gonad{
		done: make(chan struct{}),
		err:  make(chan error),
	}

	go func() {
		<-g.done
		err := <-g.err

		if err != nil {
			nextG.done <- struct{}{}
			nextG.err <- err
			return
		}

		dones := make(chan struct{}, len(fs))
		errs := make(chan error, len(fs))
		for _, f := range fs {
			go func(f func() error) {
				err := f()
				dones <- struct{}{}
				errs <- err
			}(f)
		}

		for i := 0; i < len(fs); i++ {
			<-dones
		}
		for i := 0; err == nil && i < len(fs); i++ {
			err = <-errs
		}

		nextG.done <- struct{}{}
		nextG.err <- err
	}()

	return nextG
}

// Else do an error handling operation concurrently, after any operation in
// the Gonad completes with an error.
func (g *Gonad) Else(f func(err error)) *Gonad {
	nextG := &Gonad{
		done: make(chan struct{}),
		err:  make(chan error),
	}

	go func() {
		<-g.done
		err := <-g.err

		if err != nil {
			f(err)
		}

		nextG.done <- struct{}{}
		nextG.err <- err
	}()

	return nextG
}

// Finally do an operation concurrently, after all operations in the Gonad
// have completed, regardless of any errors.
func (g *Gonad) Finally(f func()) {
	go func() {
		<-g.done
		<-g.err
		f()
	}()
}
