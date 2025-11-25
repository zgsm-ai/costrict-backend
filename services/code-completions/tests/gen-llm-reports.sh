#!/bin/bash

# 生成LLM性能报告的脚本
# 从数据目录获取输入文件，根据输出结果生成CSV文件和总结报告

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -r, --response-dir DIR  指定响应文件目录 (必要参数)"
    echo "  -d, --data-dir DIR   指定数据目录 (默认: data)"
    echo "  -o, --output DIR    指定输出目录 (默认: 使用响应文件目录)"
    echo "  -j, --perf-json-file FILE  指定性能数据JSON文件路径(默认: $response-dir/perf_data.json)"
    echo "  -h, --help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -r responses"
    echo "  $0 -d data -r responses -o results"
    echo "  $0 -r responses -j perf_data.json"
}

# 默认参数值
DATA_DIR="data"
OUTPUT_DIR=""
RESPONSE_DIR=""
PERF_JSON_FILE=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--data-dir)
            DATA_DIR="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -r|--response-dir)
            RESPONSE_DIR="$2"
            shift 2
            ;;
        -j|--perf-json-file)
            PERF_JSON_FILE="$2"
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

# 检查响应文件目录是否存在
if [ -z "$RESPONSE_DIR" ]; then
    echo "错误: 必须指定响应文件目录"
    show_usage
    exit 1
fi

if [ ! -d "$RESPONSE_DIR" ]; then
    echo "错误: 响应文件目录不存在: $RESPONSE_DIR"
    exit 1
fi

# 检查数据目录是否存在
if [ ! -d "$DATA_DIR" ]; then
    echo "错误: 数据目录不存在: $DATA_DIR"
    exit 1
fi

# 如果未指定输出目录，使用响应文件目录
if [ -z "$OUTPUT_DIR" ]; then
    OUTPUT_DIR="$RESPONSE_DIR"
fi

# 如果未指定JSON性能文件，使用响应文件目录下的perf_data.json
if [ -z "$PERF_JSON_FILE" ]; then
    PERF_JSON_FILE="$RESPONSE_DIR/perf_data.json"
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 创建结果文件
result_file="$OUTPUT_DIR/perf_results.csv"
echo "filename,filesize,fulltime,prompt_tokens,completion_tokens" > "$result_file"

echo "处理测试结果..."
processed=0

# 总体统计变量
total_count=0
total_success=0
total_time_sum=0
total_time_valid=0
total_prompt_tokens=0
total_completion_tokens=0

# 响应时间数组（用于计算P90和P99）
declare -a response_times

# 创建关联数组来存储从perf-json-file中获取的fulltime数据
declare -A fulltime_data

# 如果提供了JSON文件，读取fulltime数据
if [ -n "$PERF_JSON_FILE" ] && [ -f "$PERF_JSON_FILE" ]; then
    echo "读取JSON性能数据文件中的fulltime: $PERF_JSON_FILE"
    
    # 获取数组长度
    total_files=$(jq '. | length' "$PERF_JSON_FILE" 2>/dev/null || echo "0")
    
    # 遍历JSON数组，存储filename和对应的fulltime
    for ((i=0; i<total_files; i++)); do
        filename=$(jq -r ".[$i].filename" "$PERF_JSON_FILE" 2>/dev/null || echo "unknown")
        response_time=$(jq -r ".[$i].response_time_ms // \"0\"" "$PERF_JSON_FILE" 2>/dev/null || echo "0")
        
        # 存储到关联数组
        fulltime_data["$filename"]="$response_time"
    done
    
    echo "已从JSON文件读取 $total_files 个文件的fulltime数据"
else
    echo "未提供JSON性能数据文件或文件不存在: $PERF_JSON_FILE"
fi

# 格式化RESPONSE_DIR目录下的所有JSON文件
echo "格式化RESPONSE_DIR目录下的JSON文件..."
format_json_count=0
for json_file in "$RESPONSE_DIR"/*.json; do
    if [ -f "$json_file" ]; then
        # 使用jq格式化JSON文件，并保存回原文件
        jq '.' "$json_file" > "$json_file.tmp" && mv "$json_file.tmp" "$json_file"
        format_json_count=$((format_json_count + 1))
    fi
done
echo "已格式化 $format_json_count 个JSON文件"

# 使用响应文件目录方式处理数据
echo "使用响应文件目录方式处理数据..."

# 获取数据目录中的所有文件
file_list=$(find "$DATA_DIR" -type f | sort)
total_files=$(echo "$file_list" | wc -l)

for filepath in $file_list; do
    # 获取文件名（不包含路径）
    filename=$(basename "$filepath")
    response_file="$RESPONSE_DIR/${filename%.*}.json"
    
    processed=$((processed + 1))
    echo "[$processed/$total_files] 处理文件: $filename"
    
    # 初始化统计值
    prompt_tokens="0"
    completion_tokens="0"
    fulltime="${fulltime_data[$filename]:-0}"
    
    # 获取文件大小（字节）
    if [ -f "$filepath" ]; then
        filesize=$(stat -c%s "$filepath" 2>/dev/null || stat -f%z "$filepath" 2>/dev/null || echo "0")
    else
        filesize="0"
    fi
    
    # 增加总测试数量
    total_count=$((total_count + 1))
    
    if [ -f "$response_file" ]; then
        # 尝试从响应文件中提取信息
        if grep -q '"choices"' "$response_file"; then
            # 增加成功计数
            total_success=$((total_success + 1))
            # 从usage中提取token信息
            prompt_tokens=$(jq -r '.usage.prompt_tokens // "0"' "$response_file" 2>/dev/null)
            completion_tokens=$(jq -r '.usage.completion_tokens // "0"' "$response_file" 2>/dev/null)
            # 使用fulltime作为耗用时间
            if [ "$fulltime" != "0" ]; then
                # 累加有效响应时间
                total_time_sum=$(echo "$total_time_sum + $fulltime" | bc 2>/dev/null || echo "$total_time_sum")
                total_time_valid=$((total_time_valid + 1))
                # 将响应时间添加到数组中，用于计算百分位数
                response_times+=("$fulltime")
            fi
            
            # 累加token数量
            total_prompt_tokens=$((total_prompt_tokens + prompt_tokens))
            total_completion_tokens=$((total_completion_tokens + completion_tokens))
            
            echo "  文件: $filename, fulltime: $fulltime, 输入Token: $prompt_tokens, 输出Token: $completion_tokens"
        else
            echo "  文件: $filename, 响应文件格式错误"
        fi
    else
        echo "  文件: $filename, 响应文件不存在"
    fi
    
    # 记录结果到CSV
    echo "$filename,$filesize,$fulltime,$prompt_tokens,$completion_tokens" >> "$result_file"
done

# 生成性能测试汇总报告
echo "生成性能测试汇总报告..."
summary_file="$OUTPUT_DIR/perf_summary.txt"

echo "LLM性能测试汇总报告" > "$summary_file"
echo "测试时间: $(date)" >> "$summary_file"
echo "数据来源: $DATA_DIR" >> "$summary_file"
echo "响应文件目录: $RESPONSE_DIR" >> "$summary_file"
echo "测试文件数量: $processed" >> "$summary_file"
echo "" >> "$summary_file"

# 总体统计
echo "总体统计:" >> "$summary_file"
echo "  总测试数量: $total_count" >> "$summary_file"
echo "  总成功数量: $total_success" >> "$summary_file"

if [ $total_count -gt 0 ]; then
    success_rate=$(echo "scale=2; $total_success * 100 / $total_count" | bc)
    echo "  成功率: ${success_rate}%" >> "$summary_file"
else
    echo "  成功率: 0%" >> "$summary_file"
fi

if [ $total_time_valid -gt 0 ]; then
    total_avg_time=$(echo "scale=2; $total_time_sum / $total_time_valid" | bc)
else
    total_avg_time="0"
fi

echo "  平均响应时间: ${total_avg_time}ms" >> "$summary_file"

# 计算并输出P90和P99百分位数
if [ ${#response_times[@]} -gt 0 ]; then
    # 使用sort命令对响应时间进行排序
    IFS=$'\n' sorted_times=($(sort -n <<<"${response_times[*]}"))
    unset IFS
    
    # 计算P90（90%的请求响应时间小于等于此值）
    p90_index=$(echo "(${#sorted_times[@]} * 90) / 100" | bc)
    if [ $p90_index -ge ${#sorted_times[@]} ]; then
        p90_index=$((${#sorted_times[@]} - 1))
    fi
    p90_value=${sorted_times[$p90_index]}
    
    # 计算P99（99%的请求响应时间小于等于此值）
    p99_index=$(echo "(${#sorted_times[@]} * 99) / 100" | bc)
    if [ $p99_index -ge ${#sorted_times[@]} ]; then
        p99_index=$((${#sorted_times[@]} - 1))
    fi
    p99_value=${sorted_times[$p99_index]}
    
    echo "  P90响应时间: ${p90_value}ms" >> "$summary_file"
    echo "  P99响应时间: ${p99_value}ms" >> "$summary_file"
else
    echo "  P90响应时间: 0" >> "$summary_file"
    echo "  P99响应时间: 0" >> "$summary_file"
fi

echo "  总输入Token: $total_prompt_tokens" >> "$summary_file"
echo "  总输出Token: $total_completion_tokens" >> "$summary_file"

echo "CSV结果文件已保存到: $result_file"
echo "性能测试汇总报告已保存到: $summary_file"
cat $summary_file
