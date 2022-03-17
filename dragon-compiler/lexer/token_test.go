package lexer 

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestTokenName(t *testing.T) {
    index_token := NewToken(INDEX)
	require.Equal(t, "INDEX", index_token.ToString())

	real_token := NewToken(REAL)
	require.Equal(t, "REAL", real_token.ToString())
}