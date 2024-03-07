package printer

import "fmt"

type Printer interface {
	Println(a ...interface{})
}

type printer struct{}

func NewPrinter() Printer {
	return printer{}
}

func (p printer) Println(a ...interface{}) {
	fmt.Println(a...)
}
