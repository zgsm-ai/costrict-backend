#!/bin/bash

# 测试LLM补全性能的脚本
# 以串行方式测试补全模型对各种代码文件进行补全时的效果和性能

# 默认参数值
COMPLETION_EXPOSE=true
COMPLETION_URL=""
COMPLETION_APIKEY=""
DATA_DIR="./data"

# 如果存在./.env文件，则加载它
if [ -f "./.env" ]; then
    echo "加载./.env文件中的环境变量..."
    source ./.env
fi

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --no-expose      在执行测试前不调用 expose-completion.sh 暴露服务端口 (默认: false)"
    echo "  --key KEY        指定调用completion-via-service.sh时使用的API密钥"
    echo "  --url URL        指定调用completion-via-service.sh时使用的URL，即COMPLETION_URL变量的值"
    echo "  -d, --data DIR   指定数据目录 (默认: ./data)"
    echo "  -h, --help       显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 --no-expose   # 执行测试前不暴露服务端口"
    echo "  $0 --key myapikey # 指定API密钥"
    echo "  $0 --url http://localhost:32088/v1/completions # 指定URL"
    echo "  $0 -d ./my-data # 指定数据目录"
    echo "  $0               # 默认执行测试前先暴露服务端口"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-expose)
            COMPLETION_EXPOSE=false
            shift
            ;;
        --key)
            COMPLETION_APIKEY="$2"
            shift 2
            ;;
        --url)
            COMPLETION_URL="$2"
            shift 2
            ;;
        -d|--data)
            DATA_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 检查必要文件是否存在
if [ ! -f "./completion-via-service.sh" ]; then
    echo "错误：找不到 completion-via-service.sh 文件"
    exit 1
fi

KEY_OPT=""
if [ ! -z "$COMPLETION_APIKEY" ]; then
    KEY_OPT="-k $COMPLETION_APIKEY"
fi

# 创建结果目录
RESULTS_DIR="completion_perf_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$RESULTS_DIR"

echo "配置信息:"
echo "  数据目录: $DATA_DIR"
echo ""

# 创建性能数据JSON文件
PERF_DATA_FILE="$RESULTS_DIR/perf_data.json"
echo "[" > "$PERF_DATA_FILE"

# 用于跟踪是否是第一个JSON对象，避免在最后一个对象后添加逗号
is_first=true


# 函数：根据文件扩展名确定语言类型
get_language_by_extension() {
    local filename="$1"
    local extension="${filename##*.}"
    
    case "$extension" in
        "go") echo "go" ;;
        "c") echo "c" ;;
        "cpp"|"cc"|"cxx") echo "c++" ;;
        "py") echo "python" ;;
        "java") echo "java" ;;
        "js") echo "javascript" ;;
        "ts") echo "typescript" ;;
        "lua") echo "lua" ;;
        *) echo "unknown" ;;
    esac
}

# 根据 COMPLETION_EXPOSE 参数决定是否执行 expose-completion.sh
if [ "$COMPLETION_EXPOSE" = true ]; then
    if [ ! -f "./expose-completion.sh" ]; then
        echo "错误：找不到 expose-completion.sh 文件"
        exit 1
    fi
    # 检查 32088 端口是否已经被占用
    echo "检查 32088 端口状态..."
    if netstat -tuln | grep -q ":32088 "; then
        echo "端口 32088 已被占用，跳过 expose-completion.sh 的执行"
        EXPOSE_PID=""
    else
        echo "端口 32088 未被占用，启动 expose-completion.sh 后台进程..."
        ./expose-completion.sh &
        EXPOSE_PID=$!

        # 等待一下，确保 expose-completion.sh 启动完成
        sleep 3

        # 检查 expose-completion.sh 是否成功启动
        if ! kill -0 $EXPOSE_PID 2>/dev/null; then
            echo "错误：expose-completion.sh 启动失败"
            exit 1
        fi

        echo "expose-completion.sh 已启动，PID: $EXPOSE_PID"
    fi
else
    echo "COMPLETION_EXPOSE=false，跳过 expose-completion.sh 的执行"
    EXPOSE_PID=""
fi

# 串行执行测试
echo "开始串行执行LLM补全性能测试..."

# 获取数据目录中的所有文件
file_list=$(find "$DATA_DIR" -type f | sort)
total_files=$(echo "$file_list" | wc -l)
current_test=0

for filepath in $file_list; do
    # 获取文件名（不包含路径）
    filename=$(basename "$filepath")
    
    # 根据文件扩展名确定语言类型
    language=$(get_language_by_extension "$filename")
    
    # 跳过未知文件类型
    if [ "$language" = "unknown" ]; then
        echo "跳过未知文件类型: $filename"
        continue
    fi
    
    current_test=$((current_test + 1))
    # 计算输入文件内容长度（字符数）
    input_length=$(wc -c < "$filepath")
    
    # 记录开始时间
    start_time=$(date +%s%3N)  # 毫秒级时间戳
    # 执行补全请求
    response_file="$RESULTS_DIR/${filename%.*}.json"

    echo "[$current_test/$total_files] 测试文件: $filename (语言: $language, 大小: ${input_length}字节)"
    
    # 构建completion-via-service.sh命令
    cmd="./completion-via-service.sh -v -f "$filepath" -l $language -a "$COMPLETION_URL" $KEY_OPT -o "$response_file""
    echo "  $cmd"
    # 使用completion-via-service.sh发送请求并保存响应
    if eval "$cmd"; then
        # 记录结束时间
        end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
        status="success"
        # 尝试从响应中提取补全内容长度
        if grep -q '"choices"' "$response_file"; then
            # 如果是JSON响应，提取补全内容
            completion_length=$(jq -r '.choices[0].text | length' "$response_file" 2>/dev/null || echo "0")
        else
            # 如果不是JSON响应，可能是错误信息
            completion_length="0"
            if grep -q "error\|Error\|ERROR" "$response_file"; then
                status="error"
            fi
        fi
        echo "  请求成功，响应时间: ${response_time}ms, 状态: $status, 补全内容长度: $completion_length"
    else
        # 请求失败
        end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
        status="failed"
        completion_length="0"
        echo "  请求失败，响应时间: ${response_time}ms"
    fi
    
    # 将性能数据添加到JSON文件
    if [ "$is_first" = true ]; then
        is_first=false
    else
        echo "," >> "$PERF_DATA_FILE"
    fi
    
    # 创建JSON对象
    cat >> "$PERF_DATA_FILE" << EOF
    {
        "filename": "$filename",
        "response_time_ms": $response_time,
        "status": "$status"
    }
EOF
    
    # 添加延迟，避免请求过于频繁
    sleep 1
done

# 完成JSON数组
echo "]" >> "$PERF_DATA_FILE"

echo "所有测试已完成，响应文件已保存到: $RESULTS_DIR"

# 停止后台执行的 expose-completion.sh（如果存在）
if [ -n "$EXPOSE_PID" ]; then
    echo "停止 expose-completion.sh 进程..."
    if kill $EXPOSE_PID 2>/dev/null; then
        echo "expose-completion.sh 进程已停止"
    else
        echo "警告：无法停止 expose-completion.sh 进程"
    fi
else
    echo "expose-completion.sh 未由本脚本启动，无需停止"
fi

# 调用gen-completion-reports.sh生成CSV文件和性能报告
echo "使用gen-completion-reports.sh生成CSV文件和性能报告..."
if [ -f "./gen-completion-reports.sh" ]; then
    ./gen-completion-reports.sh -d "$DATA_DIR" -o "$RESULTS_DIR" -r "$RESULTS_DIR" -j "$PERF_DATA_FILE"
    echo "性能报告生成完成"
else
    echo "警告: 找不到gen-completion-reports.sh脚本，跳过报告生成"
fi

# 调用gen-completion-texts.sh生成补全结果文件
echo "使用gen-completion-texts.sh生成补全结果文件..."
if [ -f "./gen-completion-texts.sh" ]; then
    ./gen-completion-texts.sh -d "$DATA_DIR" -r "$RESULTS_DIR"
    echo "补全结果文件生成完成"
else
    echo "警告: 找不到gen-completion-texts.sh脚本，跳过补全结果文件生成"
fi

echo "所有任务已完成，响应文件已保存到: $RESULTS_DIR"
echo "性能数据JSON文件: $PERF_DATA_FILE"
exit 0
