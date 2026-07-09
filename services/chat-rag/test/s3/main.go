package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ---------------------------------------------------------------------------
// CLI 参数
// ---------------------------------------------------------------------------

type Config struct {
	Endpoint      string
	Bucket        string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	Region        string
	SkipSSLVerify bool
}

func parseFlags() Config {
	cfg := Config{}

	// 从命令行参数读取，使用 --key=value 格式
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--endpoint" && i+1 < len(args):
			cfg.Endpoint = args[i+1]
			i++
		case args[i] == "--bucket" && i+1 < len(args):
			cfg.Bucket = args[i+1]
			i++
		case args[i] == "--access-key" && i+1 < len(args):
			cfg.AccessKey = args[i+1]
			i++
		case args[i] == "--secret-key" && i+1 < len(args):
			cfg.SecretKey = args[i+1]
			i++
		case args[i] == "--region" && i+1 < len(args):
			cfg.Region = args[i+1]
			i++
		case args[i] == "--use-ssl" && i+1 < len(args):
			cfg.UseSSL = args[i+1] == "true"
			i++
		case args[i] == "--skip-ssl-verify" && i+1 < len(args):
			cfg.SkipSSLVerify = args[i+1] == "true"
			i++
		case args[i] == "--help":
			printUsage()
			os.Exit(0)
		default:
			fmt.Printf("未知参数: %s\n\n", args[i])
			printUsage()
			os.Exit(1)
		}
	}

	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		fmt.Println("错误: endpoint, bucket, access-key, secret-key 为必填参数\n")
		printUsage()
		os.Exit(1)
	}

	return cfg
}

func printUsage() {
	fmt.Println(`S3 存储测试工具
用法:
  go run main.go \
    --endpoint=<host:port> \
    --bucket=<bucket> \
    --access-key=<key> \
    --secret-key=<secret> \
    --use-ssl=<true|false> \
    --region=<region> \
    --skip-ssl-verify=<true|false>

必填参数:
  --endpoint      S3 服务地址 (例如: cos-hk-xc-di1.sit.cmft.com)
  --bucket        存储桶名称
  --access-key    访问密钥 ID
  --secret-key    秘密访问密钥

可选参数:
  --use-ssl         是否使用 HTTPS (默认: true)
  --region          区域 (默认: "")
  --skip-ssl-verify 是否跳过 SSL 证书校验 (默认: false)
  --help            显示此帮助信息
`)
}

// ---------------------------------------------------------------------------
// 辅助函数
// ---------------------------------------------------------------------------

// newClient 根据配置创建 minio 客户端
func newClient(cfg Config) (*minio.Client, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	}

	if cfg.UseSSL && cfg.SkipSSLVerify {
		opts.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return minio.New(cfg.Endpoint, opts)
}

// printSeparator 打印分隔线
func printSeparator(title string) {
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Printf("  %s\n", title)
	fmt.Println("============================================================")
}

// ---------------------------------------------------------------------------
// 测试步骤
// ---------------------------------------------------------------------------

// step1_GetBucketInfo 获取桶信息
func step1_GetBucketInfo(ctx context.Context, client *minio.Client, bucket string) {
	printSeparator("步骤 1: 获取桶信息 (ListBuckets / HeadBucket)")

	fmt.Printf("[S3 API] HeadBucket (HEAD /%s)\n", bucket)
	fmt.Printf("[S3 API] ListBuckets (GET /)\n")

	// 1.1 检查桶是否存在
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		fmt.Printf("[出错] BucketExists 调用失败: %v\n", err)
		return
	}
	fmt.Printf("[成功] 桶 %q 存在: %v\n", bucket, exists)

	// 1.2 列出所有桶
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		fmt.Printf("[出错] ListBuckets 调用失败: %v\n", err)
		return
	}
	fmt.Printf("[成功] 共 %d 个桶:\n", len(buckets))
	for _, b := range buckets {
		fmt.Printf("       - %s (创建时间: %s)\n", b.Name, b.CreationDate.Format(time.RFC3339))
	}
}

// step2_UploadFile 上传文件，测试各种字符串兼容性
func step2_UploadFile(ctx context.Context, client *minio.Client, bucket string) {
	printSeparator("步骤 2: 上传文件 (PutObject) — 测试各种字符串兼容性")

	fmt.Printf("[S3 API] PutObject (PUT /%s/{object})\n", bucket)

	type testCase struct {
		name     string
		key      string
		data     []byte
	}

	testCases := []testCase{
		{
			name: "普通 ASCII 文件名",
			key:  "test/ascii_file.json",
			data: []byte(`{"hello":"world"}`),
		},
		{
			name: "中文文件名",
			key:  "test/中文文件名测试.txt",
			data: []byte("这是一个中文内容测试文件"),
		},
		{
			name: "特殊字符文件名",
			key:  "test/special_chars_!@#$%^&().txt",
			data: []byte("special characters test"),
		},
		{
			name: "带空格文件名",
			key:  "test/file name with spaces.txt",
			data: []byte("file name with spaces"),
		},
		{
			name: "Unicode 混合文件名",
			key:  "test/Unicode_混合_测试_🎉.txt",
			data: []byte("Unicode混合测试 🎉✨"),
		},
		{
			name: "空数据文件",
			key:  "test/empty_file.txt",
			data: []byte{},
		},
		{
			name: "二进制数据文件",
			key:  "test/binary_data.bin",
			data: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0x7F, 0x80},
		},
		{
			name: "大文件测试 (1MB)",
			key:  "test/large_file_1mb.bin",
			data: generateLargeData(1024 * 1024), // 1MB
		},
		{
			name: "深层路径嵌套",
			key:  "test/a/b/c/d/e/f/g/deep_nested_file.txt",
			data: []byte("deep nested path test"),
		},
		{
			name: "纯数字文件名",
			key:  "test/1234567890.json",
			data: []byte(`{"id":1234567890}`),
		},
	}

	for _, tc := range testCases {
		fmt.Printf("\n--- 测试: %s ---\n", tc.name)
		fmt.Printf("    对象键: %q\n", tc.key)
		fmt.Printf("    数据大小: %d bytes\n", len(tc.data))

		reader := bytes.NewReader(tc.data)
		info, err := client.PutObject(ctx, bucket, tc.key, reader, int64(len(tc.data)), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
		if err != nil {
			fmt.Printf("  [出错] PutObject 失败: %v\n", err)
			continue
		}
		fmt.Printf("  [成功] 上传完成, ETag: %s, 大小: %d bytes\n", info.ETag, info.Size)
	}
}

// generateLargeData 生成指定大小的随机数据
func generateLargeData(size int) []byte {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

// step3_DownloadAndVerify 下载文件并对比内容一致性
func step3_DownloadAndVerify(ctx context.Context, client *minio.Client, bucket string) {
	printSeparator("步骤 3: 下载文件并对比内容一致性 (GetObject)")

	fmt.Printf("[S3 API] GetObject (GET /%s/{object})\n", bucket)

	// 定义要验证的文件列表 (key 和期望的内容)
	type verifyItem struct {
		name string
		key  string
		data []byte
	}

	items := []verifyItem{
		{
			name: "普通 ASCII 文件",
			key:  "test/ascii_file.json",
			data: []byte(`{"hello":"world"}`),
		},
		{
			name: "中文文件",
			key:  "test/中文文件名测试.txt",
			data: []byte("这是一个中文内容测试文件"),
		},
		{
			name: "特殊字符文件",
			key:  "test/special_chars_!@#$%^&().txt",
			data: []byte("special characters test"),
		},
		{
			name: "带空格文件",
			key:  "test/file name with spaces.txt",
			data: []byte("file name with spaces"),
		},
		{
			name: "Unicode 混合文件",
			key:  "test/Unicode_混合_测试_🎉.txt",
			data: []byte("Unicode混合测试 🎉✨"),
		},
		{
			name: "空文件",
			key:  "test/empty_file.txt",
			data: []byte{},
		},
		{
			name: "二进制文件",
			key:  "test/binary_data.bin",
			data: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0x7F, 0x80},
		},
		{
			name: "深层嵌套文件",
			key:  "test/a/b/c/d/e/f/g/deep_nested_file.txt",
			data: []byte("deep nested path test"),
		},
		{
			name: "纯数字文件",
			key:  "test/1234567890.json",
			data: []byte(`{"id":1234567890}`),
		},
		{
			name: "大文件 (前64KB校验)",
			key:  "test/large_file_1mb.bin",
			data: generateLargeData(1024 * 1024),
		},
	}

	for _, item := range items {
		fmt.Printf("\n--- 验证: %s ---\n", item.name)
		fmt.Printf("    对象键: %q\n", item.key)

		obj, err := client.GetObject(ctx, bucket, item.key, minio.GetObjectOptions{})
		if err != nil {
			fmt.Printf("  [出错] GetObject 调用失败: %v\n", err)
			continue
		}
		defer obj.Close()

		got, err := io.ReadAll(obj)
		if err != nil {
			fmt.Printf("  [出错] 读取对象内容失败: %v\n", err)
			continue
		}

		// 比对内容
		if len(got) != len(item.data) {
			fmt.Printf("  [内容不一致] 大小不匹配: 期望 %d bytes, 实际 %d bytes\n", len(item.data), len(got))
			continue
		}

		if !bytes.Equal(got, item.data) {
			// 计算 MD5 方便对比
			gotMD5 := md5.Sum(got)
			wantMD5 := md5.Sum(item.data)
			fmt.Printf("  [内容不一致] MD5 不匹配\n")
			fmt.Printf("    期望 MD5: %x\n", wantMD5)
			fmt.Printf("    实际 MD5: %x\n", gotMD5)
			fmt.Printf("    期望大小: %d bytes\n", len(item.data))
			fmt.Printf("    实际大小: %d bytes\n", len(got))

			// 打印前 100 字节对比
			showLen := len(item.data)
			if showLen > 100 {
				showLen = 100
			}
			fmt.Printf("    期望前 %d 字节: %x\n", showLen, item.data[:showLen])
			fmt.Printf("    实际前 %d 字节: %x\n", showLen, got[:showLen])
			continue
		}

		fmt.Printf("  [成功] 内容一致! (大小: %d bytes)\n", len(got))
	}
}

// step4_DeleteFile 删除文件
func step4_DeleteFile(ctx context.Context, client *minio.Client, bucket string) {
	printSeparator("步骤 4: 删除文件 (RemoveObject)")

	fmt.Printf("[S3 API] RemoveObject (DELETE /%s/{object})\n", bucket)

	// 要删除的文件列表（不包含大文件，避免耗时）
	keys := []struct {
		name string
		key  string
	}{
		{name: "普通 ASCII 文件", key: "test/ascii_file.json"},
		{name: "中文文件", key: "test/中文文件名测试.txt"},
		{name: "特殊字符文件", key: "test/special_chars_!@#$%^&().txt"},
		{name: "带空格文件", key: "test/file name with spaces.txt"},
		{name: "Unicode 混合文件", key: "test/Unicode_混合_测试_🎉.txt"},
		{name: "空文件", key: "test/empty_file.txt"},
		{name: "二进制文件", key: "test/binary_data.bin"},
		{name: "大文件", key: "test/large_file_1mb.bin"},
		{name: "深层嵌套文件", key: "test/a/b/c/d/e/f/g/deep_nested_file.txt"},
		{name: "纯数字文件", key: "test/1234567890.json"},
	}

	for _, k := range keys {
		fmt.Printf("\n--- 删除: %s ---\n", k.name)
		fmt.Printf("    对象键: %q\n", k.key)

		err := client.RemoveObject(ctx, bucket, k.key, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("  [出错] RemoveObject 失败: %v\n", err)
			continue
		}
		fmt.Printf("  [成功] 已删除\n")

		// 验证确认已删除
		_, err = client.StatObject(ctx, bucket, k.key, minio.StatObjectOptions{})
		if err != nil {
			fmt.Printf("  [验证] 删除确认: 对象已不存在 (%v)\n", err)
		} else {
			fmt.Printf("  [警告] 删除后对象仍然可访问!\n")
		}
	}
}

// step5_ListFiles 获取某个目录文件列表
func step5_ListFiles(ctx context.Context, client *minio.Client, bucket string) {
	printSeparator("步骤 5: 获取文件列表 (ListObjects — 前缀查询)")

	fmt.Printf("[S3 API] ListObjects (GET /%s?prefix={prefix}&delimiter={delimiter})\n", bucket)

	// 先上传一些测试文件到不同"目录"
	testFiles := []struct {
		name string
		key  string
	}{
		{name: "目录A-文件1", key: "dir_a/file_1.txt"},
		{name: "目录A-文件2", key: "dir_a/file_2.txt"},
		{name: "目录A-子目录文件", key: "dir_a/sub/file_sub.txt"},
		{name: "目录B-文件1", key: "dir_b/file_1.txt"},
		{name: "目录B-文件2", key: "dir_b/file_2.txt"},
		{name: "目录C-文件1", key: "dir_c/file_1.txt"},
		{name: "根目录文件1", key: "root_file_1.txt"},
		{name: "根目录文件2", key: "root_file_2.txt"},
	}

	fmt.Println("\n--- 准备测试文件 ---")
	for _, tf := range testFiles {
		_, err := client.PutObject(ctx, bucket, tf.key, bytes.NewReader([]byte(tf.name)), int64(len(tf.name)), minio.PutObjectOptions{})
		if err != nil {
			fmt.Printf("  [出错] 上传测试文件 %q 失败: %v\n", tf.key, err)
		} else {
			fmt.Printf("  [成功] 已创建: %s\n", tf.key)
		}
	}

	// 测试不同的前缀查询
	prefixes := []struct {
		name     string
		prefix   string
		delimiter string
	}{
		{name: "列出所有文件 (前缀为空)", prefix: "", delimiter: ""},
		{name: "列出 dir_a/ 目录下所有文件", prefix: "dir_a/", delimiter: ""},
		{name: "列出 dir_a/ 目录下直接子文件 (带分隔符 /)", prefix: "dir_a/", delimiter: "/"},
		{name: "列出 dir_b/ 目录下所有文件", prefix: "dir_b/", delimiter: ""},
		{name: "列出根目录文件 (分隔符 /)", prefix: "", delimiter: "/"},
		{name: "列出 root_ 前缀的文件", prefix: "root_", delimiter: ""},
	}

	for _, p := range prefixes {
		fmt.Printf("\n--- 查询: %s ---\n", p.name)
		fmt.Printf("    prefix=%q, delimiter=%q\n", p.prefix, p.delimiter)

		opts := minio.ListObjectsOptions{
			Prefix:    p.prefix,
			Delimiter: p.delimiter,
		}

		objCount := 0
		for objInfo := range client.ListObjects(ctx, bucket, opts) {
			if objInfo.Err != nil {
				fmt.Printf("  [出错] ListObjects 遍历错误: %v\n", objInfo.Err)
				break
			}
			objCount++
			if objInfo.IsDir {
				fmt.Printf("  [目录] %s\n", objInfo.Prefix)
			} else {
				fmt.Printf("  [文件] %s (大小: %d bytes, ETag: %s)\n",
					objInfo.Key, objInfo.Size, objInfo.ETag)
			}
		}

		fmt.Printf("  --- 共 %d 个对象 ---\n", objCount)
	}

	// 清理测试文件
	fmt.Println("\n--- 清理目录测试文件 ---")
	for _, tf := range testFiles {
		err := client.RemoveObject(ctx, bucket, tf.key, minio.RemoveObjectOptions{})
		if err != nil {
			fmt.Printf("  [出错] 删除 %q 失败: %v\n", tf.key, err)
		}
	}
	fmt.Println("  [完成] 清理完毕")
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	cfg := parseFlags()

	fmt.Println("============================================================")
	fmt.Println("  S3 存储测试工具")
	fmt.Println("============================================================")
	fmt.Printf("  目标 S3: %s\n", cfg.Endpoint)
	fmt.Printf("  存储桶:   %s\n", cfg.Bucket)
	fmt.Printf("  UseSSL:   %v\n", cfg.UseSSL)
	fmt.Printf("  Region:   %q\n", cfg.Region)
	fmt.Printf("  SkipSSL:  %v\n", cfg.SkipSSLVerify)
	fmt.Println()

	// 打印将要调用的 S3 API 列表
	fmt.Println("============================================================")
	fmt.Println("  将要调用的 S3 API 列表:")
	fmt.Println("============================================================")
	fmt.Println("  1️⃣  HeadBucket       — 检查桶是否存在")
	fmt.Println("  2️⃣  ListBuckets      — 列出所有桶")
	fmt.Println("  3️⃣  PutObject        — 上传对象 (多种字符串测试)")
	fmt.Println("  4️⃣  GetObject        — 下载对象并验证内容")
	fmt.Println("  5️⃣  RemoveObject     — 删除对象")
	fmt.Println("  6️⃣  StatObject       — 获取对象元数据 (删除验证)")
	fmt.Println("  7️⃣  ListObjects      — 按前缀列出对象")
	fmt.Println()

	ctx := context.Background()

	// 创建客户端
	client, err := newClient(cfg)
	if err != nil {
		log.Fatalf("[致命错误] 创建 S3 客户端失败: %v", err)
	}
	fmt.Println("[成功] S3 客户端创建成功")

	// 步骤 1: 获取桶信息 (失败继续)
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  开始执行测试步骤...")
	fmt.Println("════════════════════════════════════════════════════════════")

	step1_GetBucketInfo(ctx, client, cfg.Bucket)

	// 步骤 2: 上传文件
	step2_UploadFile(ctx, client, cfg.Bucket)

	// 步骤 3: 下载并验证
	step3_DownloadAndVerify(ctx, client, cfg.Bucket)

	// 步骤 4: 删除文件
	step4_DeleteFile(ctx, client, cfg.Bucket)

	// 步骤 5: 获取文件列表
	step5_ListFiles(ctx, client, cfg.Bucket)

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  所有测试步骤执行完毕!")
	fmt.Println("============================================================")
}
