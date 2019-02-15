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
