# Klik 服务器端

这是 Klik 前端项目的后端服务器实现，使用 Go 语言和 Gin 框架开发，实现了与前端 axios-mock 相同的接口功能。

## 项目结构

```
server/
├── config/         # 配置相关
├── controller/     # 控制器，处理请求
├── model/          # 数据模型
├── router/         # 路由定义
├── utils/          # 工具函数
├── go.mod          # Go 模块定义
├── main.go         # 主程序入口
└── README.md       # 说明文档
```

## 功能列表

该服务器实现了以下接口：

- `/video/recommended` - 获取推荐视频
- `/video/long/recommended` - 获取长视频推荐
- `/video/comments` - 获取视频评论
- `/video/private` - 获取私有视频
- `/video/like` - 获取喜欢的视频
- `/video/my` - 获取我的视频
- `/video/history` - 获取历史视频
- `/user/collect` - 获取用户收藏
- `/user/video_list` - 获取用户视频列表
- `/user/panel` - 获取用户面板信息
- `/user/friends` - 获取用户好友
- `/historyOther` - 获取其他历史记录
- `/post/recommended` - 获取推荐帖子
- `/shop/recommended` - 获取推荐商品

## 安装与运行

### 前提条件

- 安装 Go 1.16 或更高版本
- 确保 `public/data` 目录中有必要的数据文件

### 安装依赖

```bash
cd server
go mod download
```

### 运行服务器

```bash
cd server
go run main.go
```

服务器将在 http://localhost:8080 上启动，可以通过 `/api/...` 路径访问各个接口。

## 数据文件

服务器需要以下数据文件（与前端 mock 数据相同）：

- `public/data/resource.js` - 资源数据
- `public/data/posts6.json` - 视频数据
- `public/data/posts.md` - 帖子数据
- `public/data/users.md` - 用户数据
- `public/data/goods.md` - 商品数据
- `public/data/user_video_list/user-*.md` - 用户视频列表
- `public/data/comments/video_id_*.md` - 视频评论

## 与前端集成

修改前端项目中的 API 基础路径，指向该服务器：

```javascript
// 例如在前端配置文件中
const API_BASE_URL = 'http://localhost:8080/api';
```

然后禁用前端的 axios-mock，使用真实的 API 请求。
