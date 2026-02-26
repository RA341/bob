package main

import "testing"

var con = `version: 3

as := aasdasd

sd := 1

hello() {
    print("hello world")
}

hello2(
user: str!, 
otherP: sd = asdasd,
) {
    workdir core
    print("hello ${user}")
    print(as)
}
`

var sd = `bd(out=.build) {
    clear
    go build -o ${out}/dev
    // ./dev hello2 user=lmao
}
`

func Test_parseCmd(t *testing.T) {
	var b Bobfile

	ParseFromBytes(&b, []byte(sd))
}

func TestParse(t *testing.T) {

}
