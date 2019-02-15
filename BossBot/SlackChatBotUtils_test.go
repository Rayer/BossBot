package BossBot

import (
	"regexp"
	"testing"
)

func TestRegexCompile(t *testing.T) {
	r, err := regexp.Compile(`\[([A-Za-z 0-9_]*)]`)
	if err != nil {
		t.Fatalf("Error : %+v", err)
	}

	t.Log(r.FindAllString("[This is] a [book] about [Science and Space], and will you [read] it?", -1))

}

type OuterInterface interface {
	interface1method()
}

type OuterInterfaceImpl struct {
	num  int
	char string
}

func NewOuterInterfaceImpl(num int, char string) *OuterInterfaceImpl {
	return &OuterInterfaceImpl{num: num, char: char}
}

type testbed struct {
	OuterInterfaceImpl
}

func newTestbed() *testbed {
	return &testbed{OuterInterfaceImpl: *NewOuterInterfaceImpl(1, "ccc")}
}

func TestEmbedded(t *testing.T) {
	tb := newTestbed()
	t.Log(*tb)
	t.Logf("%+v", *tb)
}
