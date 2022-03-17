package lexer 

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestWordToken(t *testing.T) {
	word := NewWordToken("variable", ID)
	require.Equal(t, "variable", word.ToString())
	word_tag := word.Tag 
	require.Equal(t, word_tag.ToString(), "ID")
}

func TestKeyWords(t *testing.T) {
	key_words := GetKeyWords()
	require.Equal(t, len(key_words), 10)

	and_key_word := key_words[0]
	require.Equal(t, and_key_word.ToString(), "&&")

	or_key_word := key_words[1]
	require.Equal(t, or_key_word.ToString(), "||")
}