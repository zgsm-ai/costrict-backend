#!/bin/bash

# 生成参考文件的脚本
# 从案例文件目录获取文件内容，输出到参考文件目录

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -d, --data-dir DIR     指定案例文件目录 (默认: data)"
    echo "  -o, --output-dir DIR   指定输出目录 (必要参数)"
    echo "  -h, --help            显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -d data -R refs"
    echo "  $0 -d test_cases -o output_files"
}

# 默认参数值
DATA_DIR="data"
OUTPUT_DIR=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--data-dir)
            DATA_DIR="$2"
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

if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR="completion-results/$(basename "$DATA_DIR")"
fi

# 创建参考文件目录
mkdir -p "$OUTPUT_DIR"

echo "处理参考文件..."
processed=0
success=0

# 获取数据目录中的所有文件
file_list=$(find "$DATA_DIR" -type f | sort)
total_files=$(echo "$file_list" | wc -l)

for filepath in $file_list; do
    # 获取文件名（不包含路径）
    filename=$(basename "$filepath")
    base_name="${filename%.*}"
    
    processed=$((processed + 1))
    echo "[$processed/$total_files] 处理文件: $filename"
    
    # 读取案例文件内容
    case_content=$(cat "$filepath")
    
    # 将案例文件内容保存到参考文件目录
    reference_file="$OUTPUT_DIR/$base_name.txt"
    echo "$case_content" > "$reference_file"
    echo "  参考文件已保存到: $reference_file"
    
    success=$((success + 1))
done

echo ""
echo "处理完成!"
echo "总共处理了 $processed 个文件"
echo "成功生成了 $success 个参考文件"
echo "参考文件保存在: $OUTPUT_DIR"