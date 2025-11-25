#!/bin/bash

# Stress test script for completion API
# This script generates multiple concurrent requests to test the system performance

# ========================================
# Function Definitions
# ========================================

# Generate random 16-byte hex string
function generate_random_id() {
  openssl rand -hex 16 2>/dev/null || (
    # Fallback if openssl is not available
    date +%s%N | sha256sum | head -c 32
  )
}

# Function to start expose-completion.sh
function start_expose_completion() {
  if [ "$COMPLETION_EXPOSE" = true ]; then
    # Check if expose-completion.sh exists
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
}

# Function to stop expose-completion.sh
function stop_expose_completion() {
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
}

# Function to generate test report
function generate_test_report() {
  local report_file="$OUTPUT_DIR/stress_summary.txt"
  echo "生成测试报告: $report_file"

  cat > "$report_file" << EOF
压力测试报告
=============

测试时间: $(date)
测试持续时间: $ELAPSED_TIME 秒
并发任务数: $CONCURRENT_JOBS
批次超时时间: $BATCH_TIMEOUT 秒
总批次数: $TOTAL_BATCHES
测试文件数: ${#TEST_FILES[@]}
服务器地址: $COMPLETION_URL

ClientID列表:
EOF

  for i in "${!CLIENT_IDS[@]}"; do
    echo "ClientID $((i+1)): ${CLIENT_IDS[i]}" >> "$report_file"
  done

  # 调用gen-completion-reports.sh生成CSV文件和性能报告
  echo "使用gen-completion-reports.sh生成CSV文件和性能报告..."
  if [ -f "./gen-completion-reports.sh" ]; then
      ./gen-completion-reports.sh -d "$DATA_DIR" -r "$OUTPUT_DIR" -o "$OUTPUT_DIR"
      echo "性能报告生成完成"
      
      # 显示生成的报告文件
      echo ""
      echo "生成的报告文件:"
      if [ -f "$OUTPUT_DIR/perf_results.csv" ]; then
          echo "- CSV结果文件: $OUTPUT_DIR/perf_results.csv"
          echo "前5行内容:"
          head -n 5 "$OUTPUT_DIR/perf_results.csv" | while IFS=, read -r filename status fulltime prompt_tokens completion_tokens total_duration; do
              echo "  $filename | $status | $fulltime | $prompt_tokens | $completion_tokens | $total_duration"
          done
      fi
  else
      echo "警告: 找不到gen-completion-reports.sh脚本，跳过报告生成"
  fi
}

# Function to get language ID from file extension
function get_language_id() {
  local file_path="$1"
  local ext="${file_path##*.}"
  local language_id=""
  
  case "$ext" in
    py) language_id="python" ;;
    c) language_id="c" ;;
    cpp) language_id="cpp" ;;
    go) language_id="go" ;;
    java) language_id="java" ;;
    js) language_id="javascript" ;;
    ts) language_id="typescript" ;;
    lua) language_id="lua" ;;
    *) language_id="text" ;;
  esac
  
  echo "$language_id"
}

# Print usage information
function print_help() {
  echo "Usage: $0 [options]"
  echo "  --no-expose      在执行测试前不调用 expose-completion.sh 暴露服务端口 (默认: false)"
  echo "  -u url: 服务器地址 (必需)"
  echo "  -d dir: 测试文件目录 (默认: data)"
  echo "  -j jobs: 并发任务数 (默认: 10)"
  echo "  -c clients: 客户端ID数量 (默认: 10)"
  echo "  -k apikey: API密钥"
  echo "  -m model: 模型名称"
  echo "  -o dir: 输出目录 (默认: stress_xxx)"
  echo "  -t temperature: 温度值 (默认: 0.1)"
  echo "  -M max_tokens: 最大令牌数 (默认: 50)"
  echo "  -T timeout: 批次超时时间（秒）(默认: 30)"
  echo "  -h: 帮助信息"
}

# Function to execute a single request
function execute_request() {
  local file="$1"
  local index="$2"
  local client_id="${CLIENT_IDS[$(($RANDOM % $CONCURRENT_CLIENTS))]}"
  local completion_id="comp-$(generate_random_id)"
  local fname=$(basename "$file")
  local output_file="$OUTPUT_DIR/${fname%.*}.json"

  # Get language ID from file extension
  local language_id=$(get_language_id "$file")
  
  # Prepare command arguments
  local cmd_args="-a $COMPLETION_URL -f $file -c $client_id -i $completion_id -l $language_id -t $TEMPERATURE -M $MAX_TOKENS -o $output_file -n"
  
  if [ -n "$COMPLETION_APIKEY" ]; then
    cmd_args="$cmd_args -k $COMPLETION_APIKEY"
  fi
  
  if [ -n "$COMPLETION_MODEL" ]; then
    cmd_args="$cmd_args -m $COMPLETION_MODEL"
  fi
  
  # Execute the request
  echo "执行请求 $index: 文件=$file, ClientID=$client_id, CompletionID=$completion_id"
  ./completion-via-service.sh $cmd_args
  
  # Check if request was successful
  if [ $? -eq 0 ]; then
    echo "请求 $index 完成: 输出文件=$output_file"
  else
    echo "请求 $index 失败: 文件=$file"
  fi
}

# Function to execute concurrent tests
function execute_concurrent_tests() {
  local test_files=("$@")
  local concurrent_jobs="$CONCURRENT_JOBS"
  local batch_timeout="$BATCH_TIMEOUT"
  local output_dir="$OUTPUT_DIR"
  
  echo "开始批次测试，每批并发数: $concurrent_jobs"
  echo "批次超时时间: $batch_timeout 秒"
  echo "测试开始时间: $(date)"
  
  # 计算总批次数
  local total_files=${#test_files[@]}
  local total_batches=$(( (total_files + concurrent_jobs - 1) / concurrent_jobs ))
  echo "总共 $total_files 个文件，将分为 $total_batches 批执行"
  
  # 按批次执行任务
  for (( batch=0; batch<total_batches; batch++ )); do
    local batch_start=$(( batch * concurrent_jobs ))
    local batch_end=$(( batch_start + concurrent_jobs - 1 ))
    if [ $batch_end -ge $total_files ]; then
      batch_end=$((total_files - 1))
    fi
    
    echo "执行第 $((batch+1))/$total_batches 批任务 (文件 $((batch_start+1))-$((batch_end+1)))"
    
    local pids=()
    local batch_start_time=$(date +%s)
    
    # 启动当前批次的任务
    for (( i=batch_start; i<=batch_end; i++ )); do
      (
        local index=$((i+1))
        local file="${test_files[$i]}"
        execute_request "$file" "$index"
      ) &
      pids+=($!)
    done
    
    # 等待当前批次的所有任务完成或超时
    echo "等待第 $((batch+1)) 批的 ${#pids[@]} 个任务完成..."
    local batch_timeout_reached=false
    
    for pid in "${pids[@]}"; do
      local current_time=$(date +%s)
      local elapsed_batch_time=$((current_time - batch_start_time))
      
      if [ $elapsed_batch_time -ge $batch_timeout ]; then
        echo "警告: 第 $((batch+1)) 批执行时间超过 $batch_timeout 秒，强制终止剩余任务"
        for term_pid in "${pids[@]}"; do
          if kill -0 $term_pid 2>/dev/null; then
            echo "终止进程: $term_pid"
            kill -9 $term_pid 2>/dev/null
          fi
        done
        batch_timeout_reached=true
        break
      fi
      
      # 检查进程是否仍在运行
      if kill -0 $pid 2>/dev/null; then
        echo "等待进程 $pid 完成..."
        wait $pid
      fi
    done
    
    if [ "$batch_timeout_reached" = false ]; then
      echo "第 $((batch+1)) 批任务已完成"
    else
      echo "第 $((batch+1)) 批任务因超时提前结束"
    fi
  done
  
  echo "所有批次的任务已执行完成"
}

# Export functions for use in subshells
export -f generate_random_id
export -f start_expose_completion
export -f stop_expose_completion
export -f generate_test_report
export -f get_language_id
export -f print_help
export -f execute_request
export -f execute_concurrent_tests

# ========================================
# Main Script
# ========================================

# Default values
COMPLETION_EXPOSE=true
COMPLETION_URL=""
COMPLETION_APIKEY=""
COMPLETION_MODEL=""

DATA_DIR="./data"
CONCURRENT_JOBS=10
CONCURRENT_CLIENTS=10
OUTPUT_DIR="stress_$(date +%Y%m%d_%H%M%S)"
TEMPERATURE="0.1"
MAX_TOKENS="50"
BATCH_TIMEOUT=30

# 如果存在./.env文件，则加载它
if [ -f "./.env" ]; then
    echo "加载./.env文件中的环境变量..."
    source ./.env
fi

# Parse command line options
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-expose)
            COMPLETION_EXPOSE=false
            shift
            ;;
        -u)
            COMPLETION_URL="$2"
            shift 2
            ;;
        -d)
            DATA_DIR="$2"
            shift 2
            ;;
        -j)
            CONCURRENT_JOBS="$2"
            shift 2
            ;;
        -c)
            CONCURRENT_CLIENTS="$2"
            shift 2
            ;;
        -k)
            COMPLETION_APIKEY="$2"
            shift 2
            ;;
        -m)
            COMPLETION_MODEL="$2"
            shift 2
            ;;
        -o)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -t)
            TEMPERATURE="$2"
            shift 2
            ;;
        -M)
            MAX_TOKENS="$2"
            shift 2
            ;;
        -T)
            BATCH_TIMEOUT="$2"
            shift 2
            ;;
        -h)
            print_help
            exit 0
            ;;
        *)
            echo "无效选项: $1"
            print_help
            exit 1
            ;;
    esac
done

# Validate required parameters
if [ -z "$COMPLETION_URL" ]; then
  echo "错误: 缺少必需参数 '-u/--url'"
  print_help
  exit 1
fi

# Check if data directory exists
if [ ! -d "$DATA_DIR" ]; then
  echo "错误: 数据目录 '$DATA_DIR' 不存在"
  exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Generate ClientIDs based on CONCURRENT_CLIENTS
echo "生成ClientID列表..."
CLIENT_IDS=()
for i in $(seq 1 $CONCURRENT_CLIENTS); do
  CLIENT_ID=$(generate_random_id)
  CLIENT_IDS+=("cli-$CLIENT_ID")
  echo "ClientID $i: $CLIENT_ID"
done

# Get list of test files
echo "扫描测试文件..."
TEST_FILES=($(find "$DATA_DIR" -type f \( -name "*.py" -o -name "*.c" -o -name "*.cpp" -o -name "*.go" -o -name "*.java" -o -name "*.js" -o -name "*.ts" -o -name "*.lua" \) | sort))

if [ ${#TEST_FILES[@]} -eq 0 ]; then
  echo "错误: 在 '$DATA_DIR' 目录中未找到测试文件"
  exit 1
fi

echo "找到 ${#TEST_FILES[@]} 个测试文件"

# 启动 expose-completion.sh（如果需要）
start_expose_completion

# Start time for performance measurement
START_TIME=$(date +%s)

# Export variables for subshells
export CLIENT_IDS
export OUTPUT_DIR
export COMPLETION_URL
export COMPLETION_APIKEY
export COMPLETION_MODEL
export TEMPERATURE
export MAX_TOKENS

# Execute requests in batches
execute_concurrent_tests "${TEST_FILES[@]}"

# End time for performance measurement
END_TIME=$(date +%s)
ELAPSED_TIME=$((END_TIME - START_TIME))

echo "测试完成时间: $(date)"
echo "总耗时: $ELAPSED_TIME 秒"
echo "结果文件保存在: $OUTPUT_DIR"

# 停止 expose-completion.sh（如果已经启动）
stop_expose_completion

# 生成测试报告
generate_test_report