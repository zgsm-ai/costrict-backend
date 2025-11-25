#!/bin/bash

# 模型性能对比报告生成脚本
# 扫描指定目录，找到含有perf_results.csv文件的子目录，对比不同模型的性能数据

# 默认参数值
SCAN_DIR=""
OUTPUT_DIR="./comparison-$(date +%Y%m%d%H%M%S)"
OUTPUT_FORMAT="both"

# 显示使用说明
show_usage() {
    echo "用法: $0 [选项]"
    echo "选项:"
    echo "  -d, --dir DIR       指定要扫描的目录 (必要参数)"
    echo "  -o, --output DIR    指定输出目录 (默认: $OUTPUT_DIR)"
    echo "  -f, --format FILE   指定输出格式 (csv|html|both, 默认: both)"
    echo "  -h, --help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 -d ."
    echo "  $0 -d ./tests -o ./comparison_results"
    echo "  $0 -d ./tests -f html"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--dir)
            SCAN_DIR="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -f|--format)
            OUTPUT_FORMAT="$2"
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

# 检查是否指定了扫描目录
if [ -z "$SCAN_DIR" ]; then
    echo "错误: 必须指定扫描目录"
    show_usage
    exit 1
fi

# 检查扫描目录是否存在
if [ ! -d "$SCAN_DIR" ]; then
    echo "错误: 扫描目录不存在: $SCAN_DIR"
    exit 1
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 查找所有包含perf_results.csv文件的目录
echo "正在扫描目录: $SCAN_DIR"
echo "查找包含perf_results.csv文件的子目录..."

model_dirs=()
for dir in $(find "$SCAN_DIR" -type d); do
    if [ -f "$dir/perf_results.csv" ]; then
        # 检查CSV文件的格式
        # 只处理与gen-llm-reports.sh生成的perf_results.csv相同格式的文件
        # 即格式: filename,filesize,fulltime,prompt_tokens,completion_tokens
        header=$(head -n 1 "$dir/perf_results.csv")
        if [[ "$header" == "filename,filesize,fulltime,prompt_tokens,completion_tokens" ]]; then
            model_dirs+=("$dir")
        else
            echo "跳过目录 $dir，因为perf_results.csv格式不符合要求"
        fi
    fi
done

if [ ${#model_dirs[@]} -eq 0 ]; then
    echo "错误: 未找到任何包含perf_results.csv文件的子目录"
    exit 1
fi

echo "找到 ${#model_dirs[@]} 个包含perf_results.csv文件的子目录:"
for dir in "${model_dirs[@]}"; do
    echo "  - $dir"
done

# 生成统计数据
generate_statistics() {
    local all_data_file="$OUTPUT_DIR/all_data.csv"
    local summary_file="$OUTPUT_DIR/model_summary.csv"
    
    echo "正在生成统计数据..."
    # 创建临时文件来存储所有模型的数据
    echo "model,filename,filesize,fulltime,prompt_tokens,completion_tokens" > "$all_data_file"
    
    # 提取所有模型名称
    model_names=()
    for dir in "${model_dirs[@]}"; do
        # 从目录名中提取模型名称
        model_name=$(basename "$dir")
        model_names+=("$model_name")
        
        # 读取perf_results.csv文件并添加模型名称列
        # 格式: filename,filesize,fulltime,prompt_tokens,completion_tokens
        tail -n +2 "$dir/perf_results.csv" | while IFS=, read -r filename filesize fulltime prompt_tokens completion_tokens; do
            echo "$model_name,$filename,$filesize,$fulltime,$prompt_tokens,$completion_tokens" >> "$all_data_file"
        done
    done
    
    # 创建汇总CSV文件
    echo "model,total_files,avg_fulltime,p90_fulltime,p99_fulltime,avg_prompt_tokens,avg_completion_tokens" > "$summary_file"
    
    # 对每个模型计算统计信息
    for model in "${model_names[@]}"; do
        # 获取该模型的所有数据
        model_data=$(grep "^$model," "$all_data_file")
        
        if [ -n "$model_data" ]; then
            # 总文件数
            total_files=$(echo "$model_data" | wc -l)
            
            # 计算平均fulltime
            fulltimes=$(echo "$model_data" | cut -d, -f4 | grep -v '^$')
            if [ -n "$fulltimes" ]; then
                sum_fulltime=0
                count_fulltime=0
                for ft in $fulltimes; do
                    sum_fulltime=$(echo "$sum_fulltime + $ft" | bc)
                    count_fulltime=$((count_fulltime + 1))
                done
                avg_fulltime=$(echo "scale=2; $sum_fulltime / $count_fulltime" | bc)
                
                # 计算P90和P99
                sorted_fulltimes=$(echo "$fulltimes" | sort -n)
                p90_index=$(echo "scale=0; $count_fulltime * 0.9 / 1" | bc)
                p99_index=$(echo "scale=0; $count_fulltime * 0.99 / 1" | bc)
                
                p90_fulltime=$(echo "$sorted_fulltimes" | sed -n "${p90_index}p")
                p99_fulltime=$(echo "$sorted_fulltimes" | sed -n "${p99_index}p")
            else
                avg_fulltime="0"
                p90_fulltime="0"
                p99_fulltime="0"
            fi
            
            # 计算平均token数
            prompt_tokens=$(echo "$model_data" | cut -d, -f5 | grep -v '^0$' | grep -v '^$')
            completion_tokens=$(echo "$model_data" | cut -d, -f6 | grep -v '^0$' | grep -v '^$')
            
            if [ -n "$prompt_tokens" ]; then
                sum_prompt=0
                count_prompt=0
                for pt in $prompt_tokens; do
                    sum_prompt=$(echo "$sum_prompt + $pt" | bc)
                    count_prompt=$((count_prompt + 1))
                done
                avg_prompt_tokens=$(echo "scale=2; $sum_prompt / $count_prompt" | bc)
            else
                avg_prompt_tokens="0"
            fi
            
            if [ -n "$completion_tokens" ]; then
                sum_completion=0
                count_completion=0
                for ct in $completion_tokens; do
                    sum_completion=$(echo "$sum_completion + $ct" | bc)
                    count_completion=$((count_completion + 1))
                done
                avg_completion_tokens=$(echo "scale=2; $sum_completion / $count_completion" | bc)
            else
                avg_completion_tokens="0"
            fi
            
            echo "$model,$total_files,$avg_fulltime,$p90_fulltime,$p99_fulltime,$avg_prompt_tokens,$avg_completion_tokens" >> "$summary_file"
        fi
    done
    
    echo "统计数据已生成并保存到: $summary_file"
}

# 生成CSV格式的报告
generate_csv_report() {
    local all_data_file="$OUTPUT_DIR/all_data.csv"
    local details_file="$OUTPUT_DIR/model_details.csv"
    
    echo "生成CSV格式的模型对比报告..."
    
    # 创建CSV文件
    echo "filename,${model_names[*]}" | tr ' ' ',' > "$details_file"
    
    # 获取所有唯一的文件名
    unique_files=($(tail -n +2 "$all_data_file" | cut -d, -f2 | sort -u))
    
    # 对每个文件生成一行对比数据
    for file in "${unique_files[@]}"; do
        # 获取该文件在所有模型中的数据
        row_data="$file"
        
        for model in "${model_names[@]}"; do
            # 查找该模型中该文件的数据
            model_data=$(grep "^$model,$file," "$all_data_file")
            
            if [ -n "$model_data" ]; then
                fulltime=$(echo "$model_data" | cut -d, -f4)
                prompt_tokens=$(echo "$model_data" | cut -d, -f5)
                completion_tokens=$(echo "$model_data" | cut -d, -f6)
                
                # 显示数据
                row_data="$row_data,${fulltime}ms/${prompt_tokens}/${completion_tokens}"
            else
                row_data="$row_data,0"
            fi
        done
        
        echo "$row_data" >> "$details_file"
    done
    
    echo "CSV格式的模型对比报告已保存到: $details_file"
}

# 生成HTML格式的报告
generate_html_report() {
    local all_data_file="$OUTPUT_DIR/all_data.csv"
    local summary_file="$OUTPUT_DIR/model_summary.csv"
    local html_file="$OUTPUT_DIR/model_reports.html"
    
    echo "生成HTML格式的模型对比报告..."
    
    # 创建HTML文件
    cat > "$html_file" << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>模型性能对比报告</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1, h2 {
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: 30px;
            background-color: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.2);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: center;
        }
        th {
            background-color: #4CAF50;
            color: white;
            font-weight: bold;
            position: sticky;
            top: 0;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        tr:hover {
            background-color: #f1f1f1;
        }
        .best {
            background-color: #d4edda !important;
            font-weight: bold;
        }
        .worst {
            background-color: #f8d7da !important;
        }
        .summary-table th, .summary-table td {
            text-align: left;
        }
        .details-table th {
            text-align: center;
        }
        .details-table td {
            text-align: center;
        }
        .details-table .filename {
            text-align: left;
            font-weight: bold;
        }
        .metric-value {
            font-weight: bold;
        }
        .token-info {
            font-size: 0.9em;
            color: #666;
        }
        .footer {
            margin-top: 30px;
            font-size: 0.9em;
            color: #666;
            text-align: center;
        }
        .legend {
            margin-bottom: 20px;
            padding: 10px;
            background-color: white;
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.2);
        }
        .legend-item {
            display: inline-block;
            margin-right: 15px;
        }
        .legend-color {
            display: inline-block;
            width: 15px;
            height: 15px;
            margin-right: 5px;
            border: 1px solid #ddd;
        }
        .best-color {
            background-color: #d4edda;
        }
        .worst-color {
            background-color: #f8d7da;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>模型性能对比报告</h1>
        <p>生成时间: $(date)</p>
        <p>扫描目录: $SCAN_DIR</p>
        
        <div class="legend">
            <div class="legend-item">
                <span class="legend-color best-color"></span>最佳性能
            </div>
            <div class="legend-item">
                <span class="legend-color worst-color"></span>最差性能
            </div>
        </div>
        
        <h2>模型汇总表格</h2>
EOF

    # 添加汇总表格
    cat >> "$html_file" << EOF
        <table class="summary-table">
            <thead>
                <tr>
                    <th>模型名称</th>
                    <th>总文件数</th>
                    <th>平均响应时间 (ms)</th>
                    <th>P90响应时间 (ms)</th>
                    <th>P99响应时间 (ms)</th>
                    <th>平均输入Token</th>
                    <th>平均输出Token</th>
                </tr>
            </thead>
            <tbody>
EOF

    # 添加汇总数据行，并找出最佳和最差性能
    declare -a avg_times
    declare -a model_index_map
    index=0
    
    # 先读取所有数据以确定最佳和最差性能
    while IFS=, read -r model total_files avg_fulltime p90_fulltime p99_fulltime avg_prompt_tokens avg_completion_tokens; do
        avg_times[$index]="$avg_fulltime"
        model_index_map[$index]="$model"
        index=$((index + 1))
    done < <(tail -n +2 "$summary_file")
    
    # 找出最佳和最差的平均响应时间
    best_time=$(echo "${avg_times[@]}" | tr ' ' '\n' | sort -n | head -n1)
    worst_time=$(echo "${avg_times[@]}" | tr ' ' '\n' | sort -n | tail -n1)
    
    # 重新遍历并添加行，同时标记最佳和最差性能
    index=0
    while IFS=, read -r model total_files avg_fulltime p90_fulltime p99_fulltime avg_prompt_tokens avg_completion_tokens; do
        row_class=""
        if (( $(echo "$avg_fulltime == $best_time" | bc -l) )); then
            row_class="best"
        elif (( $(echo "$avg_fulltime == $worst_time" | bc -l) )); then
            row_class="worst"
        fi
        
        cat >> "$html_file" << EOF
                <tr class="$row_class">
                    <td>$model</td>
                    <td>$total_files</td>
                    <td class="metric-value">$avg_fulltime</td>
                    <td>$p90_fulltime</td>
                    <td>$p99_fulltime</td>
                    <td>$avg_prompt_tokens</td>
                    <td>$avg_completion_tokens</td>
                </tr>
EOF
        index=$((index + 1))
    done < <(tail -n +2 "$summary_file")

    # 结束汇总表格
    cat >> "$html_file" << EOF
            </tbody>
        </table>
        
        <h2>详细对比表格</h2>
        <table class="details-table">
            <thead>
                <tr>
                    <th>文件名</th>
EOF

    # 添加详细表格的表头（模型名称）
    for model in "${model_names[@]}"; do
        echo "                    <th>$model</th>" >> "$html_file"
    done

    echo "                </tr>" >> "$html_file"
    echo "            </thead>" >> "$html_file"
    echo "            <tbody>" >> "$html_file"

    # 获取所有唯一的文件名
    unique_files=($(tail -n +2 "$all_data_file" | cut -d, -f2 | sort -u))

    # 对每个文件生成一行对比数据
    for file in "${unique_files[@]}"; do
        echo "                <tr>" >> "$html_file"
        echo "                    <td class=\"filename\">$file</td>" >> "$html_file"
        
        # 收集所有模型在该文件上的响应时间，以确定最佳和最差性能
        declare -a file_times
        declare -a file_models
        
        for model in "${model_names[@]}"; do
            # 查找该模型中该文件的数据
            model_data=$(grep "^$model,$file," "$all_data_file")
            
            if [ -n "$model_data" ]; then
                fulltime=$(echo "$model_data" | cut -d, -f4)
                file_times+=("$fulltime")
                file_models+=("$model")
            fi
        done
        
        # 找出最佳和最差的响应时间
        if [ ${#file_times[@]} -gt 0 ]; then
            best_file_time=$(echo "${file_times[@]}" | tr ' ' '\n' | sort -n | head -n1)
            worst_file_time=$(echo "${file_times[@]}" | tr ' ' '\n' | sort -n | tail -n1)
        fi
        
        # 再次遍历所有模型，生成单元格并标记最佳和最差性能
        for model in "${model_names[@]}"; do
            # 查找该模型中该文件的数据
            model_data=$(grep "^$model,$file," "$all_data_file")
            
            if [ -n "$model_data" ]; then
                fulltime=$(echo "$model_data" | cut -d, -f4)
                prompt_tokens=$(echo "$model_data" | cut -d, -f5)
                completion_tokens=$(echo "$model_data" | cut -d, -f6)
                
                cell_class=""
                if (( $(echo "$fulltime == $best_file_time" | bc -l) )); then
                    cell_class="best"
                elif (( $(echo "$fulltime == $worst_file_time" | bc -l) )); then
                    cell_class="worst"
                fi
                
                echo "                    <td class=\"$cell_class\">${fulltime}ms<div class=\"token-info\">${prompt_tokens}/${completion_tokens}</div></td>" >> "$html_file"
            else
                echo "                    <td>N/A</td>" >> "$html_file"
            fi
        done
        
        echo "                </tr>" >> "$html_file"
    done

    # 结束HTML文件
    cat >> "$html_file" << EOF
            </tbody>
        </table>
        
        <div class="footer">
            <p>模型对比报告生成完成 | 生成时间: $(date)</p>
        </div>
    </div>
</body>
</html>
EOF

    echo "HTML格式的模型对比报告已保存到: $html_file"
}

# 首先生成统计数据
generate_statistics

# 根据输出格式生成报告
if [ "$OUTPUT_FORMAT" = "csv" ] || [ "$OUTPUT_FORMAT" = "both" ]; then
    generate_csv_report
fi

if [ "$OUTPUT_FORMAT" = "html" ] || [ "$OUTPUT_FORMAT" = "both" ]; then
    generate_html_report
fi

echo "模型对比报告生成完成!"
