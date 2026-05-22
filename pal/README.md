# Pal - 仙剑奇侠传重置项目

基于 Go 语言和 Ebiten 游戏引擎的仙剑奇侠传重置项目。

## 项目结构

```
pal/
├── cmd/                    # 命令行工具入口
│   ├── pal/               # 主游戏程序
│   ├── mgodump/           # MGO.MKF 文件查看工具
│   ├── rgmdump/           # RGM.MKF 文件查看工具
│   └── sssdump/           # SSS.MKF 文件查看工具
├── pkg/                    # 库代码
│   ├── game/              # 游戏核心逻辑
│   ├── mkf/               # MKF 文件解析库
│   ├── ui/                # UI 组件库
│   └── utils/             # 工具函数
├── data/                  # 游戏数据文件（M.MSG, WORD.DAT, wending.ttf 等）
└── doc/                   # 项目文档
```

## 环境要求

- Go 1.20+
- Ebiten 游戏引擎

## 构建与运行

### 构建主游戏

```bash
cd pal
go build -o pal cmd/pal/main.go
```

### 运行主游戏

```bash
./pal
```

### 构建工具

```bash
# MGO.MKF 查看工具
go build -o mgodump cmd/mgodump/main.go

# RGM.MKF 查看工具
go build -o rgmdump cmd/rgmdump/main.go

# SSS.MKF 查看工具
go build -o sssdump cmd/sssdump/main.go
```

### 运行工具

```bash
# MGO.MKF 查看工具
./mgodump -g ./data

# RGM.MKF 查看工具
./rgmdump -g ./data

# SSS.MKF 查看工具
./sssdump -f ./data/SSS.MKF
```

## 调试

项目已配置 VS Code 调试配置，可在 `.vscode/launch.json` 中查看和修改调试配置。

## 项目文档

- `doc/mkf_file_formats.md` - MKF 文件格式说明
- `doc/script_system.md` - 脚本系统说明
- `doc/ui_framework.md` - UI 框架说明

## 游戏数据

游戏运行需要以下数据文件放在 `data/` 目录：
- `M.MSG` - 消息文本文件
- `WORD.DAT` - 字库文件
- `SSS.MKF` - 脚本数据文件
- `DATA.MKF` - 游戏数据文件
- `MGO.MKF` - 图像资源文件
- `RGM.MKF` - 头像资源文件
- `wending.ttf` - 字体文件

## 许可证

本项目仅供学习和研究使用。
