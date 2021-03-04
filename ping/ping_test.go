package ping

import (
	"testing"
)

func TestPing(t *testing.T) {
	var args = []string{"ping", "play.miaoscraft.cn"}
	ret := func(resp string) {
		t.Log(resp)
	}
	Ping(args, ret)
}
