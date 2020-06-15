package lib

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Tests テスト仕様
type Tests struct {
	Name   string
	path   string
	Genres []Genre
}

// Genre ジャンル
type Genre struct {
	Num        int
	Genre      string
	Categories []Category
	// for width
	maxCategoryLen int
}

// Category カテゴリ
type Category struct {
	Num      int
	Category string
	Cases    []Case
	// for width
	maxCaseLen  int
	maxStepLen  int
	maxCheckLen int
}

// Case テスト手順
type Case struct {
	Num    int
	Case   string
	Steps  []string
	Checks []string
	// for width
	stepLen  int
	checkLen int
}

func NewTests(name string) *Tests {
	return &Tests{
		Name:   name,
		path:   "",
		Genres: make([]Genre, 0),
	}
}

func NewGenre(name string, te *Tests) *Genre {
	g := Genre{
		Num:        len(te.Genres),
		Genre:      name,
		Categories: make([]Category, 0),
	}
	te.Genres = append(te.Genres, g)
	return &g
}

func NewCategory(name string, ge *Genre) *Category {
	c := Category{
		Num:      len(ge.Categories),
		Category: name,
		Cases:    make([]Case, 0),
	}
	ge.Categories = append(ge.Categories, c)
	return &c
}

func NewCase(name string, ca *Category) *Case {
	c := Case{
		Num:    len(ca.Cases),
		Case:   name,
		Steps:  make([]string, 0),
		Checks: make([]string, 0),
	}
	ca.Cases = append(ca.Cases, c)
	return &c
}

func (t *Tests) LastGenre() *Genre {
	return &t.Genres[len(t.Genres)-1]
}

func (t *Tests) SetPath(path string) error {
	if path != "" {
		t.path = path
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		t.path = wd
	}
	return nil
}

func (g *Genre) LastCategory() *Category {
	return &g.Categories[len(g.Categories)-1]
}

func (g *Genre) GetMaxCategory() int {
	max := 0
	for _, c := range g.Categories {
		l := getStrLen(c.Category)
		if l > max {
			max = l
		}
	}
	return max
}

func (c *Category) LastCase() *Case {
	return &c.Cases[len(c.Cases)-1]
}

func (c *Category) GetMaxCase() int {
	max := 0
	for _, s := range c.Cases {
		l := getStrLen(s.Case)
		if l > max {
			max = l
		}
	}
	return max
}

func (c *Case) AddStep(step string, isNew bool) {
	if isNew {
		c.Steps = append(c.Steps, step)
	} else {
		c.Steps[len(c.Steps)-1] = fmt.Sprintf("%v\n  %v", c.Steps[len(c.Steps)-1], step)
	}
	l := getStrLen(step)
	// for width
	if l > c.stepLen {
		c.stepLen = l
	}
}

func (c *Case) GetStepCheckHeight() []int {
	stepNum := 0
	for _, s := range c.Steps {
		for _ = range strings.Split(s, "\n") {
			stepNum++
		}
	}

	checkNum := 0
	checks := make([]int, 0)
	for _, c := range c.Checks {
		cs := strings.Split(c, "\n")
		checks = append(checks, len(cs))
		for range cs {
			checkNum++
		}
	}

	max := 0
	if stepNum > checkNum {
		max = stepNum
	} else {
		max = checkNum
	}

	heights := make([]int, 0)
	for _, c := range checks {
		heights = append(heights, max*c/checkNum)
	}
	return heights
}

func (c *Case) AddCheck(check string, isNew bool) {
	if isNew {
		c.Checks = append(c.Checks, check)
	} else {
		c.Checks[len(c.Checks)-1] = fmt.Sprintf("%v\n%v", c.Checks[len(c.Checks)-1], check)
	}
	// for width
	l := getStrLen(check)
	if l > c.checkLen {
		c.checkLen = l
	}
}

func (c *Case) GetMaxStep() int {
	max := 0
	for _, s := range c.Steps {
		for _, ss := range strings.Split(s, "\n") {
			l := getStrLen(ss)
			if l > max {
				max = l
			}
		}
	}
	return max
}

func (c *Case) GetMaxCheck() int {
	max := 0
	for _, s := range c.Checks {
		l := getStrLen(s)
		if l > max {
			max = l
		}
	}
	return max
}

func (t *Tests) ChecksNum() int {
	num := 0
	for _, g := range t.Genres {
		for _, c := range g.Categories {
			for _, s := range c.Cases {
				num += len(s.Checks)
			}
		}
	}
	return num
}

func getStrLen(str string) int {
	return utf8.RuneCountInString(str)
}
