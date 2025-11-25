#!/bin/bash

# 测试LLM补全性能的脚本
# 以串行方式测试补全模型对各种代码文件进行补全时的效果和性能

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -u, --url URL      设置LLM服务URL (默认: https://oneapi.sangfor.com/v1/completions)"
    echo "  -k, --key KEY      设置LLM服务API密钥"
    echo "  -m, --model MODEL  设置LLM模型名称 (默认: qwen/qwen2.5-coder-7b-instruct)"
    echo "  -f, --fim          启用FIM(Fill-In-Middle)模式"
    echo "  -o, --output DIR   指定结果输出目录 (默认: 自动生成带时间戳的目录)"
    echo "  -d, --data DIR     指定数据目录 (默认: ./data)"
    echo "  -h, --help         显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -u https://api.example.com/v1/completions -k your-api-key -model gpt-4"
    echo "  $0 --url https://api.example.com/v1/completions --key your-api-key --model gpt-4 --fim"
    echo "  $0 -u https://api.example.com/v1/completions -k your-api-key -o ./my-results"
    echo "  $0 -u https://api.example.com/v1/completions -k your-api-key -d ./my-data"
}

# 默认参数值
LLM_URL=""
LLM_KEY=""
LLM_NAME=""
LLM_USE_FIM=false

RESULTS_DIR=""
DATA_DIR="./data"

# 如果存在./.env文件，则加载它
if [ -f "./.env" ]; then
    echo "加载./.env文件中的环境变量..."
    source ./.env
fi

# 检查必要文件是否存在
if [ ! -f "./completion-via-llm.sh" ]; then
    echo "错误：找不到 completion-via-llm.sh 文件"
    exit 1
fi

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--url)
            LLM_URL="$2"
            shift 2
            ;;
        -k|--key)
            LLM_KEY="$2"
            shift 2
            ;;
        -m|--model)
            LLM_NAME="$2"
            shift 2
            ;;
        -f|--fim)
            LLM_USE_FIM=true
            shift
            ;;
        -o|--output)
            RESULTS_DIR="$2"
            shift 2
            ;;
        -d|--data)
            DATA_DIR="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "错误: 未知选项 $1"
            show_usage
            exit 1
            ;;
    esac
done

KEY_OPT=""
FIM_OPT=""
# 检查必要参数
if [ ! -z "$LLM_KEY" ]; then
    KEY_OPT="-k $LLM_KEY"
fi

# 设置FIM选项
if [ "$LLM_USE_FIM" = true ]; then
    FIM_OPT="-i"
fi

echo "配置信息:"
echo "  URL: $LLM_URL"
echo "  模型: $LLM_NAME"
echo "  API密钥: ${LLM_KEY:0:10}..."
echo "  FIM模式: $LLM_USE_FIM"
echo "  数据目录: $DATA_DIR"
echo ""

# 创建结果目录
if [ -z "$RESULTS_DIR" ]; then
    RESULTS_DIR="llm_perf_$(date +%Y%m%d_%H%M%S)"
fi
mkdir -p "$RESULTS_DIR"

# 创建性能数据JSON文件
PERF_DATA_FILE="$RESULTS_DIR/perf_data.json"
echo "[" > "$PERF_DATA_FILE"

# 串行执行测试
echo "开始串行执行LLM补全性能测试..."

# 获取数据目录中的所有文件
file_list=$(find "$DATA_DIR" -type f | sort)
total_files=$(echo "$file_list" | wc -l)
current_test=0

# 用于跟踪是否是第一个JSON对象，避免在最后一个对象后添加逗号
is_first=true

for filepath in $file_list; do
    # 获取文件名（不包含路径）
    filename=$(basename "$filepath")
    
    current_test=$((current_test + 1))
    # 计算输入文件内容长度（字符数）
    input_length=$(wc -c < "$filepath")
    
    # 记录开始时间
    start_time=$(date +%s%3N)  # 毫秒级时间戳
    # 执行补全请求
    response_file="$RESULTS_DIR/${filename%.*}.json"
    echo "[$current_test/$total_files] 测试文件: $filename"
    # 使用completion-via-llm.sh发送请求并保存响应
    if ./completion-via-llm.sh -f "$DATA_DIR/$filename" $KEY_OPT $FIM_OPT -m "$LLM_NAME" -a "$LLM_URL" -o "$response_file"; then
        # 记录结束时间
        end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
        status="success"
        echo "  请求成功，响应时间: ${response_time}ms"
    else
        # 请求失败
        end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
        status="failed"
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

# 调用gen-llm-reports.sh生成CSV文件和性能报告
echo "生成CSV文件和性能报告..."
if [ -f "./gen-llm-reports.sh" ]; then
    echo "调用gen-llm-reports.sh生成报告..."
    ./gen-llm-reports.sh -d "$DATA_DIR" -o "$RESULTS_DIR" -r "$RESULTS_DIR" -j "$PERF_DATA_FILE"
else
    echo "警告: 找不到gen-llm-reports.sh脚本，跳过报告生成"
fi

# 调用gen-completion-texts.sh生成补全结果文件
echo "使用gen-completion-texts.sh生成补全结果文件..."
if [ -f "./gen-completion-texts.sh" ]; then
    ./gen-completion-texts.sh -d "$DATA_DIR" -r "$RESULTS_DIR"
    echo "补全结果文件生成完成"
else
    echo "警告: 找不到gen-completion-texts.sh脚本，跳过补全结果文件生成"
fi

exit 0
