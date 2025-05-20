package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	// BaseURL 基础URL
	BaseURL string
	// FileURL 文件URL
	FileURL string
	// DataPath 数据路径
	DataPath string
)

// Init 初始化配置
func Init() {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	// 获取项目根目录
	rootDir := filepath.Join(filepath.Dir(filename), "../..")
	// 设置数据路径
	DataPath = filepath.Join(rootDir, "public", "data")
	
	// 确保数据目录存在
	ensureDir(DataPath)
	ensureDir(filepath.Join(DataPath, "user_video_list"))
	ensureDir(filepath.Join(DataPath, "comments"))
	
	// 设置URL
	BaseURL = "/api"
	FileURL = "/api/file"
}

// 确保目录存在
func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}
