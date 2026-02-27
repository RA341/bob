package parser

import (
	"testing"

	"github.com/RA341/bob/util"
	"github.com/RA341/bob/vm"
	"github.com/stretchr/testify/require"
)

func TestParser_global_var(t *testing.T) {
	bobFile := `
	@someVar = value
	@someVar2 = value
	@someInt: VTInt = 0
	@someExpr = (@someVar2+@someVar)
`
	var bd Bobfile
	err := NewBobFileFromBytes(&bd, []byte(bobFile))
	require.NoError(t, err)

	vmm := new(vm.VM)
	get := bd.Program.Get()
	bd.Program.Print()

	vmm.Start(get, nil)

	val, ok := vmm.Vars["someVar"]
	require.True(t, ok)
	require.Equal(t, vm.VTString, val.Type)
	require.Equal(t, "value", val.Raw)

	val, ok = vmm.Vars["someVar2"]
	require.True(t, ok)
	require.Equal(t, vm.VTString, val.Type)
	require.Equal(t, "value", val.Raw)

	val, ok = vmm.Vars["someInt"]
	require.True(t, ok)
	require.Equal(t, vm.VTInt, val.Type)
	require.Equal(t, "0", val.Raw)

	val, ok = vmm.Vars["someExpr"]
	require.True(t, ok)
	require.Equal(t, vm.VTString, val.Type)
	require.Equal(t, "valuevalue", val.Raw)
}

func TestParser_Hello_world(t *testing.T) {
	bobFile := `hello() {
    	print("hello world")
	}

	hello2(){
		@someVar = "world"
		hello()
	}
`
	var bd Bobfile
	err := NewBobFileFromBytes(&bd, []byte(bobFile))
	require.NoError(t, err)
}

var con = `
@as=aasdasd
@sd=1

hello2(
	user!, 
	other2:,
) {
    workdir core
    print("hello ${user}")
    print(as)
}
`

var sd = `
dk:rn() {
	@img = test
	@tag = wow	

    bob dk 
    docker run --rm -p 3002:3000 ${img}:${tag}
}
`

func Test_parseCmd(t *testing.T) {
	var b Bobfile
	util.UNUSED(b)
	//ParseFromBytes(&b, []byte(sd))
}

func TestParse(t *testing.T) {

}
