package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MDParse Markdown を解析し構造体に格納
func MDParse(filename string) (*Tests, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Colud not open %v: %v", filename, err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	var tests *Tests
	var mode string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// 空行
		if line == "" {
			continue
		}

		// コードブロック
		if strings.HasPrefix(line, "```") {
			continue
		}

		// 引用はコメント扱い
		if strings.HasPrefix(line, ">") {
			mode = "comment"
			continue
		}

		// テストファイル
		if strings.HasPrefix(line, "# ") {
			mode = "tests"
			filename := line[2:]
			tests = NewTests(filename)
			continue
		}

		// テストジャンル
		if strings.HasPrefix(line, "## ") {
			mode = "genre"
			name := line[3:]
			NewGenre(name, tests)
			continue
		}

		// テストカテゴリ
		if strings.HasPrefix(line, "### ") {
			mode = "category"
			name := line[4:]
			genre, err := tests.LastGenre()
			if err != nil {
				return nil, err
			}
			NewCategory(name, genre)
			continue
		}

		// テストケース
		if strings.HasPrefix(line, "#### ") {
			mode = "case"
			name := line[5:]
			catagory, err := tests.LastCategory()
			if err != nil {
				return nil, err
			}
			NewCase(name, catagory)
			continue
		}

		// テストチェック
		if strings.HasPrefix(line, "- [ ] ") {
			mode = "check"
			lcase, err := tests.LastCase()
			if err != nil {
				return nil, err
			}
			lcase.AddCheck(line[6:], true)
			continue
		}

		// テストステップ
		i := strings.Index(line, ".")
		if i > 0 {
			_, err := strconv.Atoi(line[:i])
			if err == nil {
				mode = "step"
				lcase, err := tests.LastCase()
				if err != nil {
					return nil, err
				}
				lcase.AddStep(line[i+2:], true)
			} else {
				lcase, err := tests.LastCase()
				if err != nil {
					return nil, err
				}
				if mode == "step" {
					lcase.AddStep(line, false)
				} else if mode == "check" {
					lcase.AddCheck(line, false)
				}
			}
			continue
		}

		if mode == "step" {
			lcase, err := tests.LastCase()
			if err != nil {
				return nil, err
			}
			lcase.AddStep(line, false)
			continue
		}

		if mode == "check" {
			lcase, err := tests.LastCase()
			if err != nil {
				return nil, err
			}
			lcase.AddCheck(line, false)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("File scan error: %v", err)
	}

	return tests, nil
}
