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
SUFFIX=""
NO_DEBUG=""
OUTPUT=""
FIM=""

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

  # 更新DATA变量中的suffix字段
  DATA=`jq --arg suffix "$suffix" '.suffix = $suffix' <<< "$DATA"`
  
  # 更新DATA变量中的prompt字段
  DATA=`jq --arg prompt "$prefix" '.prompt = $prompt' <<< "$DATA"`
}

function print_help() {
  echo "Usage: $0 {-a addr} [options]"
  echo "  -a addr: 地址"
  echo "  -p prompt: 提示前缀"
  echo "  -d data: 补全消息"
  echo "  -f prompt-file: 含提示前缀的文件"
  echo "  -F data-file: 含完整补全消息的文件"
  echo "  -k apikey: 密钥"
  echo "  -m model: 模型"
  echo "  -t temperature: 温度值"
  echo "  -M max_tokens: 最大令牌数"
  echo "  -s suffix: 后缀"
  echo "  -o output: 输出文件名"
  echo "  -n: 仅输出curl获取的数据内容，不输出调试信息"
  echo "  -i: 启用FIM模式，将prefix和suffix组合为特定格式"
  echo "  -h: 帮助"
}
# 初始化选项
while getopts "a:p:d:f:F:k:m:s:t:M:o:hni" opt; do
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
    s)
      SUFFIX="$OPTARG"
      ;;
    o)
      OUTPUT="$OPTARG"
      ;;
    n)
      NO_DEBUG="true"
      ;;
    i)
      FIM="true"
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

TEMP='{
  "model": "DeepSeek-Coder-V2-Lite-Base",
  "prompt": "#!/usr/bin/env python\n# coding: utf-8\nimport time\nimport base64\ndef trace(rsp):\n    print",
  "temperature": 0.1,
  "max_tokens": 50,
  "stop": [],
  "beta_mode": false,
  "suffix": null
}'

if [ X"$DFILE" != X"" ]; then
  DATA="$(cat $DFILE)"
elif [ X"$DATA" == X"" ]; then
  DATA="$TEMP"
fi

if [ X"$PROMPT" != X"" ]; then
  DATA=`jq --arg newValue "$PROMPT" '.prompt = $newValue' <<< "$DATA"`
elif [ X"$PFILE" != X"" ]; then
  # 检查文件是否包含FIM标记
  if grep -q "<｜fim▁hole｜>" "$PFILE" 2>/dev/null; then
    # 使用FIM解析函数
    parse_fim_file "$PFILE"
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

if [ X"$SUFFIX" != X"" ]; then
  DATA=`jq --arg newValue "$SUFFIX" '.suffix = $newValue' <<< "$DATA"`
fi

# 处理FIM模式
if [ X"$FIM" == X"true" ]; then
  # 获取当前的prompt和suffix值
  CURRENT_PROMPT=$(echo "$DATA" | jq -r '.prompt')
  CURRENT_SUFFIX=$(echo "$DATA" | jq -r '.suffix')
  
  # 如果suffix为null或空，设为空字符串
  if [ "$CURRENT_SUFFIX" = "null" ] || [ -z "$CURRENT_SUFFIX" ]; then
    CURRENT_SUFFIX=""
  fi
  
  # 组合为FIM格式: <｜fim▁begin｜>prefix<｜fim▁hole｜>suffix<｜fim▁end｜>
  FIM_PROMPT="<｜fim▁begin｜>${CURRENT_PROMPT}<｜fim▁hole｜>${CURRENT_SUFFIX}<｜fim▁end｜>"
  
  # 更新DATA中的prompt字段，并清空suffix字段
  DATA=`jq --arg newValue "$FIM_PROMPT" '.prompt = $newValue | .suffix = null' <<< "$DATA"`
fi

if [ X"$ADDR" == X"" ]; then
  echo missing '-a/--addr', such as: -a "http://172.16.0.4:5001/v1/completions"
  exit 1
fi

HEADERS=("-H" "Content-Type: application/json")
if [ X"$APIKEY" != X"" ]; then
  HEADERS+=("-H" "Authorization: Bearer $APIKEY")
fi

OUTPUT_OPT=""
if [ X"$OUTPUT" != X"" ]; then
  OUTPUT_OPT="-o $OUTPUT"
fi

if [ X"$NO_DEBUG" == X"" ]; then
  echo curl $OUTPUT_OPT -sS $ADDR "${HEADERS[@]}" -X POST -d "$DATA" 1>&2
fi
curl $OUTPUT_OPT -sS $ADDR "${HEADERS[@]}" -X POST -d "$DATA" 
