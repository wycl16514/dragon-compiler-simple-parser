package lexer

type Word struct {
	lexeme string
	Tag    Token
}

func NewWordToken(s string, tag Tag) Word {
	return Word{
		lexeme: s,
		Tag:    NewToken(tag),
	}
}

func (w *Word) ToString() string {
	return w.lexeme
}

func GetKeyWords() []Word {
	key_words := []Word{}
	key_words = append(key_words, NewWordToken("&&", AND))
	key_words = append(key_words, NewWordToken("||", OR))
	key_words = append(key_words, NewWordToken("==", EQ))
	key_words = append(key_words, NewWordToken("!=", NE))
	key_words = append(key_words, NewWordToken("<=", LE))
	key_words = append(key_words, NewWordToken(">=", GE))
	key_words = append(key_words, NewWordToken("minus", MINUS))
	key_words = append(key_words, NewWordToken("true", TRUE))
	key_words = append(key_words, NewWordToken("false", FALSE))
	key_words = append(key_words, NewWordToken("t", TEMP))
	key_words = append(key_words, NewWordToken("if", IF))
	key_words = append(key_words, NewWordToken("else", ELSE))

	return key_words
}
