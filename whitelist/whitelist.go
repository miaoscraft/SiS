package whitelist

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
)

func MyID(qq int64, name string, ret func(msg string)) {
	// 查询玩家名字和ID
	Name, id, err := GetUUID(name)
	if err != nil {
		Logger.Errorf("向Mojang查询玩家UUID失败: %v", err)
		ret(fmt.Sprintf("捡到个纸团\n( ^ ω ^) \n≡⊃§⊂≡ \n打开看一眼\n( ^ ω ^)\n⊃|" + name + "|⊂\n不认识这个id呢\n( ^ ω ^) \n≡⊃§⊂≡\n \n§\n ¶\n　∩( ^ ω ^)"))
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

	// 在数据库中记录
	owner, err := data.SetWhitelist(qq, id, onOldID)
	if err != nil {
		Logger.Errorf("设置白名单失败: %v", err)
		ret("白名单貌似没有成功加上欸，怎么办ʕ •ᴥ•ʔ")
		return
	}

	// 若owner不为0说明绑定失败
	if owner != 0 {
		//自己已经占用了
		if owner == qq {
			ret(fmt.Sprintf("{\\__/}\n( • . •)\n/ >%s\n呐，你的白名单", Name))
			return
		}
		//被别人占用
		if len(Name) < 3 {
			ret(fmt.Sprintf("白名单%s现在在[CQ:at,qq=%d]手上", Name, owner))
		} else {
			ret(fmt.Sprintf("{\\__/}\n( • . •)\n/ >%s\n你要这个吗？\n\n{\\__/}\n( • - •)\n%s< \\\n这是[CQ:at,qq=%d]的", Name, Name[len(Name)-3:], owner))
		}
		return
	}
	// 添加白名单
	err = data.AddWhitelist(Name)
	if err != nil {
		Logger.Errorf("添加白名单失败: %v", err)
		ret(fmt.Sprintf("添加白名单失败: %v", err))
		if err := data.UnsetWhitelist(qq, func(_ uuid.UUID) error { return nil }); err != nil {
			Logger.Errorf("从数据库删除记录失败: %v", err)
			ret(fmt.Sprintf("从数据库删除记录失败: %v", err))
		}
		return
	}
	Logger.Infof("添加白名单%q成功", Name)
	ret(fmt.Sprintf("{\\__/}\n( • . •)\n/ >%s\n呐，你的白名单", Name))
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

		Logger.Infof("删除白名单%q成功", name)
		ret(name + "，你白名单(号)没了")
		return nil
	}
	// 删除数据库中的数据
	err := data.UnsetWhitelist(qq, onHas)
	if err != nil {
		Logger.Errorf("删除白名单失败: %v", err)
		ret("我的系统又出问题了(つД`)ノ")
		return
	}
}

// GetUUID 查询玩家的UUID
func GetUUID(name string) (string, uuid.UUID, error) {
	var id uuid.UUID

	// 发送请求
	d, status, err := get("https://api.mojang.com/users/profiles/minecraft/" + name)
	if err != nil {
		return "", id, err
	}
	defer d.Close()

	// 检查返回码
	if status != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", status)
	}

	// 解析json返回值
	err = json.NewDecoder(d).Decode(&struct {
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
	d, status, err := get("https://sessionserver.mojang.com/session/minecraft/profile/" + hex.EncodeToString(UUID[:]))
	if err != nil {
		return "", err
	}
	defer d.Close()

	// 检查返回码
	if status != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", status)
	}

	var resp struct{ Name string }
	// 解析json返回值
	err = json.NewDecoder(d).Decode(&resp)
	if err != nil {
		return "", err
	}

	return resp.Name, nil
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

var Logger interface {
	Error(str string)
	Errorf(format string, args ...interface{})

	Waring(str string)
	Waringf(format string, args ...interface{})

	Info(str string)
	Infof(format string, args ...interface{})

	Debug(str string)
	Debugf(format string, args ...interface{})
}
