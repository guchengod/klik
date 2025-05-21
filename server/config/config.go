package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server struct {
		BaseURL string `yaml:"baseURL"`
		FileURL string `yaml:"fileURL"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		DBName       string `yaml:"dbname"`
		SSLMode      string `yaml:"sslmode"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		UseDB        bool   `yaml:"useDB"`
	} `yaml:"database"`

	Paths struct {
		DataPath         string `yaml:"dataPath"`
		UserVideoListPath string `yaml:"userVideoListPath"`
		CommentsPath     string `yaml:"commentsPath"`
	} `yaml:"paths"`
}

var (
	// AppConfig 全局应用配置
	AppConfig Config

	// BaseURL 基础URL
	BaseURL string
	// FileURL 文件URL
	FileURL string
	// DataPath 数据路径
	DataPath string
	// UseDB 是否使用数据库
	UseDB bool
)

// Init 初始化配置
func Init() {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	// 获取配置文件路径
	configPath := filepath.Join(filepath.Dir(filename), "config.yaml")
	// 获取项目根目录
	rootDir := filepath.Join(filepath.Dir(filename), "../..")

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	// 解析YAML配置
	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	// 设置全局变量
	BaseURL = AppConfig.Server.BaseURL
	FileURL = AppConfig.Server.FileURL
	UseDB = AppConfig.Database.UseDB

	// 处理相对路径
	DataPath = filepath.Join(rootDir, AppConfig.Paths.DataPath)
	userVideoListPath := filepath.Join(rootDir, AppConfig.Paths.UserVideoListPath)
	commentsPath := filepath.Join(rootDir, AppConfig.Paths.CommentsPath)
	
	// 确保数据目录存在
	ensureDir(DataPath)
	ensureDir(userVideoListPath)
	ensureDir(commentsPath)

	// 初始化数据库
	if UseDB {
		InitDB()
	}

	log.Printf("配置加载成功: 服务器端口=%d, 数据库=%s@%s:%d/%s", 
		AppConfig.Server.Port, 
		AppConfig.Database.User, 
		AppConfig.Database.Host, 
		AppConfig.Database.Port, 
		AppConfig.Database.DBName)
}

// 确保目录存在
func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

// Close 关闭资源
func Close() {
	// 关闭数据库连接
	if UseDB {
		CloseDB()
	}
}
