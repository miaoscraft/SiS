# SiS
基于酷Q的MC服务器综合管理插件，涵盖了非常方便的许多功能。
- 中文错误输出，插件简单易用
- Q群成员绑定服务器白名单，玩家自助添加，退群自动取消
- MC服务器状态查询，延迟以及在线玩家汇报（敬请期待）
- 管理员可在群内执行后台命令（敬请期待）
- 支持分离游戏群与管理群，两个群都能执行命令
- 所有命令都支持权限管理，为每个管理员和命令设置权限（敬请期待）

## 配置文件
conf.toml文件为配置文件
```toml
# 游戏群群号
GroupID = 123456789

# 管理群群号（可选）
AdminID = 123456789

# Ping工具配置
[Ping]
DefaultServer = "play.miaoscraft.cn"
# 超时设置，为0时禁用
Timeout = "60s"

# RCON配置
[RCON]
Address = "127.0.0.1"
Password = "your_password"

# 数据库配置
[Database]
Address = "127.0.0.1"
User = "your_mysql_username"
Password = "your_mysql_password"
Schema = "数据库库名"
```

## 数据接口
数据库暂仅支持MySQL，计划支持SQLite。  
支持`toml`配置文件格式，通俗易懂类似`.ini`。  
通过RCON协议与MC服务器通信，同时支持官方服务器与Bukkit系列，无需服务器插件。
