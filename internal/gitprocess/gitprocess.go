package gitprocess

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

const (
	spliter = "$@#"
)

type GitProcessor struct {
	rootPath string
}

func NewGetProcessor(rootPath string) *GitProcessor {
	return &GitProcessor{rootPath: rootPath}
}

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

func removeEmptyLines(source string) (string, error) {
	return strings.ReplaceAll(source, "\n\n", "\n"), nil
}

func (gp *GitProcessor) Clone(dir string, url string) error {
	gitDir := path.Join(gp.rootPath, dir)
	if _, err := os.Stat(gitDir); os.IsExist(err) || err == nil {
		return nil
	}
	cmd := exec.Command("git", "clone", url, gitDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return xerrors.Errorf("clone git %s error: %w", url, err)
	}
	if strings.Contains(string(output), "100%") {
		return xerrors.Errorf("clone git %s error: %w", url, output)
	}
	return nil
}

func (gp *GitProcessor) GitLog(dir string) (string, error) {
	gitDir := path.Join(gp.rootPath, dir)
	cmd := exec.Command("git", GetGitLotArgs()...)
	cmd.Dir = gitDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", xerrors.Errorf("git log error: %w", err)
	}
	result, err := Process(string(output))
	if err != nil {
		return "", xerrors.Errorf("convert git log to csv error: %w", err)
	}
	return result, nil
}
