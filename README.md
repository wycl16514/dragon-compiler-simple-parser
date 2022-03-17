语法解析本质上是判断给定的字符串序列是否符合特定规则，它是编译原理中难度相当大的部分，当然也相当不好理解。举个例子，我们如何识别由数字0到9，和符号(,),+,-所形成的算术表达式，例如"1+2", "1+(3-2)", "1", "((1+2)+(((4+4))))"都是满足规则的表达式，然而"()+1"就不能满足。一种直观的做法是我们依次读入字符然后做即时判断，例如首先看第一个读到的字符是不是数字，或者是不是左括号，然后根据读入的前一个字符看看接下来读入的字符是否合法，你可以尝试用代码来实现试试，你很快会发现代码非常难写。

有没有系统化的方法来处理这样的问题呢。编译原理中的语法解析就是解决这类问题的方案。我们看看如何解决上面提到的问题，在编译原理中有一种数据结构胶backus-nour范式，它给出了一种自动化的判断给定字符串是否复合特定规则的方法，例如上面的问题对应的backus范式为：
```
list -> "(" list ")"
list -> list  "+" list 
list -> list "-" list 
list -> number
number -> "0" | "1" | "2" | "3" | "4"| "5" | "6" | "7" | "8" | "9"
```
如果你是第一次接触这个东西，你会感觉很难理解，其中一个原因在于，它使用递归的方法来定义字节，上面带有->的表达式我们称为生产式，出现在箭头左边的符号叫非终结符，只出现在右边的符号叫终结符，例如字符"0" 到 "9"，和左右括号"(",")"。所谓非终结符就是能通过箭头右边的符号进行分解，这里一个难点在于它可以自己分解自己，例如list -> ( list ) 中，左边的list 可以分解成左括号， 然后是list 和 右括号的组合，我们先看具体例子，假设给定表达式3 + 2，我们怎么用上面的生产式来判断它是否符合规定呢。

算法的基本做法是选择相应的生产式进行”套用“，直到生产式解析为终结符为止。于是对应表达式3+2，我们可以猜到可以使用list -> list + list, 因为只有它含有符号"+"。于是我们接下来的任务就是看 3 和 2是否满足list的定义，此时我们不难猜测可以使用list -> number，于是我们又得判断3, 2是否能使用number来解析，现在我们看到number右边的字符包含0到9，于是可以解析，由此表达式3 + 2满足上面生产式所规定的规则。

生产式是对字符串组合规律的一种抽象描述，所有能满足给定生产式的字符串组合就叫做生产式生成的“语言”。给定一系列字符串的组合，然后判断其是否满足给定生产式的判断过程叫“推导”，同时生产式所描述的规则就叫做"语法“。我们再看一个例子，java,c++,c代码中函数调用,例如max(x,y), 其的语法:
```
call -> ID ( optparams )
optparams -> params | "ε"
params -> params  "," param | param
(此处 param 的生产式没有给出来) 
```
语法的定义比较抽象，通过这里几个例子，大家有没有一些感性认识。我个人觉得很难用语言来描述什么叫语法，但我发现如果使用代码的话，或许能让人有“心领神会”的感觉。另外值得一提的是推导的基本逻辑，我们看到推导实际上是用生产式去”套用“字符串，看看能不能一路解析到终结符，但是生产式有若干个，我们如何确定用哪个去套呢？当我们给定(3+2)时，我们不难猜到用list -> "(" list ")" 去套，那是因为我们看到表达式的第一个字符跟字符串的第一个字符一样，这意味着在推导过程中，我们通过读取第一个字符来选择合适的生产式来进行推导。

由于语法解析不好用语言说明，有意思的是使用代码反而能描述更清楚，因此我们直接通过代码来理解如何使用生产式来匹配给定的字符串，首先我们对上次完成的lexer做一些修改：
```
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

```

上面代码的修改主要是增加了lexeme，用来记录当前读到的字符串，增加了函数UnRead，用来把当前读到的字符重新放回缓冲器，接下来我们看看解析器的实现，在上一节代码的目录parser,然后在里面增加文件simple_parser.go,并完成如下代码：
```
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

```
在上面代码实现中有几点需要注意，在list函数中，它完全模拟了对应的生产式，例如在使用list -> "(" list ")"，时，代码首先判断读入的是不是左括号，然后递归的调用list函数字节，最后再判断最终读入的是不是右括号，所以生产式本质上是指导我们代码如何实现，代码对读入字符的判断，以及自我递归等步骤都完全根据生产式来进行。最后我们看看如何调用如上代码，在main.go中输入代码如下：
```
package main

import (
	"fmt"
	"io"
	"lexer"
	"simple_parser"
)

func main() {
	source := "(1+(2+3))"
	my_lexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(my_lexer)
	err := parser.Parse()
	if err != nil && err != io.EOF {
		fmt.Println("source error: ", err)
	} else {
		fmt.Println("source is legal expression")
	}
}

```
有兴趣的同学可以运行起代码看看，以上就是语法解析的基本原理。语法解析在编译原理中是非常复杂的一个模块，这里我们通过实践的方式提前了解到其一些基本概念和原理，这对我们未来更好的深入理解其原理打下扎实基础。
