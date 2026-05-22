# UI 框架文档

## 概述

`framework/UI` 是一个基于 Ebiten 引擎构建的轻量级游戏 UI 框架，提供了组件化的界面开发能力。该框架采用面向对象设计，支持窗口管理、事件处理、动画播放等核心功能。

## 1. 架构设计

### 1.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        WinManager                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    BasicWindow                            │  │
│  │  ┌─────────────────────────────────────────────────────┐  │  │
│  │  │           BasicComponent (基类)                      │  │  │
│  │  │  ├── Button        ├── Image        ├── Animation   │  │  │
│  │  │  └── 自定义组件... └── 自定义组件... └── 其他组件... │  │  │
│  │  └─────────────────────────────────────────────────────┘  │  │
│  │  └── popUpWin (弹出窗口)                                  │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 核心接口

| 接口 | 职责 | 关键方法 |
|------|------|----------|
| `Game` | 游戏主循环 | `Update()`, `Draw()`, `Layout()` |
| `GameComponent` | 游戏组件 | `ID()`, `Rect()`, `Parent()`, 鼠标事件 |
| `ParentComponent` | 父组件 | `Notify()` |
| `Window` | 窗口容器 | `AddComponent()`, `RemoveComponent()`, `Close()` |
| `IAnimation` | 动画接口 | `Update()`, `GetImage()` |

## 2. 核心组件

### 2.1 BasicComponent - 基础组件

**定义位置**: `framework/UI/component.go`

```go
type BasicComponent struct {
    id     int
    RECT   Rect
    parent ParentComponent
    
    onClick     func(x, y int)
    onMouseDown func(x, y int)
    onUpdate    func() error
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `ID()` | 获取组件唯一标识 |
| `Rect()` | 获取组件矩形区域 |
| `Parent()` | 获取父组件 |
| `SetOnClick(f)` | 设置点击回调 |
| `MouseDown(x, y)` | 处理鼠标按下事件 |
| `MouseUp(x, y)` | 处理鼠标释放事件 |
| `MouseMove(x, y)` | 处理鼠标移动事件 |
| `MouseLeave()` | 处理鼠标离开事件 |

### 2.2 BasicWindow - 基础窗口

**定义位置**: `framework/UI/window.go`

```go
type BasicWindow struct {
    BasicComponent
    components []GameComponent
    
    popUpWin         GameComponent  // 弹出窗口
    focusedCom       GameComponent  // 焦点组件
    mouseDriftingCom GameComponent  // 鼠标漂移组件
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `Update()` | 更新所有子组件 |
| `Draw(screen)` | 绘制所有子组件 |
| `AddComponent(c)` | 添加子组件 |
| `RemoveSubWin(w)` | 移除子窗口 |
| `Pop(com)` | 弹出窗口 |
| `Notify(subId, event, msg)` | 事件通知 |
| `MouseDown/Up/Move(x, y)` | 鼠标事件分发 |

**事件分发机制**:

```
用户操作 → BasicWindow → 从后向前遍历组件 → 找到覆盖点的组件 → 触发事件
```

### 2.3 Button - 按钮组件

**定义位置**: `framework/UI/button_com.go`

```go
type Button struct {
    BasicComponent
    
    imgButtonUp    *ebiten.Image
    imgButtongDown *ebiten.Image
    isMouseDown    bool
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `Button(width, height)` | 设置按钮尺寸 |
| `SetLocation(x, y)` | 设置按钮位置 |
| `AddButtonUpImage(img)` | 设置按钮抬起状态图片 |
| `AddButtonDownImage(img)` | 设置按钮按下状态图片 |
| `Draw(screen)` | 根据状态绘制按钮 |

**链式调用示例**:

```go
btn := NewButton(parent).
    Button(100, 50).
    SetLocation(10, 10).
    AddButtonUpImage(imgUp).
    AddButtonDownImage(imgDown)
```

### 2.4 Image - 图片组件

**定义位置**: `framework/UI/image_com.go`

```go
type Image struct {
    BasicComponent
    image     *ebiten.Image
    autoScale bool
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `SetSize(width, height)` | 设置显示尺寸 |
| `SetLocation(x, y)` | 设置位置 |
| `SetAutoScale(auto)` | 设置是否自动缩放 |
| `LoadImage(img)` | 加载图片 |
| `Draw(screen)` | 绘制图片（支持自动缩放） |

### 2.5 Animation - 动画组件

**定义位置**: `framework/UI/animation.go`

```go
type Animation struct {
    images    []*ebiten.Image
    interval  int
    paused    bool
    crtImgIdx int
}
```

**核心方法**:

| 方法 | 功能 |
|------|------|
| `Pause()` | 暂停动画 |
| `Resume()` | 恢复动画 |
| `Update(count)` | 更新动画帧 |
| `GetImage()` | 获取当前帧图片 |
| `AppendImage(img)` | 添加单张图片 |
| `AppendImages(imgs)` | 添加多张图片 |

**默认参数**:

```go
const defaultInterval = 20  // 默认帧间隔
```

## 3. 窗口管理

### 3.1 WinManager - 窗口管理器

**定义位置**: `framework/UI/window_manager.go`

```go
type WinManager struct {
    winStack []*Window
}
```

**设计意图**: 管理多个窗口的层级关系，支持窗口栈操作。

### 3.2 窗口层级

| 层级 | 说明 |
|------|------|
| 主窗口 | 底层窗口，承载主要UI |
| 子组件 | 主窗口内的各种控件 |
| 弹出窗口 | 临时弹窗，位于最上层 |

## 4. 事件系统

### 4.1 事件类型

**定义位置**: `framework/UI/event.go`

```go
type ComEvent int

const (
    OnClose ComEvent = iota  // 关闭事件
)
```

### 4.2 事件处理流程

```
组件触发事件 → Notify(parent, event, msg) → 父组件处理
```

### 4.3 Input - 输入处理

**定义位置**: `framework/UI/input.go`

```go
type Input struct {
    device InputDevice
    owner  GameComponent
}
```

**输入设备类型**:

| 类型 | 值 | 说明 |
|------|-----|------|
| `NoneDevice` | 0 | 无设备 |
| `Mouse` | 1 | 鼠标 |
| `Keyboard` | 2 | 键盘 |

**输入更新流程**:

```go
func (i *Input) Update() {
    if i.device&Mouse != 0 {
        i.updateMouseEvent()
    }
}
```

## 5. 辅助工具

### 5.1 Rect - 矩形区域

**定义位置**: `framework/UI/common.go`

```go
type Rect struct {
    Top, Left, Height, Width int
}

func (r Rect) Cover(x, y int) bool  // 判断点是否在矩形内
```

**坐标系统**:

- `Top`: 顶部Y坐标
- `Left`: 左侧X坐标
- `Height`: 高度
- `Width`: 宽度

## 6. 渲染机制

### 6.1 组件绘制流程

```go
func (bw *BasicWindow) drawCompoent(screen *ebiten.Image, com GameComponent) {
    // 1. 获取组件矩形
    rect := com.Rect()
    
    // 2. 获取组件布局尺寸
    sw, sh := com.Layout(rect.Width, rect.Height)
    
    // 3. 创建中间缓冲区
    img := ebiten.NewImage(sw, sh)
    
    // 4. 组件自绘
    com.Draw(img)
    
    // 5. 应用缩放变换
    op := &ebiten.DrawImageOptions{}
    if rect.Width != sw || rect.Height != sh {
        op.GeoM.Scale(...)
    }
    
    // 6. 应用位置变换
    op.GeoM.Translate(float64(rect.Left), float64(rect.Top))
    
    // 7. 绘制到屏幕
    screen.DrawImage(img, op)
}
```

### 6.2 坐标系转换

组件内部坐标 → 窗口坐标 → 屏幕坐标

## 7. 使用示例

### 7.1 创建窗口和按钮

```go
// 创建窗口管理器
wm := NewWinManager()

// 创建主窗口
mainWin := NewBasicWindow(640, 480, nil)

// 创建按钮
btn := NewButton(mainWin).
    Button(100, 40).
    SetLocation(50, 50)

// 添加到窗口
mainWin.AddComponent(btn)
```

### 7.2 创建动画

```go
// 创建动画（每20帧切换）
anim := NewAnimation(20)

// 添加帧图片
anim.AppendImages([]*ebiten.Image{img1, img2, img3, img4})

// 在游戏循环中更新
func (g *Game) Update() error {
    anim.Update(frameCount)
    return nil
}

// 绘制当前帧
func (g *Game) Draw(screen *ebiten.Image) {
    screen.DrawImage(anim.GetImage(), nil)
}
```

### 7.3 绑定输入

```go
input := NewInput().
    Bind(mainWin).
    AddDevice(Mouse)

// 在游戏循环中更新
func (g *Game) Update() error {
    input.Update()
    return nil
}
```

## 8. 组件生命周期

```
创建 → 添加到父容器 → Update/Draw → 移除 → 销毁
    ↓           ↓              ↓
  NewXxx    AddComponent   游戏循环   RemoveComponent
```

## 9. 设计模式

| 模式 | 应用 | 说明 |
|------|------|------|
| 组合模式 | Window-Component | 窗口包含多个子组件 |
| 观察者模式 | Notify机制 | 组件事件通知 |
| 模板方法 | Draw/Update | 统一的生命周期方法 |
| 建造者模式 | 链式调用 | Button/Image的链式配置 |

## 10. 扩展建议

### 10.1 添加新组件

```go
type CustomComponent struct {
    BasicComponent
    // 自定义字段
}

func NewCustomComponent(parent ParentComponent) *CustomComponent {
    ret := CustomComponent{}
    ret.BasicComponent = *NewConponent(0, 0, 0, 0, parent)
    return &ret
}

func (c *CustomComponent) Draw(screen *ebiten.Image) {
    // 自定义绘制逻辑
}
```

### 10.2 添加新事件类型

```go
const (
    OnClose ComEvent = iota
    OnClick
    OnHover
    OnFocus
)
```

## 附录：文件清单

| 文件 | 功能 |
|------|------|
| `interface.go` | 定义核心接口 |
| `component.go` | 基础组件实现 |
| `window.go` | 窗口组件实现 |
| `window_manager.go` | 窗口管理器 |
| `button_com.go` | 按钮组件 |
| `image_com.go` | 图片组件 |
| `animation.go` | 动画组件 |
| `event.go` | 事件定义 |
| `input.go` | 输入处理 |
| `common.go` | 通用工具 |