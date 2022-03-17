package lexer 

import (
	"strconv"
	"fmt"
)

type Num struct {
	Tag Token
	value int 
}

func NewNumToken(val int) Num {
	return Num {
		Tag: NewToken(NUM),
		value: val,
	}
}

func (n *Num) ToString() string {
	return  strconv.Itoa(n.value)
}

type Real struct {
	Tag Token 
	value float64 
}

func NewRealToken(val float64) Real{
	return Real {
		value: val,
		Tag: NewToken(REAL),
	}
}

func (r *Real) ToString() (string) {
    return fmt.Sprintf("%.7f", r.value)
}