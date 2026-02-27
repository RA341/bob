package parser

import (
	"testing"

	"github.com/RA341/bob/util"
)

func TestParser_Hello_world(t *testing.T) {
	bobFile := `hello() {
    	print("hello world")
	}`
	var bd Bobfile
	ParseFromBytes(&bd, []byte(bobFile))
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
