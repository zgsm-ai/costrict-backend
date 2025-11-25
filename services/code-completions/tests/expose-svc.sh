#!/bin/sh

# 函数定义区域

# 显示帮助信息的函数
show_help() {
    echo "用法: $0 [-s 服务名] [-n 命名空间] [-p 主机端口] [-t 服务端口] [-e env文件]"
    echo "  -s 指定k8s服务名 (必需)"
    echo "  -n 指定k8s命名空间 (默认: costrict)"
    echo "  -p 指定映射后的主机端口 (默认: 30000)"
    echo "  -t 指定服务端口 (默认: 8080)"
    echo "  -e 指定env文件路径 (默认: ./.env)"
    echo "  -h 显示此帮助信息"
}

# 设置参数值的函数，按优先级从高到低选择第一个非空的值
# 用法: value=$(get_value "$cmd_value" "$env_value" "$default_value")
get_value() {
    for val in "$1" "$2" "$3"; do
        if [ -n "$val" ]; then
            echo "$val"
            return
        fi
    done
}

# 按顺序尝试从网络接口获取主机IP
get_host_ip() {
    interfaces="eth0 ens18 eth1"
    for interface in $interfaces; do
        ip=`ip -4 -o addr show $interface | awk '{print $4}' | cut -d'/' -f1`
        if [ -n "$ip" ]; then
            echo "$ip"
            return
        fi
    done
    echo ""
}

# 信号处理函数：清理子进程
cleanup() {
    echo "接收到终止信号，正在清理子进程..."
    if [ -n "$KUBECTL_PID" ]; then
        echo "终止 kubectl port-forward 进程 (PID: $KUBECTL_PID)"
        kill -9 $KUBECTL_PID 2>/dev/null
    fi
    echo "kubectl port-forward进程已终止，退出expose-svc.sh脚本"
    exit 0
}
# 函数定义区域结束

# 初始化变量
EXPOSE_SERVICE_NAME=""
EXPOSE_SERVICE_PORT=""
EXPOSE_NAMESPACE=""
EXPOSE_HOST_PORT=""

ENV_FILE="./.env"

# 解析命令行参数（最高优先级）
while getopts "s:n:p:t:e:h" opt; do
    case $opt in
    s) CMD_SERVICE_NAME="$OPTARG" ;;
    n) CMD_NAMESPACE="$OPTARG" ;;
    p) CMD_HOST_PORT="$OPTARG" ;;
    t) CMD_SERVICE_PORT="$OPTARG" ;;
    e) ENV_FILE="$OPTARG" ;;
    h)
        show_help
        exit 0
        ;;
    \?)
        echo "无效选项: -$OPTARG" >&2
        show_help
        exit 1
        ;;
    :)
        echo "选项 -$OPTARG 需要参数." >&2
        show_help
        exit 1
        ;;
    esac
done

# 检查ENV_FILE是否以/或.开头，如果不是，则自动补上./前缀
if [ "$ENV_FILE" != "" ] && [ "$ENV_FILE" = "${ENV_FILE#/}" ] && [ "$ENV_FILE" = "${ENV_FILE#.}" ]; then
    ENV_FILE="./$ENV_FILE"
fi

# 加载env文件（第二优先级）
if [ -f "$ENV_FILE" ]; then
    echo "加载env文件: $ENV_FILE"
    source "$ENV_FILE"
fi

# 按照优先级设置参数值：
# 1. 命令行选项（最高优先级）
# 2. ENV_FILE指定的文件中设定的值（第二优先级）
# 3. 默认值（最低优先级）
EXPOSE_SERVICE_NAME=$(get_value "$CMD_SERVICE_NAME" "$EXPOSE_SERVICE_NAME" "")
EXPOSE_NAMESPACE=$(get_value "$CMD_NAMESPACE" "$EXPOSE_NAMESPACE" "costrict")
EXPOSE_HOST_PORT=$(get_value "$CMD_HOST_PORT" "$EXPOSE_HOST_PORT" "30000")
EXPOSE_SERVICE_PORT=$(get_value "$CMD_SERVICE_PORT" "$EXPOSE_SERVICE_PORT" "8080")

HOST_IP=$(get_host_ip)

# 检查必需参数
if [ -z "$EXPOSE_SERVICE_NAME" ]; then
    echo "错误: 必须指定k8s服务名 (-s 选项)"
    show_help
    exit 1
fi

if [ -z "$HOST_IP" ]; then
    echo "错误: HOST_IP为空"
    show_help
    exit 1
fi

echo "配置信息:"
echo "  服务名: $EXPOSE_SERVICE_NAME"
echo "  命名空间: $EXPOSE_NAMESPACE"
echo "  主机端口: $EXPOSE_HOST_PORT"
echo "  服务端口: $EXPOSE_SERVICE_PORT"
echo "  Env文件: $ENV_FILE"
echo "  主机IP: $HOST_IP"
echo ""


# 设置信号陷阱：捕获 SIGINT, SIGTERM, SIGQUIT 等信号
trap cleanup INT TERM QUIT

echo "启动 kubectl port-forward..."
# 启动 kubectl port-forward 并记录其 PID
kubectl port-forward svc/$EXPOSE_SERVICE_NAME $EXPOSE_HOST_PORT:$EXPOSE_SERVICE_PORT -n $EXPOSE_NAMESPACE --address $HOST_IP &
KUBECTL_PID=$!

echo "kubectl port-forward 已启动，PID: $KUBECTL_PID"
echo "将 svc/$EXPOSE_SERVICE_NAME 的 $EXPOSE_SERVICE_PORT 端口映射到 $HOST_IP:$EXPOSE_HOST_PORT"

# 等待 kubectl 进程完成
wait $KUBECTL_PID

echo "kubectl port-forward 进程已正常退出"
