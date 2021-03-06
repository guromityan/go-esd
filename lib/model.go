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

// NewTests テスト仕様生成
func NewTests(name string) *Tests {
	return &Tests{
		Name:   name,
		path:   "",
		Genres: make([]Genre, 0),
	}
}

// NewGenre ジャンル生成
func NewGenre(name string, te *Tests) *Genre {
	g := Genre{
		Num:        len(te.Genres),
		Genre:      name,
		Categories: make([]Category, 0),
	}
	te.Genres = append(te.Genres, g)
	return &g
}

// NewCategory カテゴリ生成
func NewCategory(name string, ge *Genre) *Category {
	c := Category{
		Num:      len(ge.Categories),
		Category: name,
		Cases:    make([]Case, 0),
	}
	ge.Categories = append(ge.Categories, c)
	return &c
}

// NewCase ケース生成
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

// LastGenre 最後のジャンル取得
func (t *Tests) LastGenre() (*Genre, error) {
	if len(t.Genres) == 0 {
		return nil, fmt.Errorf("None Genre: %v", t.Name)
	}
	return &t.Genres[len(t.Genres)-1], nil
}

// LastCategory 最後のカテゴリ取得
func (t *Tests) LastCategory() (*Category, error) {
	genre, err := t.LastGenre()
	if err != nil {
		return nil, err
	}
	if len(genre.Categories) == 0 {
		return nil, fmt.Errorf("No Categories: %v", genre.Genre)
	}
	return &genre.Categories[len(genre.Categories)-1], nil
}

// LastCase 最後のケース取得
func (t *Tests) LastCase() (*Case, error) {
	category, err := t.LastCategory()
	if err != nil {
		return nil, err
	}
	if len(category.Cases) == 0 {
		return nil, fmt.Errorf("No Cases: %v", category.Category)
	}
	return &category.Cases[len(category.Cases)-1], nil

}

// SetPath 生成した Excel を保存するディレクトリを設定
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

// GetMaxCategory 最長カテゴリを取得
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

// GetMaxCase 最長ケースを取得
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

// AddStep ステップを追加
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

// GetStepCheckHeight ステップとチェックから最適な高さを取得
func (c *Case) GetStepCheckHeight() []int {
	stepNum := 0
	for _, s := range c.Steps {
		for range strings.Split(s, "\n") {
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

// AddCheck チェックを追加
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

// GetMaxStep 最長のステップを取得
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

// GetMaxCheck 最長のチェックを取得
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

func getStrLen(str string) int {
	return utf8.RuneCountInString(str)
}
