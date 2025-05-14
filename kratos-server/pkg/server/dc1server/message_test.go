package dc1server

import (
	"testing"
)

func TestAns(t *testing.T) {
	answer := &Answer{Uuid: "12321323", Status: CODE_SUCCESS, Msg: "device identified"}
	t.Log(string(answer.ToMsg()))

}
