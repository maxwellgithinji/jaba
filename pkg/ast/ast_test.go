package ast

import (
	"testing"

	"github.com/maxwellgithinji/jaba/pkg/token"
)

func TestString(t *testing.T) {
	Program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "var1",
					},
					Value: "var1",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "var2",
					},
					Value: "var2",
				},
			},
		},
	}

	if Program.String() != "let var1 = var2;" {
		t.Errorf("Expected 'let var1 = var2;' got '%s'", Program.String())
	}
}
