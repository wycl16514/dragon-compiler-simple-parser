package simple_parser

import (
	"errors"
	"lexer"
)

type SimpleParser struct {
	lexer lexer.Lexer
}

func NewSimpleParser(lexer lexer.Lexer) *SimpleParser {
	return &SimpleParser{
		lexer: lexer,
	}
}

func (s *SimpleParser) list() error {
	//根据读取的第一个字符决定选取哪个生产式
	token, err := s.lexer.Scan()
	if err != nil {
		return err
	}

	if token.Tag == lexer.LEFT_BRACKET {
		//选择 list -> ( list )
		s.list()
		token, err = s.lexer.Scan()
		if token.Tag != lexer.RIGHT_BRACKET {
			err := errors.New("Missinf of right bracket")
			return err
		}
	}

	if token.Tag == lexer.NUM {
		// list -> number
		err = s.number()
		if err != nil {
			return err
		}
	}

	token, err = s.lexer.Scan()
	if err != nil {
		return err
	}

	if token.Tag == lexer.PLUS || token.Tag == lexer.MINUS {
		s.list() // list -> list + list , list -> list - list
	} else {
		s.lexer.ReverseScan()
	}

	return err
}

func (s *SimpleParser) number() error {
	if len(s.lexer.Lexeme) > 1 {
		err := errors.New("Number only allow 0-9")
		return err
	}

	return nil
}

func (s *SimpleParser) Parse() error {
	return s.list()
}
