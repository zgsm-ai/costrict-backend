#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
代码库定义查询测试脚本

该脚本用于查询代码片段中，所有在本项目其它文件定义的函数、对象等。
通过调用代码库索引器API，获取代码片段中引用的符号定义信息。

主要功能：
- 发送代码片段到代码库索引器
- 获取符号定义信息，包括文件路径、名称、类型、内容等
- 支持指定客户端ID、代码库路径、文件路径等参数

使用方法：
    python codebase-test.py [选项]

命令行选项：
    --clientId        客户端ID (默认: 123)
    --codebasePath    代码库路径 (默认: g:\\tmp\\projects\\c\\redis)
    --filePath        待补全的文件路径 (默认: g:\\tmp\\projects\\c\\redis\\src\\ae.c)
    --snippet         代码片段的URL编码 (默认: Redis ae.c文件头部代码)
    --snippet-file    从指定文件读取代码片段内容
    --token           Authorization头中的access_token (默认: 使用内置token)
    --token-file      从指定文件读取token内容
    --host            服务器主机地址 (默认: localhost)
    --port            服务器端口 (默认: 9001)
    --help            显示帮助信息

作者：AI Assistant
创建时间：2025-10-21
"""

import urllib.request
import urllib.parse
import json
import argparse

snippetDefault = """
/*
Copyright © 2022 zbc <zbc@sangfor.com.cn>
*/
package main

import (
	_ "github.com/zgsm-ai/smc/cmd"
	"github.com/zgsm-ai/smc/cmd/common"
)

func main() {
	common.Execute()
}

"""

# 创建命令行参数解析器
parser = argparse.ArgumentParser(description='代码库定义查询测试脚本')

# 添加命令行选项
parser.add_argument('--clientId', type=str, default='f1c91100ff3cda160edb4b3f3f8b5934f3dce52b252e33db27d86a1fe36f45c3.f4be6518',
                    help='客户端ID (默认: f1c91100ff3cda160edb4b3f3f8b5934f3dce52b252e33db27d86a1fe36f45c3.f4be6518)')
parser.add_argument('--codebasePath', type=str, default='/mnt/d/shenma/smc',
                    help='代码库路径 (默认: /mnt/d/shenma/smc)')
parser.add_argument('--filePath', type=str, default='/mnt/d/shenma/smc/main.go',
                    help='待补全的文件路径 (默认: /mnt/d/shenma/smc/main.go)')
parser.add_argument('--snippet', type=str, default=snippetDefault,
                    help='代码片段的URL编码 (默认: Redis ae.c文件头部代码)')
parser.add_argument('--snippet-file', type=str, help='从指定文件读取代码片段内容')
parser.add_argument('--token', type=str, help='Authorization头中的access_token')
parser.add_argument('--token-file', type=str, help='从指定文件读取token内容')
parser.add_argument('--host', type=str, default='localhost', help='服务器主机地址 (默认: localhost)')
parser.add_argument('--port', type=str, default='9001', help='服务器端口 (默认: 9001)')

# 解析命令行参数
args = parser.parse_args()

# 使用解析后的参数
clientid = args.clientId
codebasePath = args.codebasePath # 工作区路径
filePath = args.filePath # 待补全的文件路径
snippet = args.snippet # 代码片段urlencode（最好是文件头部的import 部分 + 当前待补全的代码片段）
host = args.host # 服务器主机地址
port = args.port # 服务器端口
token = args.token # Authorization头中的access_token

# 如果指定了snippet-file，从文件中读取代码片段内容
if args.snippet_file:
    try:
        with open(args.snippet_file, 'r', encoding='utf-8') as f:
            snippet = f.read()
    except Exception as e:
        print("读取代码片段文件失败: {}".format(e))
        exit(1)

# 如果指定了token-file，从文件中读取token内容
if args.token_file:
    try:
        with open(args.token_file, 'r', encoding='utf-8') as f:
            token = f.read().strip()
    except Exception as e:
        print("读取token文件失败: {}".format(e))
        exit(1)

url = "http://{}:{}/codebase-indexer/api/v1/search/definition?clientId={}&codebasePath={}&filePath={}&codeSnippet={}".format(host, port, clientid, codebasePath, filePath, urllib.parse.quote(snippet))

headers = {
  'Authorization': 'Bearer {}'.format(token)
}

# 创建请求对象
req = urllib.request.Request(url, headers=headers)

# 发送请求并获取响应
try:
    with urllib.request.urlopen(req) as response:
        response_data = response.read().decode('utf-8')
        print(response_data)
except urllib.error.URLError as e:
    print("请求失败: {}".format(e))
except Exception as e:
    print("发生错误: {}".format(e))
# response
"""
{
    "code": "0",
    "message": "ok",
    "success": true,
    "data": {
        "list": [
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\monotonic.h",
                "name": "monotime",
                "type": "definition.class",
                "content": "typedef uint64_t monotime;",
                "position": {
                    "startLine": 22,
                    "startColumn": 1,
                    "endLine": 22,
                    "endColumn": 27
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeEventLoop",
                "type": "definition.class",
                "content": "typedef struct aeEventLoop {\n    int maxfd;   /* highest file descriptor currently registered */\n    int setsize; /* max number of file descriptors tracked */\n    long long timeEventNextId;\n    int nevents; /* Size of Registered events */\n    aeFileEvent *events; /* Registered events */\n    aeFiredEvent *fired; /* Fired events */\n    aeTimeEvent *timeEventHead;\n    int stop;\n    void *apidata; /* This is used for polling API specific data */\n    aeBeforeSleepProc *beforesleep;\n    aeBeforeSleepProc *aftersleep;\n    int flags;\n    void *privdata[2];\n} aeEventLoop;",
                "position": {
                    "startLine": 79,
                    "startColumn": 9,
                    "endLine": 93,
                    "endColumn": 2
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeEventLoop",
                "type": "definition.class",
                "content": "typedef struct aeEventLoop {\n    int maxfd;   /* highest file descriptor currently registered */\n    int setsize; /* max number of file descriptors tracked */\n    long long timeEventNextId;\n    int nevents; /* Size of Registered events */\n    aeFileEvent *events; /* Registered events */\n    aeFiredEvent *fired; /* Fired events */\n    aeTimeEvent *timeEventHead;\n    int stop;\n    void *apidata; /* This is used for polling API specific data */\n    aeBeforeSleepProc *beforesleep;\n    aeBeforeSleepProc *aftersleep;\n    int flags;\n    void *privdata[2];\n} aeEventLoop;",
                "position": {
                    "startLine": 79,
                    "startColumn": 1,
                    "endLine": 93,
                    "endColumn": 15
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeTimeEvent",
                "type": "definition.class",
                "content": "typedef struct aeTimeEvent {\n    long long id; /* time event identifier. */\n    monotime when;\n    aeTimeProc *timeProc;\n    aeEventFinalizerProc *finalizerProc;\n    void *clientData;\n    struct aeTimeEvent *prev;\n    struct aeTimeEvent *next;\n    int refcount; /* refcount to prevent timer events from being\n  \t\t   * freed in recursive time event calls. */\n} aeTimeEvent;",
                "position": {
                    "startLine": 60,
                    "startColumn": 9,
                    "endLine": 70,
                    "endColumn": 2
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeTimeEvent",
                "type": "definition.class",
                "content": "typedef struct aeTimeEvent {\n    long long id; /* time event identifier. */\n    monotime when;\n    aeTimeProc *timeProc;\n    aeEventFinalizerProc *finalizerProc;\n    void *clientData;\n    struct aeTimeEvent *prev;\n    struct aeTimeEvent *next;\n    int refcount; /* refcount to prevent timer events from being\n  \t\t   * freed in recursive time event calls. */\n} aeTimeEvent;",
                "position": {
                    "startLine": 60,
                    "startColumn": 1,
                    "endLine": 70,
                    "endColumn": 15
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeFileEvent",
                "type": "definition.class",
                "content": "typedef struct aeFileEvent {\n    int mask; /* one of AE_(READABLE|WRITABLE|BARRIER) */\n    aeFileProc *rfileProc;\n    aeFileProc *wfileProc;\n    void *clientData;\n} aeFileEvent;",
                "position": {
                    "startLine": 52,
                    "startColumn": 9,
                    "endLine": 57,
                    "endColumn": 2
                }
            },
            {
                "filePath": "g:\\tmp\\projects\\c\\redis\\src\\ae.h",
                "name": "aeFileEvent",
                "type": "definition.class",
                "content": "typedef struct aeFileEvent {\n    int mask; /* one of AE_(READABLE|WRITABLE|BARRIER) */\n    aeFileProc *rfileProc;\n    aeFileProc *wfileProc;\n    void *clientData;\n} aeFileEvent;",
                "position": {
                    "startLine": 52,
                    "startColumn": 1,
                    "endLine": 57,
                    "endColumn": 15
                }
            }
        ]
    }
}

"""
