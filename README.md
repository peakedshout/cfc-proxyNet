***
# cfc-proxyNet <img src="./asset/box/cfcproxynet_logo.png">
###### *A communication connection agent based on golang development, the goal is simple, efficient, stable, secure, and extensible.*
***
### [简体中文](./README_CN.md)/English
***
## What is cfc-proxyNet？
- This is a simple proxy communication tool, mainly proxy http, https, socks protocol communication
- Using [therecipe/qt](https://github.com/therecipe/qt) to write qt client
- Working model: rawData <=> local listener service <=> proxy listener service <=> target service
***
## What can cfc-proxyNet do?
- Proxy communicates with http, https, and socks
- over
***
## Why is this thing here?
- I'm not satisfied that go-CFC's features stop updating. It shouldn't stop there. As the first library in my life, I want to develop it as much as possible
- [cFc-fileManage](https://github.com/peakedshout/cfc-fileManage) is [BUDAI-AZ](https://github.com/BUDAI-AZ) and I cooperation work, but the effect is not very good, As a rookie, the back end vs. front end was really tough, so I decided to try to write a simple ui page application myself
- But obviously, I underestimated, even if go to develop qt, I was also deeply tortured, such as a variety of environment configuration, compilation failure and so on, as a qt zero based rookie, now the degree has been OK, enough
- Well, I needed to use it too, so I built it myself, and it's not as good as other people's, but I feel proud to use my own software
- If you enjoy it, please give me ✨
***
## How to use it?
- At present, there are only macos amd64 and windows amd64, as these are the only machines I have, and I think they are relatively popular
- Server
  - Put the corresponding application on your server and run the developed configuration file. [here](./asset/complete_resources/server) you can see the finished product about the server
  - ``` 
    ./cfc-proxyNetS -c ./config.json
    ```
  - If you don't know how to fill out a profile, Can refer to the configuration notes -> ([CN](./asset/complete_resources/server/configCN/config.json)/[EN](./asset/complete_resources/server/configEN/config.json))
- MacOS amd64 client
    - Run dmg or untar and just run it
    - It's not hard to do it according to the page, right
    - MacOS needs to manually go to the system setting and point the proxy setting to the running port. I know this is very tedious, but I don't know how to use golang to set the system setting of macos. Please tell me if you know, thank you
- Windows amd64 client
    - unzip, run
    - It's not hard to do it according to the page, right
    - Windows will help you complete the proxy Settings, normal shutdown application will also help you turn off the proxy mode, if you find that the application crashes and you can not access the Internet, you can run it and close the normal, or go to the proxy Settings to turn off the proxy switch
- Android
    - Do not do
    - Do not know how to interconnect vpn mode, do not know whether the existing compilation mode can handle successfully
    - If you know how to operate, please let me know, thanks
- iOS
    - Do not do
    - I don't have an iOS device
***
## How to compile?
- If you're not satisfied with the finished product, you can pull the source code and compile it yourself, but I don't recommend that. You'll have plenty of time to fiddle around and find alternatives
    - 0.Ok, you still want to compile, please follow these steps:
    - 1.To install golang, I am currently using 1.19.+, you can try other versions, there is nothing too high version in it, probably?
    - 2.Install [therecipe/qt](https://github.com/therecipe/qt), it's a nightmare, please prepare
    - 3.run ``go run ./client``, If no errors are reported, then./asset/complete_resources should have a finished product that runs against the current system. If you want to compile to another system, Please use [therecipe/qt](https://github.com/therecipe/qt) approach to compile ./asset/cfc-proxyNet
    - x.Since I haven't implemented any other systems, there may be some problems. Good luck
    - y.If you want to compile the server, run  ``go run ./server``, Some common systems are already encapsulated in it, if you want to compile other systems, please run``go build ./asset/cfc-proxyNetS``
    - z.Of course, if you give up thinking, you can run``go run ./main.go``, It performs the ``go run./client`` and ``go run./server`` operations. If you don't want the finished product to be a folder, you can package it manually without affecting the correctness of the program
***
## Features you want to support (later?)
- [x] http/https
- [x] socks
- [x] windows/macos client
- [ ] Android client
- [ ] Customize the proxy route
- [ ] Perfect it
***