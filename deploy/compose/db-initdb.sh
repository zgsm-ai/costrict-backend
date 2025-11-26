#!/bin/sh

. ./utils.sh
. ./configure.sh

DEPLOY="$(basename "$(dirname "$(readlink -f "$0")")")"

retry "docker compose -f docker-compose-initdb.yml up -d" 60 5 || fatal "Failed to start postgres"

sleep 5

# 等待PostgreSQL完全启动
echo "等待PostgreSQL完全启动..."
for i in {1..30}; do
    if docker exec "${DEPLOY}-postgres-1" pg_isready -U $POSTGRES_USER -q; then
        echo "PostgreSQL已准备就绪"
        break
    else
        echo "等待PostgreSQL启动... ($i/30)"
        sleep 2
    fi
done

# 先执行所有创建表的SQL文件（create）
for sql_file in postgres/scripts/*create*.sql; do
    # 短路返回：如果不是文件则跳过
    [ ! -f "$sql_file" ] && continue
    
    # 获取文件名（不包含路径）
    filename=$(basename "$sql_file")
    # 去掉.sql后缀得到数据库名前缀
    db=$(basename "$sql_file" .sql)
    # 从create文件名中提取实际数据库名（处理如casdoor-create1.sql -> casdoor的情况）
    db=$(echo "$db" | sed -E 's/-create[0-9]*$//')
    echo "正在执行表创建SQL文件: $filename，目标数据库: $db"
    
    # 执行SQL文件的命令，使用提取的数据库名
    docker exec -i "${DEPLOY}-postgres-1" /usr/local/bin/psql -U $POSTGRES_USER -d $db -f "/scripts/$filename"

    # 执行完成，等待5s
    echo "执行表创建SQL文件: $filename 完成, 等待5秒..."
    sleep 5
done

# 然后执行所有修改表的SQL文件（alter）
for sql_file in postgres/scripts/*alter*.sql; do
    # 短路返回：如果不是文件则跳过
    [ ! -f "$sql_file" ] && continue
    
    # 获取文件名（不包含路径）
    filename=$(basename "$sql_file")
    # 去掉.sql后缀得到数据库名前缀
    db=$(basename "$sql_file" .sql)
    # 从alter文件名中提取实际数据库名（处理如casdoor-alter1.sql -> casdoor的情况）
    db=$(echo "$db" | sed -E 's/-alter[0-9]*$//')
    echo "正在执行表修改SQL文件: $filename，目标数据库: $db"
    
    # 执行SQL文件的命令，使用提取的数据库名
    docker exec -i "${DEPLOY}-postgres-1" /usr/local/bin/psql -U $POSTGRES_USER -d $db -f "/scripts/$filename"

    # 执行完成，等待5s
    echo "执行表修改SQL文件: $filename 完成, 等待5秒..."
    sleep 5
done

# 最后执行所有插入数据的SQL文件（insert）
for sql_file in postgres/scripts/*insert*.sql; do
    # 短路返回：如果不是文件则跳过
    [ ! -f "$sql_file" ] && continue
    
    # 获取文件名（不包含路径）
    filename=$(basename "$sql_file")
    # 去掉.sql后缀得到数据库名前缀
    db=$(basename "$sql_file" .sql)
    # 从insert文件名中提取实际数据库名（处理如casdoor-insert1.sql -> casdoor的情况）
    db=$(echo "$db" | sed -E 's/-insert[0-9]*$//')
    echo "正在执行数据插入SQL文件: $filename，目标数据库: $db"
    
    # 执行SQL文件的命令，使用提取的数据库名
    docker exec -i "${DEPLOY}-postgres-1" /usr/local/bin/psql -U $POSTGRES_USER -d $db -f "/scripts/$filename"

    # 执行完成，等待5s
    echo "执行数据插入SQL文件: $filename 完成, 等待5秒..."
    sleep 5
done

SQL="SELECT datname AS database_name FROM pg_database ORDER BY datname;"
docker exec -i "${DEPLOY}-postgres-1" /usr/local/bin/psql -U $POSTGRES_USER -c "$SQL"

# 定义需要创建的数据库列表
DBLIST="chatgpt auth oneapi quota_manager codereview casdoor codebase_indexer"
for db in $DBLIST; do
    SQL="SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename;"
    echo "--------------------------------------------------------------------------"
    echo "db: ${db}"
    echo "$SQL"
    echo "--------------------------------------------------------------------------"
    docker exec -i "${DEPLOY}-postgres-1" /usr/local/bin/psql -U $POSTGRES_USER -d $db -c "$SQL"
done

retry "docker compose -f docker-compose-initdb.yml down" 60 3 || fatal "Failed to stop postgres"
