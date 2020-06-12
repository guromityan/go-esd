package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MDParser(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Colud not open %v: %v", filename, err)
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

			name := line[4:]
			NewCategory(name, tests.LastGenre())
			continue
		}

		// テストケース
		if strings.HasPrefix(line, "#### ") {
			mode = "case"
			name := line[5:]
			NewCase(name, tests.LastGenre().LastCategory())
			continue
		}

		// テストチェック
		if strings.HasPrefix(line, "- [ ] ") {
			mode = "check"
			checks := tests.LastGenre().LastCategory().LastCase().Checks
			checks = append(checks, line)
			continue
		}

		// テストステップ
		i := strings.Index(line, ".")
		if i > 0 {
			_, err := strconv.Atoi(line[:i])
			if err == nil {
				mode = "step"
				lastCase := tests.LastGenre().LastCategory().LastCase()
				lastCase.Steps = append(lastCase.Steps, line[i+2:])
			} else {
				if mode == "step" {
					lastCase := tests.LastGenre().LastCategory().LastCase()
					if len(lastCase.Steps) != 0 {
						lastCase.Steps[len(lastCase.Steps)-1] = fmt.Sprintf("%v\n%v", lastCase.Steps[len(lastCase.Steps)-1], line)
					}
				}
			}
			continue
		}

		if mode == "step" {
			lastCase := tests.LastGenre().LastCategory().LastCase()
			lastCase.Steps[len(lastCase.Steps)-1] = fmt.Sprintf("%v\n%v", lastCase.Steps[len(lastCase.Steps)-1], line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("File scan error: %v", err)
	}

	return nil
}
