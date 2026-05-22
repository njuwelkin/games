package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/njuwelkin/games/pal/pkg/mkf"
	"github.com/njuwelkin/games/pal/pkg/utils"
)

// 方向常量（与 pal/common.go 中的定义一致）
const (
	kDirSouth = iota
	kDirWest
	kDirNorth
	kDirEast
	kDirUnknown
)

// 脚本操作码常量
const (
	OP_STOP = 0x0000
)

var (
	filePath            = flag.String("f", "SSS.MKF", "SSS.MKF 文件路径")
	listScenes          = flag.Bool("scenes", false, "列出所有场景")
	listObjects         = flag.Bool("objects", false, "列出所有事件对象")
	listScripts         = flag.Bool("scripts", false, "列出所有脚本条目")
	listAll             = flag.Bool("all", false, "显示所有内容")
	sceneIndex          = flag.Int("scene", -1, "显示指定索引的场景详情")
	objectIndex         = flag.Int("object", -1, "显示指定索引的事件对象详情")
	sceneEnterScript    = flag.Int("scene-enter-script", -1, "显示指定场景的进入脚本内容")
	sceneTeleportScript = flag.Int("scene-teleport-script", -1, "显示指定场景的传送脚本内容")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "SSS.MKF 文件解析工具\n\n")
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -all\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -scenes\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -scene 10\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -objects\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -object 5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -scene-enter-script 10\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f SSS.MKF -scene-teleport-script 10\n", os.Args[0])
	}

	flag.Parse()

	if *filePath == "" {
		fmt.Println("错误: 必须指定文件路径")
		flag.Usage()
		os.Exit(1)
	}

	sss, err := mkf.NewSSSMkf(*filePath)
	if err != nil {
		fmt.Printf("无法打开文件 %s: %v\n", *filePath, err)
		os.Exit(1)
	}
	defer sss.Close()

	// 获取脚本块（用于显示脚本内容）
	scriptChunk, err := sss.GetScriptEntryChunk()
	if err != nil {
		fmt.Printf("读取脚本数据失败: %v\n", err)
	}

	// 获取消息偏移表和消息数据（用于显示文本内容）
	var msgOffsetChunk *mkf.MsgOffsetChunk
	var msgData []byte

	msgOffsetChunk, err = sss.GetMsgOffsetChunk()
	if err != nil {
		fmt.Printf("读取消息偏移表失败: %v\n", err)
	}

	msgData, err = os.ReadFile("../../M.MSG")
	if err != nil {
		fmt.Printf("读取 M.MSG 文件失败: %v\n", err)
	}

	if *listAll || *listScenes || *sceneIndex >= 0 {
		showScenes(&sss, scriptChunk, msgOffsetChunk, msgData, *sceneIndex, *listAll || *listScenes)
	}

	if *listAll || *listObjects || *objectIndex >= 0 {
		showEventObjects(&sss, *objectIndex, *listAll || *listObjects)
	}

	if *listAll || *listScripts {
		showScripts(&sss)
	}

	// 显示场景进入脚本内容
	if *sceneEnterScript >= 0 {
		showSceneEnterScript(&sss, scriptChunk, msgOffsetChunk, msgData, *sceneEnterScript)
	}

	// 显示场景传送脚本内容
	if *sceneTeleportScript >= 0 {
		showSceneTeleportScript(&sss, scriptChunk, msgOffsetChunk, msgData, *sceneTeleportScript)
	}

	if !*listAll && !*listScenes && !*listObjects && !*listScripts &&
		*sceneIndex < 0 && *objectIndex < 0 && *sceneEnterScript < 0 && *sceneTeleportScript < 0 {
		flag.Usage()
	}
}

func showScenes(sss *mkf.SSSMkf, scriptChunk *mkf.ScriptEntryChunk, msgOffsetChunk *mkf.MsgOffsetChunk, msgData []byte, idx int, listAll bool) {
	chunk, err := sss.GetSceneChunk()
	if err != nil {
		fmt.Printf("读取场景数据失败: %v\n", err)
		return
	}

	count := chunk.Len()

	if idx >= 0 {
		if idx >= count {
			fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", idx, count-1)
			return
		}
		scene := chunk.GetScene(idx)
		fmt.Printf("\n=== 场景 #%d ===\n", idx)
		printScene(scene)

		// 如果有脚本块，显示脚本内容
		if scriptChunk != nil {
			if scene.ScriptOnEnter != 0 {
				fmt.Println("\n  --- 进入脚本内容 ---")
				printScriptContent(scriptChunk, msgOffsetChunk, msgData, int(scene.ScriptOnEnter))
			}
			if scene.ScriptOnTeleport != 0 {
				fmt.Println("\n  --- 传送脚本内容 ---")
				printScriptContent(scriptChunk, msgOffsetChunk, msgData, int(scene.ScriptOnTeleport))
			}
		}
	} else if listAll {
		fmt.Println("\n=== 场景列表 ===")
		for i := 0; i < count; i++ {
			scene := chunk.GetScene(i)
			if scene.MapNum == 0 && scene.ScriptOnEnter == 0 && scene.ScriptOnTeleport == 0 && scene.EventObjectIndex == 0 {
				continue
			}
			fmt.Printf("[%4d] 地图:%4d 进入脚本:%4d 传送脚本:%4d 对象索引:%4d\n",
				i, scene.MapNum, scene.ScriptOnEnter, scene.ScriptOnTeleport, scene.EventObjectIndex)
		}
	}
}

func printScene(scene mkf.Scene) {
	fmt.Printf("  地图编号:         %d\n", scene.MapNum)
	fmt.Printf("  进入脚本:         %d\n", scene.ScriptOnEnter)
	fmt.Printf("  传送脚本:         %d\n", scene.ScriptOnTeleport)
	fmt.Printf("  事件对象起始索引: %d\n", scene.EventObjectIndex)
}

func showSceneEnterScript(sss *mkf.SSSMkf, scriptChunk *mkf.ScriptEntryChunk, msgOffsetChunk *mkf.MsgOffsetChunk, msgData []byte, sceneIdx int) {
	sceneChunk, err := sss.GetSceneChunk()
	if err != nil {
		fmt.Printf("读取场景数据失败: %v\n", err)
		return
	}

	count := sceneChunk.Len()
	if sceneIdx >= count {
		fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", sceneIdx, count-1)
		return
	}

	scene := sceneChunk.GetScene(sceneIdx)
	if scene.ScriptOnEnter == 0 {
		fmt.Printf("场景 #%d 没有进入脚本\n", sceneIdx)
		return
	}

	fmt.Printf("\n=== 场景 #%d 的进入脚本 (索引: %d) ===\n", sceneIdx, scene.ScriptOnEnter)
	printScriptContent(scriptChunk, msgOffsetChunk, msgData, int(scene.ScriptOnEnter))
}

func showSceneTeleportScript(sss *mkf.SSSMkf, scriptChunk *mkf.ScriptEntryChunk, msgOffsetChunk *mkf.MsgOffsetChunk, msgData []byte, sceneIdx int) {
	sceneChunk, err := sss.GetSceneChunk()
	if err != nil {
		fmt.Printf("读取场景数据失败: %v\n", err)
		return
	}

	count := sceneChunk.Len()
	if sceneIdx >= count {
		fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", sceneIdx, count-1)
		return
	}

	scene := sceneChunk.GetScene(sceneIdx)
	if scene.ScriptOnTeleport == 0 {
		fmt.Printf("场景 #%d 没有传送脚本\n", sceneIdx)
		return
	}

	fmt.Printf("\n=== 场景 #%d 的传送脚本 (索引: %d) ===\n", sceneIdx, scene.ScriptOnTeleport)
	printScriptContent(scriptChunk, msgOffsetChunk, msgData, int(scene.ScriptOnTeleport))
}

func printScriptContent(scriptChunk *mkf.ScriptEntryChunk, msgOffsetChunk *mkf.MsgOffsetChunk, msgData []byte, startIndex int) {
	if scriptChunk == nil {
		fmt.Println("  无法读取脚本数据")
		return
	}

	count := scriptChunk.Len()
	if startIndex >= count {
		fmt.Printf("  脚本索引 %d 超出范围\n", startIndex)
		return
	}

	idx := startIndex
	for {
		if idx >= count {
			break
		}

		entry := scriptChunk.GetScriptEntry(idx)

		// 显示脚本指令
		fmt.Printf("    [%5d] 操作码: 0x%04X 操作数: [%d, %d, %d]  %s",
			idx, entry.Operation, entry.Operand[0], entry.Operand[1], entry.Operand[2],
			opcodeToString(entry.Operation))

		// 如果是显示文字操作码，显示文本内容
		if entry.Operation == 0xFFFF && msgOffsetChunk != nil && msgData != nil {
			msgID := entry.Operand[0]
			msgText := getMessageText(msgOffsetChunk, msgData, int(msgID))
			if msgText != "" {
				fmt.Printf("  文本: %s", msgText)
			}
		}

		fmt.Println()

		// 遇到停止指令退出
		if entry.Operation == OP_STOP {
			break
		}

		idx++
	}
}

// getMessageText 根据消息ID获取消息文本
func getMessageText(msgOffsetChunk *mkf.MsgOffsetChunk, msgData []byte, msgID int) string {
	if msgOffsetChunk == nil || msgData == nil {
		return ""
	}

	count := msgOffsetChunk.Len()
	if msgID < 0 || msgID >= count-1 {
		return ""
	}

	offsetCrt := *((*mkf.DWORD)(msgOffsetChunk.Get(msgID)))
	offsetNext := *((*mkf.DWORD)(msgOffsetChunk.Get(msgID + 1)))

	if int(offsetCrt) >= len(msgData) || int(offsetNext) > len(msgData) {
		return ""
	}

	// BIG5 编码解码
	mbs := msgData[offsetCrt:offsetNext]
	wcs := make([]rune, 0, len(mbs))
	state := 0
	for i := 0; i < len(mbs); i++ {
		v := mbs[i]
		if v == 0 {
			break
		}
		if state == 0 {
			if v <= 0x80 {
				wcs = append(wcs, rune(v))
			} else if v == 0xff {
				wcs = append(wcs, 0xf8f8)
			} else {
				state = 1
				continue
			}
		} else {
			if v < 0x40 || v >= 0x7f && v <= 0xa0 {
				wcs = append(wcs, 0x3f) // 无效字符
			} else if v <= 0x7e {
				wcs = append(wcs, utils.Cptbl_big5[mbs[i-1]-0x81][v-0x40])
			} else {
				wcs = append(wcs, utils.Cptbl_big5[mbs[i-1]-0x81][v-0x60])
			}
			state = 0
		}
	}

	// 移除末尾的 '1' 字符（游戏中用于标记）
	if len(wcs) > 0 && wcs[len(wcs)-1] == '1' {
		wcs = wcs[:len(wcs)-1]
	}

	return string(wcs)
}

func showEventObjects(sss *mkf.SSSMkf, idx int, listAll bool) {
	chunk, err := sss.GetEventObjectChunk()
	if err != nil {
		fmt.Printf("读取事件对象数据失败: %v\n", err)
		return
	}

	count := chunk.Len()

	if idx >= 0 {
		if idx >= count {
			fmt.Printf("错误: 事件对象索引 %d 超出范围（最大索引: %d）\n", idx, count-1)
			return
		}
		obj := chunk.GetEventObject(idx)
		fmt.Printf("\n=== 事件对象 #%d ===\n", idx)
		printEventObject(obj)
	} else if listAll {
		fmt.Println("\n=== 事件对象列表 ===")
		for i := 0; i < count; i++ {
			obj := chunk.GetEventObject(i)
			if obj.X == 0 && obj.Y == 0 && obj.SpriteNum == 0 && obj.TriggerScript == 0 && obj.AutoScript == 0 {
				continue
			}
			fmt.Printf("[%4d] X:%4d Y:%4d 图层:%2d 精灵:%4d 触发脚本:%4d 状态:%d\n",
				i, obj.X, obj.Y, obj.Layer, obj.SpriteNum, obj.TriggerScript, obj.State)
		}
	}
}

func printEventObject(obj mkf.EventObject) {
	fmt.Printf("  X 坐标:              %d\n", obj.X)
	fmt.Printf("  Y 坐标:              %d\n", obj.Y)
	fmt.Printf("  图层:                %d\n", obj.Layer)
	fmt.Printf("  触发脚本:            %d\n", obj.TriggerScript)
	fmt.Printf("  自动脚本:            %d\n", obj.AutoScript)
	fmt.Printf("  状态:                %d\n", obj.State)
	fmt.Printf("  触发模式:            %d\n", obj.TriggerMode)
	fmt.Printf("  精灵编号:            %d\n", obj.SpriteNum)
	fmt.Printf("  精灵帧数(触发):      %d\n", obj.SpriteFrames)
	fmt.Printf("  方向:                %s\n", directionToString(obj.Direction))
	fmt.Printf("  当前帧:              %d\n", obj.CurrentFrameNum)
	fmt.Printf("  精灵帧数(自动):      %d\n", obj.SpriteFramesAuto)
	fmt.Printf("  消失时间:            %d\n", obj.VanishTime)
}

func directionToString(dir mkf.WORD) string {
	switch dir {
	case kDirSouth:
		return "南"
	case kDirWest:
		return "西"
	case kDirNorth:
		return "北"
	case kDirEast:
		return "东"
	default:
		return fmt.Sprintf("未知(%d)", dir)
	}
}

func showScripts(sss *mkf.SSSMkf) {
	chunk, err := sss.GetScriptEntryChunk()
	if err != nil {
		fmt.Printf("读取脚本数据失败: %v\n", err)
		return
	}

	fmt.Println("\n=== 脚本条目列表 ===")
	count := chunk.Len()
	for i := 0; i < count; i++ {
		entry := chunk.GetScriptEntry(i)
		if entry.Operation == 0 && entry.Operand[0] == 0 && entry.Operand[1] == 0 && entry.Operand[2] == 0 {
			continue
		}
		fmt.Printf("[%5d] 操作码: 0x%04X 操作数: [%d, %d, %d]  %s\n",
			i, entry.Operation, entry.Operand[0], entry.Operand[1], entry.Operand[2],
			opcodeToString(entry.Operation))
	}
}
