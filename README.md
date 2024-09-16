[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/blum/app?startapp=ref_YbE9XKVsqA)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/BlumBot/main/assets/Screenshot_3.png)

## Recommendation before use

# 🔥🔥 Go Version Tested 1.23.1 🔥🔥

## Features

|       Feature        | Supported |
| :------------------: | :-------: |
|    Multithreading    |    ✅     |
|    Use Query Data    |    ✅     |
|    Auto Check In     |    ✅     |
|     Auto Farming     |    ✅     |
|    Auto Play Game    |    ✅     |
| Auto Completing Task |    ✅     |
| Auto Claim Ref Point |    ⏳     |
| Auto Connect Wallet  |    ⏳     |
|  Random User Agent   |    ✅     |

## [Settings](https://github.com/ehhramaaa/BlumBot/blob/main/config.yml)

|     Settings     |                          Description                          |
| :--------------: | :-----------------------------------------------------------: |
| **GAME_POINTS**  |      Amount point wanna get from Game (e.g. MIN:200, MAX:300)       |
|  **MAX_THREAD**  |       Max Thread Worker Run Parallel Recommend 10 - 100       |
| **RANDOM_SLEEP** | Delay before the next lap (e.g. MIN:3600, MAX:7200) in second |

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) **version > 1.22**

## Installation

You can download the [**repository**](https://github.com/ehhramaaa/BlumBot.git) by cloning it to your system and installing the necessary dependencies:

```shell
git clone https://github.com/ehhramaaa/BlumBot.git
cd BlumBot
go mod tidy
go run .
```

Then you can do build application by typing:

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

**If You Want Auto Select Choice In Terminal**

For Option 1

```shell
go run . -c 1
```

For Option 2

```shell
go run . -c 2
```