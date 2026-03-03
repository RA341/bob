package parser

import (
	"errors"
	"os"

	"github.com/RA341/bob/vm"
)

func ParseBobFromPath(fPath string) ([]vm.Ins, error) {
	cont, err := os.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	return ParseBobFromContents(cont)
}

func ParseBobFromContents(cont []byte) ([]vm.Ins, error) {
	lex := RunLexer(cont)
	if lex.errs != nil {
		return nil, errors.Join(lex.errs...)
	}

	parsed := RunParser(lex.tokens)
	if parsed.errs != nil {
		return nil, errors.Join(lex.errs...)
	}

	return parsed.instructions, nil
}
