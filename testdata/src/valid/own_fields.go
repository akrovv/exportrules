package valid

type Human struct {
	name string
	age  int
}

func (h Human) Name() string {
	return h.name
}

func (h Human) Age() {
	h.age++
}

func (h Human) UpdateAge() {
	h.age = 10
	h.age--
	h.age *= 2
	h.age += 2
	h.age++
	h.age = h.age + 2
}

func (h Human) SumName() {
	h.name = h.name + " Multiply"
}
