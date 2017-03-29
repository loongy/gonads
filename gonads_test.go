package gonads

import "testing"

func TestDoThenFinally(t *testing.T) {
	wait := make(chan struct{})

	do := false
	then1 := false
	then2 := false

	Do(func() error {
		if do || then1 || then2 {
			t.Fatal("Do must execute 1st")
		}
		do = true
		return nil
	}).Then(func() error {
		if !do || then1 || then2 {
			t.Fatal("Then #1 must execute 2nd")
		}
		then1 = true
		return nil
	}).ThenAll(func() error {
		if !do || !then1 || then2 {
			t.Fatal("ThenAll #1 must execute 3rd, 4th, or 5th")
		}
		return nil
	}, func() error {
		if !do || !then1 || then2 {
			t.Fatal("ThenAll #2 must execute 3rd, 4th, or 5th")
		}
		return nil
	}, func() error {
		if !do || !then1 || then2 {
			t.Fatal("ThenAll #3 must execute 3rd, 4th, or 5th")
		}
		return nil
	}).Then(func() error {
		if !do || !then1 || then2 {
			t.Fatal("Then #2 must execute 6th")
		}
		then2 = true
		return nil
	}).Finally(func() {
		if !do || !then1 || !then2 {
			t.Fatal("Finally must execute 7th")
		}
		wait <- struct{}{}
	})

	<-wait
}
