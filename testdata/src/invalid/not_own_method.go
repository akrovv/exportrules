package invalid

type person struct{}

type human struct{}

func (h human) foo() {
	p := person{}

	p.foo()
}

func (p *person) foo() {}
