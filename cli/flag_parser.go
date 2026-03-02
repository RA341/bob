package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// Flag any val starting with --
type Flag struct {
	name  string
	val   any
	isSet bool

	help string
}

type FlagParser struct {
	flags []*Flag

	current int
	len     int
}

func IntFlag(name, help string) *Flag {
	// use 2 to indicate int otherwise it is considered bool
	return NewFlag(name, help, 2)
}

func StrFlag(name, help string) *Flag {
	return NewFlag(name, help, "")
}

func BoolFlag(name, help string) *Flag {
	return NewFlag(name, help, false)
}

func NewFlag(name string, help string, val any) *Flag {
	return &Flag{
		name: name,
		val:  val,
		help: help,
	}
}

func ParseFlags(flags []*Flag, args []string) (FlagParser, error) {
	fp := FlagParser{
		flags: flags,
	}

	err := fp.parse(args)
	return fp, err
}

func (f *FlagParser) parse(in []string) error {
	peek := func() string {
		if f.current >= len(in) {
			return ""
		}
		return in[f.current]
	}

	isAtEnd := func() bool {
		if f.current >= len(in) {
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
		f.current++
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

		f.current++

		return nx, nil
	}

	for !isAtEnd() {
		nex := next()
		val, ok := strings.CutPrefix(nex, "--")

		for _, flag := range f.flags {
			if ok && val == flag.name {
				switch v := flag.val.(type) {
				case bool:
					flag.val = true
				case string:
					load, err := loadArg()
					if err != nil {
						return err
					}

					flag.val = load
				case int:
					loaded, err := loadArg()
					if err != nil {
						return err
					}

					atoi, err := strconv.Atoi(loaded)
					if err != nil {
						return err
					}

					flag.val = atoi

				default:
					return fmt.Errorf("unsupported type: %q for flag %q", v, flag.name)
				}
				flag.isSet = true
			}

		}
	}

	return nil
}
