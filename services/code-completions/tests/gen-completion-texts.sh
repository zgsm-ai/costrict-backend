#!/bin/bash

# 生成代码补全结果的脚本
# 从案例文件目录获取文件内容，从模型输出结果目录获取补全内容，替换或拼接后输出到结果目录

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -r, --response-dir DIR 指定模型输出结果目录 (必要参数)"
    echo "  -d, --data-dir DIR     指定案例文件目录 (默认: data)"
    echo "  -o, --output-dir DIR   指定输出目录 (默认: completion-results+$response_dir)"
    echo "  -h, --help            显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -d data -r responses"
    echo "  $0 -d test_cases -r model_outputs -o results"
}

# 默认参数值
DATA_DIR="data"
RESPONSE_DIR=""
OUTPUT_DIR=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--data-dir)
            DATA_DIR="$2"
            shift 2
            ;;
        -r|--response-dir)
            RESPONSE_DIR="$2"
            shift 2
            ;;
        -o|--output-dir)
            OUTPUT_DIR="$2"
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

# 检查jq命令是否存在
if ! command -v jq >/dev/null 2>&1; then
    echo "错误: jq命令不存在，此脚本需要jq来处理JSON数据"
    echo "请安装jq后再运行此脚本"
    echo "安装方法:"
    echo "  Ubuntu/Debian: sudo apt-get install jq"
    echo "  CentOS/RHEL: sudo yum install jq"
    echo "  macOS: brew install jq"
    echo "  Windows: 使用WSL或下载jq二进制文件"
    exit 1
fi

# 检查案例文件目录是否存在
if [ -z "$DATA_DIR" ]; then
    echo "错误: 必须指定案例文件目录"
    show_usage
    exit 1
fi

if [ ! -d "$DATA_DIR" ]; then
    echo "错误: 案例文件目录不存在: $DATA_DIR"
    exit 1
fi

# 检查响应文件目录是否存在
if [ -z "$RESPONSE_DIR" ]; then
    echo "错误: 必须指定模型输出结果目录"
    show_usage
    exit 1
fi

if [ ! -d "$RESPONSE_DIR" ]; then
    echo "错误: 模型输出结果目录不存在: $RESPONSE_DIR"
    exit 1
fi

# 如果未指定输出目录，使用默认值 "completion-results" 加上 $RESPONSE_DIR
if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR="completion-results/$(basename "$RESPONSE_DIR")"
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "处理补全结果..."
processed=0
success=0

# 获取数据目录中的所有文件
file_list=$(find "$DATA_DIR" -type f | sort)
total_files=$(echo "$file_list" | wc -l)

for filepath in $file_list; do
    # 获取文件名（不包含路径）
    filename=$(basename "$filepath")
    base_name="${filename%.*}"
    response_file="$RESPONSE_DIR/$base_name.json"
    
    processed=$((processed + 1))
    echo "[$processed/$total_files] 处理文件: $filename"
    
    # 读取案例文件内容
    case_content=$(cat "$filepath")
    
    # 检查响应文件是否存在
    if [ ! -f "$response_file" ]; then
        echo "  警告: 响应文件不存在: $response_file"
        continue
    fi
    
    # 从响应文件中提取补全内容
    completion_text=$(jq -r '.choices[0].text // empty' "$response_file" 2>/dev/null)
    
    if [ -z "$completion_text" ]; then
        echo "  警告: 无法从响应文件中提取补全内容: $response_file"
        continue
    fi
    
    # 检查案例文件内容中是否包含<｜fim▁hole｜>标记
    if [[ "$case_content" == *"<｜fim▁hole｜>"* ]]; then
        # 替换标记为补全内容
        result_content="${case_content//"<｜fim▁hole｜>"/$completion_text}"
        echo "  替换了<｜fim▁hole｜>标记"
    else
        # 拼接到末尾
        result_content="$case_content$completion_text"
        echo "  拼接补全内容到末尾"
    fi
    
    # 输出结果到TXT文件
    result_txt_file="$OUTPUT_DIR/$base_name.txt"
    echo "$result_content" > "$result_txt_file"
    
    echo "  结果已保存到: $result_txt_file"
    success=$((success + 1))
done

echo ""
echo "处理完成!"
echo "总共处理了 $processed 个文件"
echo "成功生成了 $success 个结果文件"
echo "结果保存在: $OUTPUT_DIR"