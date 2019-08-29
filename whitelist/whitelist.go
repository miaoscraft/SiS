package whitelist

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
)

func MyID(qq int64, name string, ret func(msg string)) {
	// 查询玩家名字和ID
	name, id, err := getUUID(name)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("向Mojang查询玩家UUID失败: %v", err))
		ret("我不要你觉得，我要我觉得" + name + "是个假名字")
		return
	}

	onOldID := func(oldID uuid.UUID) error {
		// 删除用户的旧白名单
		oldName, err := getName(oldID)
		if err != nil {
			return fmt.Errorf("向Mojang查询玩家Name失败: %v", err)
		}

		// 删除旧白名单
		err = data.RemoveWhitelist(oldName)
		if err != nil {
			return fmt.Errorf("删除白名单失败: %v", err)
		}

		return nil
	}

	onSuccess := func() error {
		// 添加白名单
		err = data.AddWhitelist(name)
		if err != nil {
			return fmt.Errorf("添加白名单失败: %v", err)
		}
		return nil
	}

	// 在数据库中记录
	owner, err := data.SetWhitelist(qq, id, onOldID, onSuccess)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("设置白名单失败: %v", err))
		ret("白名单貌似没有成功加上欸，怎么办ʕ •ᴥ•ʔ")
		return
	}

	// 若owner是当前处理的用户则说明绑定成功，否则就是失败
	if owner != qq {
		if len(name) < 3 {
			ret(fmt.Sprintf("白名单%s现在在[CQ:at,qq=%d]手上", name, owner))
		} else {
			ret(fmt.Sprintf(`{\\__/}
( • . •)
/ >%s
你要这个吗？

{\\__/}
( • - •)
%s< \\
这是[CQ:at,qq=%d]的`, name, "..."+name[len(name)-3:], owner))
		}
		return
	}
	ret(fmt.Sprintf(`{\\__/}
( • . •)
/ >%s
呐，你的白名单`, name))
}

func RemoveWhitelist(qq int64, ret func(msg string)) {
	onHas := func(ID uuid.UUID) error {
		name, err := getName(ID)
		if err != nil {
			return fmt.Errorf("查询QQ%d游戏名失败: %v", qq, err)
		}

		err = data.RemoveWhitelist(name)
		if err != nil {
			return fmt.Errorf("删除%s白名单失败: %v", name, err)
		}

		ret(name + "，你白名单(号)没了")
		return nil
	}
	// 删除数据库中的数据
	err := data.UnsetWhitelist(qq, onHas)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("删除白名单失败: %v", err))
		ret("我的系统又出问题了(つД`)ノ")
		return
	}
}

// getUUID 查询玩家的UUID
func getUUID(name string) (string, uuid.UUID, error) {
	var id uuid.UUID

	// 发送请求
	data, status, err := get("https://api.mojang.com/users/profiles/minecraft/" + name)
	if err != nil {
		return "", id, err
	}
	defer data.Close()

	// 检查返回码
	if status != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", status)
	}

	// 解析json返回值
	err = json.NewDecoder(data).Decode(&struct {
		Name *string
		ID   *uuid.UUID
	}{&name, &id})
	if err != nil {
		return name, id, err
	}

	return name, id, err
}

// getName 查询玩家的Name
func getName(UUID uuid.UUID) (string, error) {
	data, status, err := get("https://api.mojang.com/user/profiles/" + hex.EncodeToString(UUID[:]) + "/names")
	if err != nil {
		return "", err
	}
	defer data.Close()

	// 检查返回码
	if status != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", status)
	}

	var resp []struct{ Name string }
	// 解析json返回值
	err = json.NewDecoder(data).Decode(&resp)
	if err != nil {
		return "", err
	}

	if len(resp) < 1 {
		return "", errors.New("(ﾟﾍﾟ?)???没有查询到值")
	}

	return resp[0].Name, nil
}

// 发送GET请求
func get(url string) (io.ReadCloser, int, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	// Golang默认的User-agent被屏蔽了
	request.Header.Set("User-agent", "SiS")

	// 发送Get请求
	resp, err := new(http.Client).Do(request)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.StatusCode, nil
}
