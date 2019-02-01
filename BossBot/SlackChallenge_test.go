package BossBot

import (
	"fmt"
	"testing"
)

type ParentType interface {
	A() string
	B() int
}

type Selector struct {
}

func (s *Selector) A() string {
	return "Selector"
}

func (s *Selector) B() int {
	return 1
}

type Modifier struct {
}

func (m *Modifier) A() string {
	return "Modifier"
}

func (m *Modifier) B() int {
	return 2
}

func PutIntoSlice(slice *[]ParentType, a ParentType) {
	*slice = append(*slice, a)
}

func TestServer(t *testing.T) {
	//Server()
	var slice []ParentType
	s1 := &Selector{}
	s2 := &Selector{}
	m1 := &Modifier{}

	PutIntoSlice(&slice, s1)
	PutIntoSlice(&slice, m1)
	PutIntoSlice(&slice, s2)

	for _, s := range slice {
		fmt.Printf("S : %s, D: %d\n", (s).A(), (s).B())
	}

}
