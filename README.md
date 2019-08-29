# SiS
基于酷Q的MC服务器综合管理插件，涵盖了非常方便的许多功能。

## 功能
1. 白名单  
群员在群内发送`MyID=<正版游戏名>`自助获取白名单，例如`MyID=Tnze`；
群员更换账号时同样发送`MyID=<正版游戏名>`更新，例如`MyID=Xi_Xi_Mi`；
群员退群时白名单将被自动取消。**MyID指令无需艾特机器人** 。
2. Ping  
群员在群内**艾特机器人** 并发送指令`ping`查询服务器状态，例如`@robot ping`，
插件自己实现了MC PingAndList功能，不调用第三方API；
群员可在ping后接第三方服务器地址，例如`@robot play.miaoscraft.cn`，
则ping目标转至群员指定的服务器；
群员可进一步指定服务器端口，有两种格式均可行：
	- `@robot ping play.miaoscraft.cn:25565`
	- `@robot ping play.miaoscraft.cn 25565`。

## 配置文件
请务必修改配置文件

## 数据接口
~~数据库暂仅支持MySQL，计划支持SQLite。~~
改为使用bolt数据库  
支持`toml`配置文件格式，通俗易懂类似`.ini`。  
通过RCON协议与MC服务器通信，同时支持官方服务器与Bukkit系列，无需服务器插件。
