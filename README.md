# Monads for Go

Gonads allow you to chain sets of concurrent operations.

```go
import "errors"
import "github.com/mushu/gonads"
import "log"

func main() {

  gonads.Do(func() error {
    log.Println("This happens 1st")
    return nil
  }).Then(func() error {
    log.Println("This happens 2nd")
    return nil
  }).ThenAll(func() error {
    log.Println("This happens 3rd, 4th, or 5th")
    return nil
  }, func() error {
    log.Println("This happens 3rd, 4th, or 5th")
    return nil
  }, func() error {
    log.Println("This happens 3rd, 4th, or 5th")
    return nil
  }).Then(func() error {
    log.Println("This happens 6th")
    return errors.New("This is an error")
  }).Then(func() error {
    log.Println("This does not happen")
    return nil
  }).Else(func(err error) error {
    log.Println(err)
    return err
  }).Finally(func() {
    log.Println("This always happens last")
  })

}
```
