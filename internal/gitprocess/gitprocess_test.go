package gitprocess

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	gitLogStr = `a621336$@#2af69e1$@#chinalichen$@#chinalichen@126.com$@#chinalichen$@#2021-09-17$@#chinalichen$@#chinalichen@126.com$@#chinalichen$@#2021-09-17$@#add string process code
 4 files changed, 75 insertions(+), 11 deletions(-)

2af69e1$@#4256812$@#lc$@#lc@lcs-iMac.local$@#lc$@#2021-09-15$@#lc$@#lc@lcs-iMac.local$@#lc$@#2021-09-15$@#init
 3 files changed, 28 insertions(+), 1 deletion(-)

4256812$@#$@#lichen$@#chinalichen@126.com$@#chinalichen$@#2021-09-15$@#GitHub$@#noreply@github.com$@#noreply$@#2021-09-15$@#Initial commit
 1 file changed, 1 insertion(+)
`
)

func TestProcess(t *testing.T) {
	actual, err := Process(gitLogStr)
	if err != nil {
		t.Errorf("replace comma error: %s", err)
	}
	expect := `a621336,2af69e1,chinalichen,chinalichen@126.com,chinalichen,2021-09-17,chinalichen,chinalichen@126.com,chinalichen,2021-09-17,add string process code,4,75,11
2af69e1,4256812,lc,lc@lcs-iMac.local,lc,2021-09-15,lc,lc@lcs-iMac.local,lc,2021-09-15,init,3,28,1
4256812,,lichen,chinalichen@126.com,chinalichen,2021-09-15,GitHub,noreply@github.com,noreply,2021-09-15,Initial commit,1,1
`
	if actual != expect {
		t.Errorf("expect: \n%s\n, but got: \n%s\n", expect, actual)
	}
}

func TestReplaceComma(t *testing.T) {
	actual, err := replaceComma(gitLogStr)
	if err != nil {
		t.Errorf("replace comma error: %s", err)
	}
	index := strings.Index(actual, ",")
	if index != -1 {
		t.Errorf("acutal has invalid ',' at %d ", index)
	}
}

func TestGitOperations(t *testing.T) {
	gp := GitProcessor{rootPath: os.TempDir()}
	dir := "gitlog"
	url := "https://github.com/chinalichen/gitlog.git"
	if err := gp.Clone(dir, url); err != nil {
		t.Fatalf("git clone error %v", err)
	}
	csv, err := gp.GitLog(dir)
	if err != nil {
		t.Fatalf("git log error %v", err)
	}
	fmt.Printf("gitlog.csv: %s", csv)
}
