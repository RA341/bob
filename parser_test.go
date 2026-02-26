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

func Test_parseCmd(t *testing.T) {
	var b Bobfile

	parseCmd(&b, CleanLine{})
}

func TestParse(t *testing.T) {

}
