package context

import (
	"fmt"
	"time"
)

func ExampleContext() {
	var c = Background()
	var child, cancel = WithCancel(WithValue(c, "K", "V"), time.Second*10)
	go func() {
		time.Sleep(time.Second * 5)
		cancel()
	}()
	<-child.Done()
	fmt.Println(child.Value("K"))
	// Output:
	// V
}
