package lib

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
}

// Category カテゴリ
type Category struct {
	Num      int
	Category string
	Cases    []Case
}

// Case テスト手順
type Case struct {
	Num    int
	Case   string
	Steps  []string
	Checks []string
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

func (g *Genre) LastCategory() *Category {
	return &g.Categories[len(g.Categories)-1]
}

func (c *Category) LastCase() *Case {
	return &c.Cases[len(c.Cases)-1]
}
