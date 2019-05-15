package core

import (
	"log"
	"time"
)

type errContainer struct {
	errs     []error
	clearing bool
}

func NewErrorContainer() *errContainer {
	x := &errContainer{
		make([]error, 0),
		false,
	}

	return x
}

func (x *errContainer) clear() {
	time.Sleep(5 * time.Second)

	x.errs = make([]error, 0)
	x.clearing = true

}

func (x *errContainer) ListErrors() {
	for _, v := range x.errs {
		log.Printf("\n\n\n %v \n\n\n", v.Error())
	}
}

func (x *errContainer) Errors() []error {
	return x.errs
}

func (x *errContainer) printErr(es string) {
	log.Print("\n")
	log.Printf("\n--- !ERROR! --- %v\n", es)
	log.Print("\n")
}

func (x *errContainer) AddError(e error) {

	x.printErr(e.Error())
	x.errs = append(x.errs, e)
	if !x.clearing {
		go x.clear()
	}
}

//func (x *errContainer) AddErrors(e []error) {
//
//	log.Printf("--- adding errors --- %v\n", e)
//	for _, v := range e {
//		x.AddError(v)
//	}
//	go x.clear()
//}

func (x *errContainer) ErrorCount() int {
	return len(x.errs)
}
