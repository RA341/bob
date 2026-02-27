package vm

import "fmt"

//type BuiltInFunc func(env *RuntimeEnv, vals ...string) error
//
//var builtInMaps = map[string]BuiltInFunc{
//	"print": func(env *RuntimeEnv, vals ...string) error {
//		vals = append([]string{green.Sprint("=>")}, vals...)
//
//		var interfaceArgs []interface{}
//		for _, val := range vals {
//			interfaceArgs = append(interfaceArgs, val)
//		}
//
//		fmt.Println(interfaceArgs...)
//		return nil
//	},
//	"workdir": func(env *RuntimeEnv, vals ...string) error {
//		if len(vals) < 1 {
//			env.workingDir = env.originalWorkingDir
//			return nil
//		}
//
//		env.workingDir = vals[0]
//		return nil
//	},
//}

var DefaultFns = map[string]FnDef{
	"print": NewFnDev(VariadicArgCount, func(args ...string) error {
		var interfaceArgs []interface{}
		for _, arg := range args {
			interfaceArgs = append(interfaceArgs, arg)
		}
		fmt.Println(interfaceArgs...)

		return nil
	}),
	"printf": NewFnDev(VariadicArgCount, func(args ...string) error {
		if len(args) < 1 {
			return fmt.Errorf("in sufficient args to call printf")
		}

		var interfaceArgs []interface{}
		for _, arg := range args[1:] {
			interfaceArgs = append(interfaceArgs, arg)
		}

		fmt.Printf(args[0], interfaceArgs...)
		return nil
	}),
}
