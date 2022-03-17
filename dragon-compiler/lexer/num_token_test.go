package lexer  
import (
	"testing"
	"github.com/stretchr/testify/require"
	
)

func TestNumToken(t *testing.T) {
	num_token := NewNumToken(123)
	require.Equal(t, num_token.value, 123)
	require.Equal(t, num_token.ToString(), "123")
	num_tag := num_token.Tag 
	require.Equal(t, num_tag.ToString(), "NUM")
}

func TestRealToken(t *testing.T) {
	real_token := NewRealToken(3.1415926)
	require.Equal(t, real_token.value, 3.1415926)
	require.Equal(t, real_token.ToString(), "3.1415926")

    real_tag := real_token.Tag 
	require.Equal(t, real_tag.ToString(), "REAL")
}