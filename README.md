[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/blum/app?startapp=ref_YbE9XKVsqA)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/BlumBot/main/assets/demo.png)

# 🔥🔥 Blum Bot Auto Claim And Auto Completing Task 🔥🔥

### Tested on Windows and Docker Alpine Os with a 4-core CPU using 5 threads.

**Go Version Tested 1.23.1**

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) Must >= 1.23.

- #### Rename config.yml.example to config.yml.
- #### Rename query.txt.example to query.txt and place your query data.
- #### Rename proxy.txt.example to proxy.txt and place your query data.
- #### If you don’t have a query data, you can obtain it from [Telegram Web Tools](https://github.com/ehhramaaa/telegram-web-tools)
- #### It is recommended to use an IP info token to improve request efficiency when checking IPs.

## Features

|       Feature        | Supported |
| :------------------: | :-------: |
|        Proxy         |    ✅     |
|    Multithreading    |    ✅     |
|    Use Query Data    |    ✅     |
|    Auto Check In     |    ✅     |
|     Auto Farming     |    ✅     |
|    Auto Play Game    |    ✅     |
| Auto Completing Task |    ✅     |
| Auto Claim Ref Point |    ✅     |
|  Random User Agent   |    ✅     |

## [Settings](https://github.com/ehhramaaa/BlumBot/blob/main/config.yml)

|     Settings     |                           Description                            |
| :--------------: | :--------------------------------------------------------------: |
|  **USE_PROXY**   |                     For Activated Proxy Mode                     |
| **IPINFO_TOKEN** | For Increase Check Ip Efficiency. Put Your Own Token If You Have |
| **GAME_POINTS**  |     Amount point wanna get from Game (e.g. MIN:200, MAX:300)     |
|  **MAX_THREAD**  |        Max Thread Worker Run Parallel Recommend 10 - 100         |
| **RANDOM_SLEEP** |  Delay before the next lap (e.g. MIN:3600, MAX:7200) in second   |

## Installation

```shell
git clone https://github.com/ehhramaaa/BlumBot.git
cd BlumBot
go run .
```

## Or you can do build application by typing:

Windows:

```shell
go build -o BlumBot.exe
```

Linux:

```shell
go build -o BlumBot
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```