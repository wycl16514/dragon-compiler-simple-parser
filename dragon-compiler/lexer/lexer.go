package lexer

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	Lexeme    string
	peek      byte
	line      int
	reader    *bufio.Reader
	key_words map[string]Token
}

func NewLexer(source string) Lexer {
	str := strings.NewReader(source)
	source_reader := bufio.NewReaderSize(str, len(source))
	lexer := Lexer{
		line:      1,
		reader:    source_reader,
		key_words: make(map[string]Token),
	}

	lexer.reserve()

	return lexer
}

func (l *Lexer) ReverseScan() {
	back_len := len(l.Lexeme)
	for i := 0; i < back_len; i++ {
		l.reader.UnreadByte()
	}
}

func (l *Lexer) reserve() {
	key_words := GetKeyWords()
	for _, key_word := range key_words {
		l.key_words[key_word.ToString()] = key_word.Tag
	}
}

func (l *Lexer) Readch() error {
	char, err := l.reader.ReadByte() //提前读取下一个字符
	l.peek = char
	return err
}

func (l *Lexer) ReadCharacter(c byte) (bool, error) {
	chars, err := l.reader.Peek(1)
	if err != nil {
		return false, err
	}

	peekChar := chars[0]
	if peekChar != c {
		return false, nil
	}

	l.Readch() //越过当前peek的字符
	return true, nil
}

func (l *Lexer) UnRead() error {
	return l.reader.UnreadByte()
}

func (l *Lexer) Scan() (Token, error) {

	for {
		err := l.Readch()
		if err != nil {
			return NewToken(ERROR), err
		}

		if l.peek == ' ' || l.peek == '\t' {
			continue
		} else if l.peek == '\n' {
			l.line = l.line + 1
		} else {
			break
		}
	}

	l.Lexeme = ""

	switch l.peek {
	case '{':
		l.Lexeme = "{"
		return NewToken(LEFT_BRACE), nil
	case '}':
		l.Lexeme = "}"
		return NewToken(RIGHT_BRACE), nil
	case '+':
		l.Lexeme = "+"
		return NewToken(PLUS), nil
	case '-':
		l.Lexeme = "-"
		return NewToken(MINUS), nil
	case '(':
		l.Lexeme = "("
		return NewToken(LEFT_BRACKET), nil
	case ')':
		l.Lexeme = ")"
		return NewToken(RIGHT_BRACKET), nil
	case '&':
		l.Lexeme = "&"
		if ok, err := l.ReadCharacter('&'); ok {
			l.Lexeme = "&&"
			word := NewWordToken("&&", AND)
			return word.Tag, err
		} else {
			return NewToken(AND_OPERATOR), err
		}
	case '|':
		l.Lexeme = "|"
		if ok, err := l.ReadCharacter('|'); ok {
			l.Lexeme = "||"
			word := NewWordToken("||", OR)
			return word.Tag, err
		} else {
			return NewToken(OR_OPERATOR), err
		}

	case '=':
		l.Lexeme = "="
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "=="
			word := NewWordToken("==", EQ)
			return word.Tag, err
		} else {
			return NewToken(ASSIGN_OPERATOR), err
		}

	case '!':
		l.Lexeme = "!"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "!="
			word := NewWordToken("!=", NE)
			return word.Tag, err
		} else {
			return NewToken(NEGATE_OPERATOR), err
		}

	case '<':
		l.Lexeme = "<"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "<="
			word := NewWordToken("<=", LE)
			return word.Tag, err
		} else {
			return NewToken(LESS_OPERATOR), err
		}

	case '>':
		l.Lexeme = ">"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = ">="
			word := NewWordToken(">=", GE)
			return word.Tag, err
		} else {
			return NewToken(GREATER_OPERATOR), err
		}

	}

	if unicode.IsNumber(rune(l.peek)) {
		var v int
		var err error
		for {
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				l.UnRead() //将字符放回以便下次扫描
				break
			}
			v = 10*v + num
			l.Readch()

			l.Lexeme += string(l.peek)
		}

		if l.peek != '.' {
			return NewToken(NUM), err
		}
		l.Lexeme += string(l.peek)

		x := float64(v)
		d := float64(10)
		for {
			l.Readch()
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				l.UnRead()
				break
			}

			x = x + float64(num)/d
			d = d * 10
			l.Lexeme += string(l.peek)
		}

		return NewToken(REAL), err
	}

	if unicode.IsLetter(rune(l.peek)) {
		var buffer []byte
		for {
			buffer = append(buffer, l.peek)
			l.Lexeme += string(l.peek)

			l.Readch()
			if !unicode.IsLetter(rune(l.peek)) {
				l.UnRead()
				break
			}
		}

		s := string(buffer)
		token, ok := l.key_words[s]
		if ok {
			return token, nil
		}

		return NewToken(ID), nil
	}

	return NewToken(EOF), nil
}
