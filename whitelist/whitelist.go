package whitelist

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"io"
	"net/http"
)

func MyID(qq int64, name string, ret func(msg string)) {
	// 查询玩家名字和ID
	name, id, err := getUUID(name)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("向Mojang查询玩家UUID失败: %v", err))
		ret("啊咧？查不到这个账号的UUID呀！")
		return
	}

	// 在数据库中记录
	owner, oldID, err := data.SetWhitelist(qq, id)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("数据库操作失败: %v", err))
		ret("我访问不到数据库啦？！")
		return
	}

	// 若owner是当前处理的用户则说明绑定成功，否则就是失败
	if owner != qq {
		ret(fmt.Sprintf("你想要%q的白名？没门儿！因为已经被[CQ:at,qq=%d]占有啦！", name, owner))
		return
	}

	// 删除旧的白名单
	if oldID != uuid.Nil {
		// 获取旧用户名
		oldName, err := getName(oldID)
		if err != nil {
			cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("向Mojang查询玩家Name失败: %v", err))
			ret("诶诶？？现在查不到之前绑定的游戏名耶")
		} else {
			// 删除旧白名单
			err = data.RemoveWhitelist(oldName)
			if err != nil {
				cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("删除白名单失败: %v", err))
				ret("服务器貌似不想消掉你旧的白名单(ﾟﾍﾟ?)???")
			}
		}
	}

	// 添加白名单
	err = data.AddWhitelist(name)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("添加白名单失败: %v", err))
		ret("添加白名单居然失败了……服务器有她自己的想法_(:з」∠)_")
		return
	}

	ret(fmt.Sprintf("kira~已为您添加白名单: %s", name))
}

func RemoveWhitelist(qq int64, ret func(msg string)) {
	// 删除数据库中的数据
	id, err := data.UnsetWhitelist(qq)
	if err != nil {
		ret(fmt.Sprintf("(ﾟﾍﾟ?)???为QQ=%d的白名单操作数据库时出现了一些问题: %v", qq, err))
		return
	}

	if id != uuid.Nil { // 若这个QQ绑定了白名
		name, err := getName(id)
		if err != nil {
			ret(fmt.Sprintf("(ﾟﾍﾟ?)???查询QQ=%d的游戏名时出现了一些问题: %v", qq, err))
			return
		}

		err = data.RemoveWhitelist(name)
		if err != nil {
			ret(fmt.Sprintf("(ﾟﾍﾟ?)???消除白名单%s时出现了一些问题: %v", name, err))
			return
		}

		ret("消除了" + name + "的白名单")
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
