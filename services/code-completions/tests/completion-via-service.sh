#!/bin/bash

ADDR=""
DATA=""
PROMPT=""
DFILE=""
PFILE=""
APIKEY=""
MODEL=""
TEMPERATURE=""
MAX_TOKENS=""
COMPLETION_ID=""
CLIENT_ID=""
PROJECT_PATH="$PWD"
FILE_PROJECT_PATH=""
LANGUAGE_ID=""
STREAM=""
TRIGGER_MODE=""
NO_DEBUG=""
PROTOCOL_V2=""
VERBOSE=""
OUTPUT=""

# 解析FIM格式的文件
function parse_fim_file() {
  local file_path="$1"
  
  # 检查文件是否存在
  if [ ! -f "$file_path" ]; then
    echo "错误: 文件 $file_path 不存在"
    exit 1
  fi
  
  # 读取文件内容
  local file_content=$(cat "$file_path")
    
  if ! echo "$file_content" | grep -q "<｜fim▁hole｜>"; then
    echo "错误: 文件中缺少 <｜fim▁hole｜> 标记"
    exit 1
  fi
    
  # 提取prefix: 从文件开头到<｜fim▁hole｜>的内容
  local hole_marker="<｜fim▁hole｜>"
  local prefix=$(echo "$file_content" | sed -n "1,/$hole_marker/p" | sed '$d')
  
  # 提取suffix: 从<｜fim▁hole｜>到文件结尾的内容
  local suffix=$(echo "$file_content" | sed -n "/$hole_marker/,\$p" | sed '1d')
  
  # 提取cursor_line_prefix: <｜fim▁hole｜>所在行前面的内容
  local hole_line=$(echo "$file_content" | grep -n "$hole_marker" | cut -d: -f1)
  local cursor_line_prefix=$(echo "$file_content" | sed -n "${hole_line}p" | sed "s/$hole_marker.*//")
  
  # 提取cursor_line_suffix: <｜fim▁hole｜>所在行后面的内容
  local cursor_line_suffix=$(echo "$file_content" | sed -n "${hole_line}p" | sed "s/.*$hole_marker//")
  
  # 更新DATA变量中的prompt_options字段，如果不存在则创建
  DATA=`jq --arg prefix "$prefix" --arg suffix "$suffix" --arg cursor_line_prefix "$cursor_line_prefix" --arg cursor_line_suffix "$cursor_line_suffix" '. + {prompt_options: {prefix: $prefix, suffix: $suffix, cursor_line_prefix: $cursor_line_prefix, cursor_line_suffix: $cursor_line_suffix}}' <<< "$DATA"`
  
  # 同时清除prompt字段
  DATA=`jq --arg prompt "" '.prompt = $prompt' <<< "$DATA"`
}

function print_help() {
  echo "Usage: $0 [-a addr] [options]"
  echo "  -a addr: 地址"
  echo "  -p prompt: 提示前缀"
  echo "  -d data: 补全消息"
  echo "  -f prompt-file: 含提示前缀的文件，支持FIM格式(...<｜fim▁hole｜>...)"
  echo "  -F data-file: 含完整补全消息的文件"
  echo "  -o output: 输出文件名"
  echo "  -k apikey: 密钥"
  echo "  -m model: 模型"
  echo "  -t temperature: 温度值"
  echo "  -M max_tokens: 最大令牌数"
  echo "  -i completion_id: 补全ID"
  echo "  -c client_id: 客户端ID"
  echo "  -P project_path: 项目路径"
  echo "  -C file_project_path: 文件项目路径"
  echo "  -l language_id: 语言ID"
  echo "  -r trigger_mode: 触发模式"
  echo "  -s: 开启流式输出"
  echo "  -n: 不输出调试信息"
  echo "  -v: 开启verbose模式"
  echo "  -2: 启用V2版本协议"
  echo "  -h: 帮助"
}
# 初始化选项
while getopts "a:p:d:f:F:k:m:i:c:P:C:l:t:M:r:o:nsv2h" opt; do
  case "$opt" in
    a)
      ADDR="$OPTARG"
      ;;
    p) 
      PROMPT="$OPTARG"
      ;;
    d)
      DATA="$OPTARG"
      ;;
    f) 
      PFILE="$OPTARG"
      ;;
    F)
      DFILE="$OPTARG"
      ;;
    k)
      APIKEY="$OPTARG"
      ;;
    m)
      MODEL="$OPTARG"
      ;;
    t)
      TEMPERATURE="$OPTARG"
      ;;
    M)
      MAX_TOKENS="$OPTARG"
      ;;
    i)
      COMPLETION_ID="$OPTARG"
      ;;
    c)
      CLIENT_ID="$OPTARG"
      ;;
    P)
      PROJECT_PATH="$OPTARG"
      ;;
    C)
      FILE_PROJECT_PATH="$OPTARG"
      ;;
    l)
      LANGUAGE_ID="$OPTARG"
      ;;
    r)
      TRIGGER_MODE="$OPTARG"
      ;;
    s)
      STREAM="true"
      ;;
    n)
      NO_DEBUG="true"
      ;;
    v)
      VERBOSE="true"
      ;;
    2)
      PROTOCOL_V2="true"
      ;;
    o)
      OUTPUT="$OPTARG"
      ;;
    h)
      print_help
      exit 0
      ;;
    *) 
      echo "无效选项"
      print_help
      exit 1
      ;;
  esac
done

DATA_V1='{
  "model": "DeepSeek-Coder-V2-Lite-Base",
  "prompt": "#!/usr/bin/env python\n# coding: utf-8\nimport time\nimport base64\ndef trace(rsp):\n    print",
  "temperature": 0.1,
  "max_tokens": 50,
  "stop": [],
  "beta_mode": false,
  "stream": false,
  "verbose": false,
  "language_id": "python",
  "trigger_mode": "manual",
  "project_path": "'"$PWD"'",
  "file_project_path": "test.py",
  "client_id": "zbc-test",
  "completion_id": "666-94131415"
}'

DATA_V2='{
  "model": "DeepSeek-Coder-V2-Lite-Base",
  "temperature": 0.1,
  "max_tokens": 50,
  "stop": [],
  "beta_mode": false,
  "stream": false,
  "verbose": false,
  "language_id": "python",
  "trigger_mode": "manual",
  "project_path": "'"$PWD"'",
  "file_project_path": "test.py",
  "client_id": "zbc-test",
  "completion_id": "666-94131415",
  "prompt_options": {
    "prefix": "#!/usr/bin/env python\n# coding: utf-8\nimport time\nimport base64\ndef trace(rsp):\n    print",
    "suffix": ""
  }
}'

# 设定消息BODY的样板
if [ X"$DFILE" != X"" ]; then
  DATA="$(cat $DFILE)"
elif [ X"$PROTOCOL_V2" == X"true" ]; then
  DATA="$DATA_V2"
elif [ X"$DATA" == X"" ]; then
  DATA="$DATA_V1"
fi

# 补全消息的提示词部分
if [ X"$PROMPT" != X"" ]; then
  if [ X"$PROTOCOL_V2" == X"true" ]; then
    DATA=`jq --arg newValue "$PROMPT" '.prompt_options.prefix = $newValue' <<< "$DATA"`
  else
    DATA=`jq --arg newValue "$PROMPT" '.prompt = $newValue' <<< "$DATA"`
  fi
elif [ X"$PFILE" != X"" ]; then
  # 检查文件是否包含FIM标记
  if grep -q "<｜fim▁hole｜>" "$PFILE" 2>/dev/null; then
    # 使用FIM解析函数
    parse_fim_file "$PFILE"
  elif [ X"$PROTOCOL_V2" == X"true" ]; then
    DATA=`jq --arg newValue "$(cat $PFILE)" '.prompt_options.prefix = $newValue' <<< "$DATA"`
  else
    # 使用原有逻辑，将整个文件内容作为prompt
    DATA=`jq --arg newValue "$(cat $PFILE)" '.prompt = $newValue' <<< "$DATA"`
  fi
fi

if [ X"$TEMPERATURE" != X"" ]; then
  DATA=`jq --argjson newValue "$TEMPERATURE" '.temperature = $newValue' <<< "$DATA"`
fi

if [ X"$MAX_TOKENS" != X"" ]; then
  DATA=`jq --argjson newValue "$MAX_TOKENS" '.max_tokens = $newValue' <<< "$DATA"`
fi

if [ X"$MODEL" != X"" ]; then
  DATA=`jq --arg newValue "$MODEL" '.model = $newValue' <<< "$DATA"`
fi

if [ X"$CLIENT_ID" != X"" ]; then
  DATA=`jq --arg newValue "$CLIENT_ID" '.client_id = $newValue' <<< "$DATA"`
fi

if [ X"$COMPLETION_ID" == X"" ]; then
  COMPLETION_ID=$(dd if=/dev/urandom bs=16 count=1 2>/dev/null | hexdump -v -e '/1 "%02x"')
fi

if [ X"$COMPLETION_ID" != X"" ]; then
  DATA=`jq --arg newValue "$COMPLETION_ID" '.completion_id = $newValue' <<< "$DATA"`
fi

if [ X"$PROJECT_PATH" != X"" ]; then
  DATA=`jq --arg newValue "$PROJECT_PATH" '.project_path = $newValue' <<< "$DATA"`
fi

if [ X"$FILE_PROJECT_PATH" == X"" ]; then
  if [ X"$PFILE" != X"" ]; then
    # 转换为绝对路径，兼容Windows和Unix系统
    case "$PFILE" in
      /*|[A-Za-z]:/*|[A-Za-z]:\\*)
        # Unix绝对路径或Windows绝对路径
        FILE_PROJECT_PATH="$PFILE"
        ;;
      *)
        # 相对路径，转换为绝对路径
        FILE_PROJECT_PATH="$(pwd)/$PFILE"
        ;;
    esac
  fi
fi

if [ X"$FILE_PROJECT_PATH" != X"" ]; then
  DATA=`jq --arg newValue "$FILE_PROJECT_PATH" '.file_project_path = $newValue' <<< "$DATA"`
fi

if [ X"$LANGUAGE_ID" != X"" ]; then
  DATA=`jq --arg newValue "$LANGUAGE_ID" '.language_id = $newValue' <<< "$DATA"`
fi

if [ X"$TRIGGER_MODE" != X"" ]; then
  DATA=`jq --arg newValue "$TRIGGER_MODE" '.trigger_mode = $newValue' <<< "$DATA"`
fi

if [ X"$STREAM" != X"" ]; then
  DATA=`jq '.stream = true' <<< "$DATA"`
fi

if [ X"$VERBOSE" != X"" ]; then
  DATA=`jq '.verbose = true' <<< "$DATA"`
fi

if [ X"$ADDR" == X"" ]; then
  echo missing '-a/--addr', such as: -a "http://172.16.0.4:5001/v1/completions"
  exit 1
fi

HEADERS=("-H" "Content-Type: application/json" "-H" "X-Costrict-Version: 2.0.0" "-H" "X-Request-Id: $COMPLETION_ID")
if [ X"$APIKEY" != X"" ]; then
  HEADERS+=("-H" "Authorization: Bearer $APIKEY")
fi

OUTPUT_OPT=""
if [ X"$OUTPUT" != X"" ]; then
  OUTPUT_OPT="-o $OUTPUT"
fi

# 执行curl命令
if [ X"$NO_DEBUG" == X"" ]; then
  echo curl $OUTPUT_OPT -sS $ADDR "${HEADERS[@]}" -X POST -d "$DATA" 1>&2
fi
curl $OUTPUT_OPT -sS $ADDR "${HEADERS[@]}" -X POST -d "$DATA"
