package main

import (
	"regexp"
	"testing"
)

func TestFilterTags(t *testing.T) {
	tagsLineRe, err := regexp.Compile(`tags\s*=\s*\[[^\]]*\]`)
	if err != nil {
		t.Errorf("regexp is invalid: %s", err)
	}
	linesMatching := []string{
		`tags=["ala", "ola"]`,
		`tags = ["m","asd"]`}
	for _, l := range linesMatching {
		if !tagsLineRe.MatchString(l) {
			t.Errorf("%s expected to match %s", l, tagsLineRe.String())
		}
	}

}

type test struct {
	input string
	exp   string
}

type testArr struct {
	input string
	exp   []string
}

func TestGetGroupOfTags(t *testing.T) {
	tagsGroupRe, err := regexp.Compile(`tags\s*=\s*\[([^\]]*)\]`)
	if err != nil {
		t.Errorf("regexp is invalid: %s", err)
	}

	cases := []test{
		{`tags=["ala", "ola"]`, `"ala", "ola"`},
		{`tags = ["m","asd"]`, `"m","asd"`},
	}

	for _, tc := range cases {
		m := tagsGroupRe.FindStringSubmatch(tc.input)
		if len(m) <= 1 || m[1] != tc.exp {
			t.Errorf("%s expected to extract %s using re %s", tc.input, tc.exp, tagsGroupRe.String())
		}
	}
}

func TestGetTagsFromGroup(t *testing.T) {
	var tagsExtractRe, err = regexp.Compile(`"([^"]*)"`)
	if err != nil {
		t.Errorf("regexp is invalid: %s", err)
	}

	cases := []testArr{
		{`"ala", "ola"`, []string{"ala", "ola"}},
		{`"m","asd"`, []string{"m", "asd"}},
	}

	for _, tc := range cases {
		subm := tagsExtractRe.FindAllStringSubmatch(tc.input, -1)
		if subm == nil || len(subm) < 1 {
			for i := 0; i < len(tc.exp); i++ {
				if subm[i][1] != tc.exp[i] {
					t.Errorf("%s expected to be %s but was %s (with re %s)", tc.input, tc.exp, subm, tagsExtractRe.String())
				}
			}
		}
	}
}

func TestTagUpdate(t *testing.T) {
	to := NewTags()

	if to.nameToTag == nil {
		t.Error("to.nameToTag should not be nil")
	}
	to.Update("Ewa")
	to.Update("Kasia")
	to.Update("ala")
	to.Update("ala")
	to.Update("ala")
	to.Update("ola")
	to.Update("ola")
	to.Update("ula")
	to.Update("Xula")
	to.Update("Xula")
	to.Update("Xula")
	to.Update("Xula")

	type tcase struct {
		tname    string
		expcount int
	}
	tcases := []tcase{
		{"ala", 3},
		{"ola", 2},
		{"ula", 1},
		{"Xula", 4},
	}

	t.Run("Check map contents", func(t *testing.T) {
		for _, tc := range tcases {
			tg, ok := to.nameToTag[tc.tname]
			if !ok {
				t.Errorf("Tags should have key %s", tc.tname)
			}
			if tg.name != tc.tname {
				t.Errorf("expected Tag name to be %s, was %s", tc.tname, tg.name)
			}
			if tg.count != tc.expcount {
				t.Errorf("expected Tag count to be %d, was %d", tc.expcount, tg.count)
			}
		}
	})

	t.Run("Check sorting", func(t *testing.T) {

		arr := to.Sorted()
		expected := TagA{Tag{"Xula", 4}, Tag{"ala", 3}, Tag{"ola", 2}, Tag{"Ewa", 1}, Tag{"Kasia", 1}, Tag{"ula", 1}}

		for i := 0; i < len(expected); i++ {
			if arr[i] != expected[i] {
				t.Errorf("Expected position at %d contain %+v, got %+v", i, expected[i], arr[i])
			}
		}
	})

}
