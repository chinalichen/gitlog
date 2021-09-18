package logprocess

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

const (
	spliter = "$@#"
)

func GetGitLotArgs() []string {
	formatArgs := strings.Join([]string{"%h", "%p", "%an", "%ae", "%al", "%as", "%cN", "%ce", "%cl", "%cs", "%s"}, spliter)
	formatStr := fmt.Sprintf("--pretty=format:'%s'", formatArgs)
	return []string{"log", "--shortstat", formatStr}
}

type Processor = func(source string) (string, error)

func Process(source string) (string, error) {
	processors := []Processor{
		replaceComma,
		replaceSpliter,
		parseFile,
		parseDeletionOnly,
		parseInsertion,
		parseDeletion,
		removeEmptyLines,
	}

	result := source
	for index, f := range processors {
		if res, err := f(result); err != nil {
			return "", xerrors.Errorf("process %d error: %w, input:%s", index, err, result)
		} else {
			result = res
		}
	}

	return result, nil
}

/**
 * 因为 commit message 中可能存在“，”影响 csv 的分割
 * 使用“ ”替换所有的“，”。
 */
func replaceComma(source string) (string, error) {
	return strings.ReplaceAll(source, ",", " "), nil
}

func replaceSpliter(source string) (string, error) {
	return strings.ReplaceAll(source, spliter, ","), nil
}

func genRegexpReplacer(expr string, repl string) func(source string) (string, error) {
	return func(source string) (string, error) {
		reg, err := regexp.Compile(expr)
		if err != nil {
			return "", err
		}
		ok := reg.ReplaceAll([]byte(source), []byte(repl))
		return string(ok), nil
	}
}

var parseFile = genRegexpReplacer(`\n ([0-9]+) file[s]? `, ",$1")
var parseDeletionOnly = genRegexpReplacer(`changed  ([0-9]+) deletion[s]?\(\+\)`, ",0,$1")
var parseInsertion = genRegexpReplacer(`changed  ([0-9]+) insertion[s]?\(\+\)`, ",$1")
var parseDeletion = genRegexpReplacer(`  ([0-9]+) deletion[s]?\(-\)`, ",$1")

// func formatFileAndLines(source string) (string, error) {
// 	fileExp, err := regexp.Compile(`\n ([0-9]+) file[s]? `)
// 	if err != nil {
// 		return "", err
// 	}
// 	fileOk := fileExp.ReplaceAll([]byte(source), []byte(",$1"))

// 	deletionOnlyExp, err := regexp.Compile(`changed  ([0-9]+) deletion[s]?\(\+\)`)
// 	if err != nil {
// 		return "", err
// 	}
// 	deletionOnlyOk := deletionOnlyExp.ReplaceAll(fileOk, []byte(",0,$1"))

// 	insertionExp, err := regexp.Compile(`changed  ([0-9]+) insertion[s]?\(\+\)`)
// 	if err != nil {
// 		return "", err
// 	}
// 	insertionOk := insertionExp.ReplaceAll(deletionOnlyOk, []byte(",$1"))

// 	deletionExp, err := regexp.Compile(`  ([0-9]+) deletion[s]?\(-\)`)
// 	if err != nil {
// 		return "", err
// 	}
// 	deletionOk := deletionExp.ReplaceAll(insertionOk, []byte(",$1"))

// 	return string(deletionOk), nil
// }

func removeEmptyLines(source string) (string, error) {
	return strings.ReplaceAll(source, "\n\n", "\n"), nil
}
