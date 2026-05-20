# MKF 文件格式文档

## 概述

MKF (Multi-Kingdom File) 是《仙剑奇侠传》DOS版使用的一种资源打包格式。该格式将多个数据块（chunk）存储在单个文件中，支持YJ1_压缩算法。本项目实现了对多种MKF文件的解析和处理。

## 核心数据类型

```go
type (
    INT    uint32  // 32位无符号整数
    WORD   = uint16 // 16位无符号整数
    SHORT  = int16  // 16位有符号整数
    USHORT = uint16 // 16位无符号整数
    DWORD  = uint32 // 32位无符号整数
)
```

## 1. MKF 文件结构

### 1.1 文件头

MKF文件采用统一的文件头结构：

| 偏移 | 大小 | 字段 | 说明 |
|------|------|------|------|
| 0x00 | 4字节 | ChunkCount | 块数量 + 1（实际块数 = ChunkCount - 4 >> 2） |
| 0x04 | 4 * N 字节 | OffsetTable | 每个块的偏移量表 |

### 1.2 块访问方法

```go
// 获取块数量
func (mkf *Mkf) GetChunkCount() (INT, error)

// 获取指定块的偏移
func (mkf *Mkf) GetChunkOffset(chunkNum INT) (INT, error)

// 获取指定块的大小
func (mkf *Mkf) GetChunkSize(chunkNum INT) (INT, error)

// 读取指定块的原始数据
func (mkf *Mkf) ReadChunk(chunkNum INT) ([]byte, error)
```

## 2. YJ1_ 压缩算法

### 2.1 压缩文件头

```go
type _YJ1_FILEHEADER struct {
    Signature          INT    // "YJ1_" (0x315F4A59)
    UncompressedLength INT    // 解压后长度
    CompressedLength   INT    // 压缩后长度（含头部）
    BlockCount         uint16 // 块数量
    Unknown            uint8  // 未知
    HuffmanTreeLength  uint8  // Huffman树长度
}
```

### 2.2 压缩块结构

```go
type _YJ_1_BLOCKHEADER struct {
    UncompressedLength        uint16 // 解压后块大小（最大0x4000）
    CompressedLength          uint16 // 压缩后块大小
    LZSSRepeatTable           [4]uint16
    LZSSOffsetCodeLengthTable [4]uint8
    LZSSRepeatCodeLengthTable [3]uint8
    CodeCountCodeLengthTable  [3]uint8
    CodeCountTable            [2]uint8
}
```

### 2.3 解压流程

1. 读取文件头，验证签名
2. 构建Huffman树
3. 逐块解压：
    - 如果 `CompressedLength == 0`：直接复制原始数据
    - 否则：使用Huffman解码 + LZSS重复机制解压

## 3. DATA.MKF - 游戏数据文件

### 3.1 块索引

| Chunk编号 | 内容 | 数据结构 |
|-----------|------|----------|
| 0 | 商店数据 | `StoreChunk` |
| 1 | 敌人数据 | `EnemyChunk` |
| 2 | 敌人队伍 | `EnemyTeamChunk` |
| 3 | 玩家角色 | `PlayerRoles` |
| 4 | 魔法数据 | `MagicChunk` |
| 5 | 战场数据 | `BattleFieldChunk` |
| 6 | 升级魔法 | `LevelUpMagicChunk` |
| 11 | 战斗效果索引 | `[10][2]WORD` |
| 13 | 敌人位置 | `EnemyPos` |
| 14 | 升级经验 | `[MAX_LEVELS+1]WORD` |

### 3.2 核心数据结构

#### 3.2.1 Enemy（敌人数据）

| 字段 | 类型 | 说明 |
|------|------|------|
| IdleFrames | WORD | 空闲动画帧数 |
| MagicFrames | WORD | 施法动画帧数 |
| AttackFrames | WORD | 攻击动画帧数 |
| IdleAnimSpeed | WORD | 空闲动画速度 |
| Health | WORD | 最大HP |
| Exp | WORD | 经验值 |
| Cash | WORD | 金钱 |
| Level | WORD | 等级 |
| Magic | WORD | 魔法编号 |
| AttackStrength | WORD | 攻击力 |
| MagicStrength | WORD | 魔法攻击力 |
| Defense | WORD | 防御力 |
| Dexterity | WORD | 敏捷度 |
| PoisonResistance | WORD | 抗毒能力 |
| ElemResistance | [5]WORD | 元素抗性 |

#### 3.2.2 Magic（魔法数据）

| 字段 | 类型 | 说明 |
|------|------|------|
| Effect | WORD | 效果精灵编号 |
| Type | WORD | 魔法类型 |
| Speed | SHORT | 施法速度 |
| CostMP | WORD | MP消耗 |
| BaseDamage | WORD | 基础伤害 |
| Elemental | WORD | 元素属性 |
| Sound | SHORT | 音效编号 |

#### 3.2.3 PlayerRoles（玩家角色）

| 字段 | 类型 | 说明 |
|------|------|------|
| Avatar | [6]WORD | 头像编号 |
| SpriteNum | [6]WORD | 场景精灵编号 |
| SpriteNumInBattle | [6]WORD | 战斗精灵编号 |
| Name | [6]WORD | 名称编号 |
| Level | [6]WORD | 等级 |
| MaxHP | [6]WORD | 最大HP |
| MaxMP | [6]WORD | 最大MP |
| HP | [6]WORD | 当前HP |
| MP | [6]WORD | 当前MP |
| Equipment | [6][6]WORD | 装备 |
| Magic | [32][6]WORD | 魔法 |

## 4. FBP.MKF - 背景图片文件

### 4.1 用途

存储游戏中的背景图片，包括：
- 主菜单背景（Chunk 60）
- 状态界面背景（Chunk 0）
- 其他界面背景

### 4.2 关键常量

```go
const (
    STATUS_BACKGROUND_FBPNUM   = 0    // 状态界面背景
    MAINMENU_BACKGROUND_FBPNUM = 60   // 主菜单背景
    BITMAPNUM_SPLASH_UP        = 0x26 // 开场动画上半部分
    BITMAPNUM_SPLASH_DOWN      = 0x27 // 开场动画下半部分
)
```

### 4.3 数据结构

```go
type FbpChunk struct{ CompressedChunk }

// 解压并转换为位图
func (c *FbpChunk) GetBmp() *BitMap
```

## 5. MGO.MKF - 精灵图形文件

### 5.1 用途

存储游戏角色和NPC的精灵图形。

### 5.2 关键常量

```go
const (
    SPRITENUM_SPLASH_TITLE = 0x47 // 开场标题精灵
    SPRITENUM_SPLASH_CRANE = 0x49 // 开场仙鹤精灵
)
```

### 5.3 数据结构

```go
type MgoChunk struct{ BitMapChunk }

// 获取指定块（每帧大小为32000字节）
func (mm *MgoMkf) GetChunk(chunkNum INT) (MgoChunk, error)
```

## 6. SSS.MKF - 脚本与场景文件

### 6.1 块索引

| Chunk编号 | 内容 | 数据结构 |
|-----------|------|----------|
| 0 | 事件对象 | `EventObjectChunk` |
| 1 | 场景数据 | `SceneChunk` |
| 2 | 对象定义 | Object数组 |
| 3 | 消息偏移 | `MsgOffsetChunk` |
| 4 | 脚本条目 | `ScriptEntryChunk` |

### 6.2 核心数据结构

#### 6.2.1 EventObject（事件对象）

| 字段 | 类型 | 说明 |
|------|------|------|
| X, Y | WORD | 地图坐标 |
| Layer | SHORT | 图层（地面/空中） |
| TriggerScript | WORD | 触发脚本入口 |
| AutoScript | WORD | 自动脚本入口 |
| State | SHORT | 对象状态 |
| TriggerMode | WORD | 触发模式 |
| SpriteNum | WORD | 精灵编号 |
| SpriteFrames | USHORT | 精灵帧数 |
| Direction | WORD | 朝向 |
| CurrentFrameNum | WORD | 当前帧 |

#### 6.2.2 Scene（场景）

| 字段 | 类型 | 说明 |
|------|------|------|
| MapNum | WORD | 地图编号 |
| ScriptOnEnter | WORD | 进入场景脚本 |
| ScriptOnTeleport | WORD | 传送脚本 |
| EventObjectIndex | WORD | 事件对象起始索引 |

#### 6.2.3 ScriptEntry（脚本条目）

| 字段 | 类型 | 说明 |
|------|------|------|
| Operation | WORD | 操作码 |
| Operand[3] | WORD | 操作数（最多3个） |

#### 6.2.4 Object（对象定义）

对象定义支持多种类型：

**物品对象 (ObjectItem)**：
| 字段 | 类型 | 说明 |
|------|------|------|
| Bitmap | WORD | 图标编号 |
| Price | WORD | 价格 |
| ScriptOnEquip | WORD | 装备脚本 |
| ScriptOnThrow | WORD | 投掷脚本 |
| ScriptDesc | WORD | 描述脚本 |
| Flags | WORD | 属性标志 |

**魔法对象 (ObjectMagic)**：
| 字段 | 类型 | 说明 |
|------|------|------|
| MagicNumber | WORD | 魔法编号 |
| ScriptOnSuccess | WORD | 成功脚本 |
| ScriptOnUse | WORD | 使用脚本 |

**敌人对象 (ObjectEnemy)**：
| 字段 | 类型 | 说明 |
|------|------|------|
| EnemyID | WORD | 敌人编号 |
| ResistanceToSorcery | WORD | 抗魔性 |
| ScriptOnTurnStart | WORD | 回合开始脚本 |
| ScriptOnBattleEnd | WORD | 战斗结束脚本 |

## 7. 位图格式

### 7.1 位图结构

```go
type BitMap struct {
    data []byte  // 像素数据
    w, h INT     // 宽度和高度
    rle  bool    // 是否为RLE压缩
}
```

### 7.2 RLE 压缩格式

RLE（行程长度编码）格式用于压缩位图数据：

| 类型 | 格式 | 说明 |
|------|------|------|
| 位置跳转 | `0x80 + offset` | 跳过指定像素数 |
| 数据块 | `count + [count]byte` | 连续count个像素 |

### 7.3 解码流程

```go
func (bmp *BitMap) rleToImage(plt []color.RGBA, shadow bool) *ebiten.Image {
    // 1. 初始化图像
    // 2. 遍历RLE数据
    // 3. 处理位置跳转或数据块
    // 4. 应用调色板转换
}
```

## 8. 调色板系统

### 8.1 调色板结构

游戏使用256色调色板，每个颜色由RGBA四个分量组成。

### 8.2 像素转换

```go
func pixToRGBA(pix byte, plt []color.RGBA) color.Color {
    return plt[pix]
}
```

## 9. 关键辅助函数

### 9.1 PlaneChunk - 平面数据块

```go
type PlaneChunk struct {
    data    []byte  // 原始数据
    eleSize int     // 每个元素的大小
}

func (pc *PlaneChunk) Len() int              // 获取元素数量
func (pc *PlaneChunk) Get(idx int) unsafe.Pointer // 获取指定索引的元素指针
```

### 9.2 FrameChunk - 帧数据块

```go
type FrameChunk struct {
    data      []byte
    frameSize int  // 每帧大小
}

func (c *FrameChunk) GetCount() INT            // 获取帧数
func (c *FrameChunk) GetFrame(frameNum INT) ([]byte, error) // 获取指定帧
```

## 10. 文件格式总结

| 文件 | 用途 | 压缩方式 | 主要内容 |
|------|------|----------|----------|
| DATA.MKF | 游戏数据 | 无 | 敌人、魔法、角色、商店 |
| FBP.MKF | 背景图片 | YJ1_ | 菜单、状态界面背景 |
| MGO.MKF | 精灵图形 | YJ1_ | 角色、NPC精灵 |
| SSS.MKF | 脚本场景 | 无 | 事件对象、场景、脚本 |

## 附录：常量定义

```go
const (
    MAX_STORE_ITEM            = 9        // 商店最大物品数
    NUM_MAGIC_ELEMENTAL       = 5        // 元素类型数
    MAX_ENEMIES_IN_TEAM       = 5        // 最大敌人数
    MAX_PLAYABLE_PLAYER_ROLES = 5        // 可玩角色数
    MAX_LEVELS                = 99       // 最大等级
    MAX_PLAYER_ROLES          = 6        // 玩家角色数
    MAX_PLAYER_EQUIPMENTS     = 6        // 装备槽数
    MAX_PLAYER_MAGICS         = 32       // 最大魔法数
    MAX_SCENES                = 300      // 最大场景数
    MAX_OBJECTS               = 600      // 最大对象数
)
```