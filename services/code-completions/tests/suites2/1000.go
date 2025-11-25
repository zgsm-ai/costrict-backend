package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type DatabasePool struct {
	config *DatabaseConfig
	db     *sql.DB
}

func NewDatabasePool(config *DatabaseConfig) (*DatabasePool, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username, config.Password,
		config.Host, config.Port, config.Database)

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("连接失败: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("测试连接失败: %v", err)
	}

	return &DatabasePool{config: config, db: db}, nil
}

func (dp *DatabasePool) GetDB() *sql.DB {
	return dp.db
}

func (dp *DatabasePool) Close() error {
	return dp.db.Close()
}

func (dp *DatabasePool) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return dp.db.QueryContext(ctx, query, args...)
}

func (dp *DatabasePool) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := dp.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (dp *DatabasePool) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := dp.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("执行失败: %v, 回滚失败: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func main() {
	config := &DatabaseConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "testdb",
	}

	pool, err := NewDatabasePool(config)
	if err != nil {
		log.Fatalf("创建连接池失败: %v", err)
	}
	defer pool.Close()

	fmt.Println("数据库连接池创建成功")

	// 执行简单查询
	rows, err := pool.Query(context.Background(), "SELECT 1")
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		rows.Close()
		fmt.Println("简单查询执行成功")
	}

	// 执行事务
	err = pool.Transaction(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.ExecContext(context.Background(),
			"CREATE TABLE IF NOT EXISTS test (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))")
		if err != nil {
			return fmt.Errorf("创建表失败: %v", err)
		}

		_, err = tx.ExecContext(context.Background(),
			"INSERT INTO test (name) VALUES (?)", "测试数据")
		if err != nil {
			return fmt.Errorf("插入数据失败: %v", err)
		}
		return nil
	})

	if err != nil {
		log.Printf("执行事务失败: %v", err)
	} else {
		fmt.Println("事务执行成功")
	}<｜fim▁hole｜>

	fmt.Println("数据库操作演示完成")
}
