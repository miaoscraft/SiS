package ping

import (
	"testing"
)

func TestPing(t *testing.T) {
	var args = []string{"ping", "my.hypixel.net"}
	ret := func(resp string) {
		t.Log(resp)
	}
	Ping(args, ret)
}
