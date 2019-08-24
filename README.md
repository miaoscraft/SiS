# CoolQ-Golang-Plugin
这里是你用Golang酷Q插件开发的起点

## 开始
赶快点击右上角的`Use this template`绿色按钮开始吧！
用本模板新建一个项目（到你自己的Github账号上），然后将你的项目克隆至本地。
或者直接下载本模板项目。

## 安装环境
1. 下载并安装[Go语言编译器](https://golang.google.cn/)；
2. 下载并安装[gcc编译器](http://tdm-gcc.tdragon.net/)；  
3. 运行脚本安装SDK以及cqcfg小工具：运行`setup.bat`。

## 修改路径
要修改的地方有几处：
1. go.mod文件第一行，改为你自己项目的地址
2. app.go文件main函数前`// cqp:`开头的注释，修改名称、版本、作者和简介
3. app.go文件init函数内，修改你的AppID
4. build.bat脚本第5行DevDir，后面的路径要修改为你的酷Q的dev文件夹的路径

## 启动酷Q的开发者模式
请查看酷Q官方的[文档](https://d.cqp.me/Pro/%E5%BC%80%E5%8F%91/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8)

## 测试运行
运行一下`build.bat`试试吧，这个脚本会帮你编译插件、自动生成app.json，
然后帮你把app.dll和app.json移动到酷Q的开发目录下。

最后，在酷Q的菜单-应用管理中，点击重载应用，你应该就能看到你的插件了。