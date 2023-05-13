***
# cfc-proxyNet <img src="./asset/box/cfcproxynet_logo.png">
###### *一个基于go-CFC的代理网络工具，支持http、https、socks代理，目标是简单、方便、好用的。*
***
### 简体中文/[English](./README.md)
***
## 什么是cfc-proxyNet？
- 这是一个简单的代理通信工具，主要代理http、https、socks协议的通信（懂得都懂）
- 采用[therecipe/qt](https://github.com/therecipe/qt)的方式去编写qt客户端
- 工作模型：rawData <=> local listener service <=> proxy listener service <=> target service
***
## cfc-proxyNet能做什么？
- 代理http、https、socks协议的通信
- 没了
***
## 为什么会出现这玩意？
- 我不满足go-CFC的功能停止更新，它不应该止步于此，作为人生第一个库，我想尽可能地去开发
- [cfc-fileManage](https://github.com/peakedshout/cfc-fileManage)是我与[BUDAI-AZ](https://github.com/BUDAI-AZ)的合作作品，但效果不是很好，作为两个菜鸟，后端和前端的对决实在是太难受，所以我打算尝试自己写一个简单的ui页面程序
- 但显然，我低估了，就算是采用go去开发qt，我也被深深的折磨，比如各种环境配置，编译失败等等，作为一个qt零基础的菜鸟，现在的程度已经可以了，够了
- 嗯，这玩意我也需要去用，所以自己做了一个，虽然比不上其他人的，但自己用自己的软件有种自豪感
- 如果你用的愉快，请麻烦给个✨吧
***
## 怎么使用？
- 目前只有macos amd64 和 windows amd64 ,因为本人手头只有这两种机器，并且，我想这类机器是比较大众的
- Server
  - 将对应的程序，放在你的服务器上运行制定的配置文件，在[这里](./asset/complete_resources/server)可以查看有关服务端的成品
  - ```
    ./cfc-proxyNetS -c ./config.json
    ```
  - 如果你不知道如何填写配置文件，可以参考下配置注释->([CN](./asset/complete_resources/server/configCN/config.json)/[EN](./asset/complete_resources/server/configEN/config.json))
- MacOS amd64 client
  - 运行dmg或者是解压tar运行即可
  - 根据页面进行操作，这没什么困难的
  - MacOS需要手动去系统设置，将代理设置指向运行的端口，我知道这很繁琐，但我不清楚怎么用golang去设置macos的系统设置，如果你知道请告诉我，谢谢
- Windows amd64 client
  - 解压，运行
  - 根据页面进行操作，这没什么困难的
  - Windows会帮你完成代理设置，正常关闭程序也会帮你关闭代理模式，如果你发现程序奔溃导致无法上网，可以运行它再正常关闭，或者去代理设置把代理开关关上
- Android
  - 没做
  - 不知道vpn模式如何对接，不知道现有编译模式能不能处理成功
  - 如果你知道如果操作，请告诉我，谢谢
- iOS
  - 没做
  - 本人没有iOS的设备
***
## 怎么编译？
- 如果你对现有的成品不满意，你可以拉去源码自行编译，但我不建议你这么做，折腾编译的时间足够去寻找代替品了
  - 0.好吧，你还是执意去编译了，那请你按照以下步骤：
  - 1.安装golang，目前我使用的是1.19.+，你可以尝试其他的版本，里面没有用什么太高版本的东西，大概？
  - 2.安装[therecipe/qt](https://github.com/therecipe/qt)，这是一个噩梦，请做好心理准备
  - 3.运行 ``go run ./client``，如果没有报错，那么./asset/complete_resources应该就有相应的成品了，它是根据当前系统去运行的，如果想编译其他系统，请用[therecipe/qt](https://github.com/therecipe/qt)的方法去编译./asset/cfc-proxyNet
  - x.因为我没有执行过其他的系统，可能多少会有点问题，祝你好运。
  - y.如果你想编译服务端，请运行 ``go run ./server``，里面已经封装好了一些常见的系统，如果你想编译其他系统，请去``go build ./asset/cfc-proxyNetS``
  - z.当然，如果放弃了思考，可以运行``go run ./main.go``，它会执行``go run ./client``和``go run ./server``的操作。如果你嫌弃编译的成品是个文件夹，你可以手动打包，这不会影响程序的正确性
***
## 想支持的功能（后续完成？）
  - [x] http/https协议
  - [ ] socks协议
  - [x] windows/macos客户端
  - [ ] Android客户端
  - [ ] 自定义代理路线
  - [ ] 完善它
***