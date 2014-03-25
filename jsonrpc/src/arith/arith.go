package arith

import (
	"errors"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith struct {
	LastCall string
}

func (t *Arith) Multiply(args *Args, reply *int) error {
	t.LastCall = "Multiply"
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	t.LastCall = "Divide"
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
