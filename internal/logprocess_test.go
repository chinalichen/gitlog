package logprocess

import (
	"strings"
	"testing"
)

const (
	gitLogStr = `
2af69e1,4256812,lc,lc@lcs-iMac.local,lc,2021-09-15,lc,lc@lcs-iMac.local,lc,2021-09-15,init
3 files changed, 28 insertions(+), 1 deletion(-)

4256812,,lichen,chinalichen@126.com,chinalichen,2021-09-15,GitHub,noreply@github.com,noreply,2021-09-15,Initial commit
1 file changed, 1 insertion(+)
`
)

func TestProcess(t *testing.T) {

}

func TestReplaceComma(t *testing.T) {
	actual := replaceComma(gitLogStr)
	index := strings.Index(actual, ",")
	if index != -1 {
		t.Errorf("acutal has invalid ',' at %d ", index)
	}
}
