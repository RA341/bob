package flag

import (
	"fmt"
	"strconv"
	"strings"
)

// Flag any val starting with --
type Flag struct {
	Name  string
	Val   any
	IsSet bool
	Help  string
}

type Parser struct {
	Flags []*Flag

	Current int
	len     int
}

func ParseFlags(flags []*Flag, args []string) (Parser, error) {
	fp := Parser{
		Flags: flags,
	}

	err := fp.Parse(args)
	return fp, err
}

func (f *Parser) Parse(in []string) error {
	peek := func() string {
		if f.Current >= len(in) {
			return ""
		}
		return in[f.Current]
	}

	isAtEnd := func() bool {
		if f.Current >= len(in) {
			return true
		}

		if !strings.HasPrefix(peek(), "--") {
			return true
		}

		return false
	}

	next := func() string {
		if isAtEnd() {
			return ""
		}

		va := peek()
		f.Current++
		return va
	}

	loadArg := func() (string, error) {
		nx := peek()

		if strings.HasPrefix(nx, "--") {
			return "", fmt.Errorf("expected value got flag: %s", nx)
		}

		if nx == "" {
			return "", fmt.Errorf("expected value got no value")
		}

		f.Current++
		return nx, nil
	}

	for !isAtEnd() {
		nex := next()
		val, ok := strings.CutPrefix(nex, "--")

		for _, flag := range f.Flags {
			if ok && val == flag.Name {
				switch v := flag.Val.(type) {
				case bool:
					flag.Val = true
				case string:
					load, err := loadArg()
					if err != nil {
						return err
					}

					flag.Val = load
				case int:
					loaded, err := loadArg()
					if err != nil {
						return err
					}

					atoi, err := strconv.Atoi(loaded)
					if err != nil {
						return err
					}

					flag.Val = atoi

				default:
					return fmt.Errorf("unsupported type: %q for flag %q", v, flag.Name)
				}
				flag.IsSet = true
			}

		}
	}

	return nil
}
