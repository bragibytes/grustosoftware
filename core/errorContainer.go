package core

import (
	"log"
	"time"
)

type errContainer struct {
	errs []err
}

func NewErrorContainer() *errContainer {
	x := &errContainer{
		make([]err, 0),
	}

	return x
}

func (x *errContainer) clear() {
	time.Sleep(10 * time.Second)

	x.errs = make([]err, 0)

}

func (x *errContainer) Errors() []err {
	return x.errs
}

func (x *errContainer) AddError(e error) {
	log.Printf("\n--- !ERROR! --- %v\n", e.Error())
	x.errs = append(x.errs, e.(err))
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
