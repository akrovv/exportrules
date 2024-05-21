package tests

type Person struct{}

func New() Person {
	return Person{}
}

func (p *Person) foo() {}

func (p *Person) Getter() {
	p.foo()
}
