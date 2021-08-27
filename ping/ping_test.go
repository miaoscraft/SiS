package ping

import (
	"github.com/BaiMeow/SimpleBot/message"
	"testing"
)

func TestPing(t *testing.T) {
	var args = message.New().Text("ping").Text("play.miaoscraft.cn")
	ret := func(resp string) {
		t.Log(resp)
	}
	Ping(args, ret)
}
