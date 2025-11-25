#!/bin/sh

COMPLETION_SERVICE_ADDR=""
COMPLETION_APIKEY=""

# 初始化选项
OUTPUT=""
STATS=""
LOGS=""
DETAILS=""
DEBUG="false"

# 如果存在./.env文件，则加载它
if [ -f "./.env" ]; then
    if [ X"$DEBUG" == X"true" ]; then
        echo "加载./.env文件中的环境变量..."
    fi
    source ./.env
fi

function print_help() {
  echo "Usage: $0 [-a addr] [-k apikey] [-o output] [-s] [-l level] [-d] [-g]"
  echo "  -a addr: 地址(default: $COMPLETION_SERVICE_ADDR)"
  echo "  -k apikey: 密钥(default: ${COMPLETION_APIKEY:0:10}***)"
  echo "  -o output: 输出文件名"
  echo "  -l level: 设置日志级别"
  echo "  -s: 获取服务统计数据"
  echo "  -d: 获取补全服务详情"
  echo "  -g: 启用调试输出"
  echo "  -h: 帮助"
}

while getopts "a:k:o:sdl:gh" opt; do
  case "$opt" in
    a)
      COMPLETION_SERVICE_ADDR="$OPTARG"
      ;;
    k)
      COMPLETION_APIKEY="$OPTARG"
      ;;
    o)
      OUTPUT="$OPTARG"
      ;;
    s)
      STATS="true"
      ;;
    l)
      LOGS="$OPTARG"
      ;;
    d)
      DETAILS="true"
      ;;
    g)
      DEBUG="true"
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

if [ X"$COMPLETION_SERVICE_ADDR" == X"" ]; then
  echo missing '-a/--addr', such as: -a "http://172.16.0.4:5001"
  exit 1
fi

if [ X"$DEBUG" == X"true" ]; then
  echo "调试模式已启用"
  echo "服务地址: $COMPLETION_SERVICE_ADDR"
  if [ X"$COMPLETION_APIKEY" != X"" ]; then
    echo "API密钥: ${COMPLETION_APIKEY:0:10}***"
  fi
fi

HEADERS=("-H" "Content-Type: application/json" "-H" "X-Costrict-Version: 2.0.0")
if [ X"$COMPLETION_APIKEY" != X"" ]; then
  HEADERS+=("-H" "Authorization: Bearer $COMPLETION_APIKEY")
fi

OUTPUT_OPT=""
if [ X"$OUTPUT" != X"" ]; then
  OUTPUT_OPT="-o $OUTPUT"
fi

# 检查是否指定了stats、logs或details选项
if [ X"$STATS" == X"" ] && [ X"$LOGS" == X"" ] && [ X"$DETAILS" == X"" ]; then
  echo "错误：必须指定 -s (获取统计数据)、-l level (设置日志级别) 或 -d (获取补全服务详情) 选项"
  print_help
  exit 1
fi

if [ X"$DEBUG" == X"true" ]; then
  echo "准备发送请求..."
  if [ X"$STATS" != X"" ]; then
    echo "请求类型: 获取服务统计数据"
  elif [ X"$LOGS" != X"" ]; then
    echo "请求类型: 设置日志级别为 $LOGS"
  elif [ X"$DETAILS" != X"" ]; then
    echo "请求类型: 获取补全服务详情"
  fi
fi

# 如果指定了stats选项，发送GET请求到/api/stats路径
if [ X"$STATS" != X"" ]; then
  STATS_URL="$COMPLETION_SERVICE_ADDR/api/stats"
  if [ X"$DEBUG" == X"true" ]; then
    echo curl $OUTPUT_OPT -sS $STATS_URL "${HEADERS[@]}" -X GET 1>&2
  fi
  curl $OUTPUT_OPT -sS $STATS_URL "${HEADERS[@]}" -X GET
# 如果指定了logs选项，发送POST请求到/api/logs路径
elif [ X"$LOGS" != X"" ]; then
  LOGS_URL="$COMPLETION_SERVICE_ADDR/api/logs"
  LOGS_DATA="{\"level\": \"$LOGS\"}"
  if [ X"$DEBUG" == X"true" ]; then
    echo curl $OUTPUT_OPT -sS $LOGS_URL "${HEADERS[@]}" -X POST -d "$LOGS_DATA" 1>&2
  fi
  curl $OUTPUT_OPT -sS $LOGS_URL "${HEADERS[@]}" -X POST -d "$LOGS_DATA"
# 如果指定了details选项，发送GET请求到/api/details路径
elif [ X"$DETAILS" != X"" ]; then
  DETAILS_URL="$COMPLETION_SERVICE_ADDR/api/details"
  if [ X"$DEBUG" == X"true" ]; then
    echo curl $OUTPUT_OPT -sS $DETAILS_URL "${HEADERS[@]}" -X GET 1>&2
  fi
  curl $OUTPUT_OPT -sS $DETAILS_URL "${HEADERS[@]}" -X GET
fi

if [ X"$DEBUG" == X"true" ]; then
  echo "请求已完成"
fi
