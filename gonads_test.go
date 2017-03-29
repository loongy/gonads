package gonads

import "testing"
import "log"

func TestDoThenFinally(t *testing.T) {
	wait := make(chan struct{})

	Do(func() error {
		log.Println("Do")
		return nil
	}).Then(func() error {
		log.Println("Then #1")
		return nil
	}).ThenAll(func() error {
		log.Println("ThenAll #1")
		return nil
	}, func() error {
		log.Println("ThenAll #2")
		return nil
	}, func() error {
		log.Println("ThenAll #3")
		return nil
	}, func() error {
		log.Println("ThenAll #4")
		return nil
	}).Then(func() error {
		log.Println("Then #2")
		return nil
	}).Finally(func() {
		log.Println("Finally")
		wait <- struct{}{}
	})

	<-wait
	t.Parallel()
}
