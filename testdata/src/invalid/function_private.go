package invalid

import "fmt"

type Alien struct{}

func (a Alien) call() {
	fmt.Println("Alien was called")
}

func calling() {
	a := Alien{}
	a.call()
}
