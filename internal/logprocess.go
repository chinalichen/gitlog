package logprocess

import (
	"strings"
)

const (
	Spliter = "$%#"
)

type Processor = func(source string) string

func Process(source string) string {
	processors := []Processor{replaceComma, replaceSpliter}

	result := source
	for _, f := range processors {
		result = f(result)
	}

	return result
}

/**
 * 因为 commit message 中可能存在“，”影响 csv 的分割
 * 使用“ ”替换所有的“，”。
 */
func replaceComma(source string) string {
	return strings.ReplaceAll(source, ",", " ")
}

func replaceSpliter(source string) string {
	return strings.ReplaceAll(source, Spliter, ",")
}
