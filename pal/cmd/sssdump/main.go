package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/njuwelkin/games/pal/pkg/game"
	"github.com/njuwelkin/games/pal/pkg/mkf"
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
	gamePath            = flag.String("gamepath", ".", "游戏数据目录路径")
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
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -all\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -scenes\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -scene 10\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -objects\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -object 5\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -scene-enter-script 10\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -gamepath /path/to/game -scene-teleport-script 10\n", os.Args[0])
	}

	flag.Parse()

	if *gamePath != "" {
		os.Setenv("PAL_GAME_PATH", *gamePath)
	}

	game.InitGlobalSetting()

	if *listAll || *listScenes || *sceneIndex >= 0 {
		showScenes(*sceneIndex, *listAll || *listScenes)
	}

	if *listAll || *listObjects || *objectIndex >= 0 {
		showEventObjects(*objectIndex, *listAll || *listObjects)
	}

	if *listAll || *listScripts {
		showScripts()
	}

	if *sceneEnterScript >= 0 {
		showSceneEnterScript(*sceneEnterScript)
	}

	if *sceneTeleportScript >= 0 {
		showSceneTeleportScript(*sceneTeleportScript)
	}

	if !*listAll && !*listScenes && !*listObjects && !*listScripts &&
		*sceneIndex < 0 && *objectIndex < 0 && *sceneEnterScript < 0 && *sceneTeleportScript < 0 {
		flag.Usage()
	}
}

func showScenes(idx int, listAll bool) {
	scenes := game.Globals.G.Scenes

	if idx >= 0 {
		if idx >= len(scenes) {
			fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", idx, len(scenes)-1)
			return
		}
		scene := scenes[idx]
		if scene.MapNum == 0 && scene.ScriptOnEnter == 0 &&
			scene.ScriptOnTeleport == 0 && scene.EventObjectIndex == 0 {
			fmt.Printf("场景 #%d 为空\n", idx)
			return
		}
		fmt.Printf("\n=== 场景 #%d ===\n", idx)
		printScene(scene)

		if scene.ScriptOnEnter != 0 {
			fmt.Println("\n  --- 进入脚本内容 ---")
			printScriptContent(int(scene.ScriptOnEnter))
		}
		if scene.ScriptOnTeleport != 0 {
			fmt.Println("\n  --- 传送脚本内容 ---")
			printScriptContent(int(scene.ScriptOnTeleport))
		}
	} else if listAll {
		fmt.Println("\n=== 场景列表 ===")
		for i := 0; i < len(scenes); i++ {
			scene := scenes[i]
			if scene.MapNum == 0 && scene.ScriptOnEnter == 0 &&
				scene.ScriptOnTeleport == 0 && scene.EventObjectIndex == 0 {
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

func showSceneEnterScript(sceneIdx int) {
	scenes := game.Globals.G.Scenes

	if sceneIdx >= len(scenes) {
		fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", sceneIdx, len(scenes)-1)
		return
	}

	scene := scenes[sceneIdx]
	if scene.ScriptOnEnter == 0 {
		fmt.Printf("场景 #%d 没有进入脚本\n", sceneIdx)
		return
	}

	fmt.Printf("\n=== 场景 #%d 的进入脚本 (索引: %d) ===\n", sceneIdx, scene.ScriptOnEnter)
	printScriptContent(int(scene.ScriptOnEnter))
}

func showSceneTeleportScript(sceneIdx int) {
	scenes := game.Globals.G.Scenes

	if sceneIdx >= len(scenes) {
		fmt.Printf("错误: 场景索引 %d 超出范围（最大索引: %d）\n", sceneIdx, len(scenes)-1)
		return
	}

	scene := scenes[sceneIdx]
	if scene.ScriptOnTeleport == 0 {
		fmt.Printf("场景 #%d 没有传送脚本\n", sceneIdx)
		return
	}

	fmt.Printf("\n=== 场景 #%d 的传送脚本 (索引: %d) ===\n", sceneIdx, scene.ScriptOnTeleport)
	printScriptContent(int(scene.ScriptOnTeleport))
}

func printScriptContent(startIndex int) {
	scriptEntries := game.Globals.G.ScriptEntries

	count := len(scriptEntries)
	if startIndex >= count {
		fmt.Printf("  脚本索引 %d 超出范围\n", startIndex)
		return
	}

	idx := startIndex
	for {
		if idx >= count {
			break
		}

		entry := scriptEntries[idx]

		fmt.Printf("    [%5d] 操作码: 0x%04X 操作数: [%d, %d, %d]  %s",
			idx, entry.Operation, entry.Operand[0], entry.Operand[1], entry.Operand[2],
			opcodeToString(entry.Operation))

		if entry.Operation == 0xFFFF {
			msgID := entry.Operand[0]
			msgText := game.Globals.Text.GetMessage(int(msgID))
			if msgText != "" {
				fmt.Printf("  文本: %s", msgText)
			}
		}

		fmt.Println()

		if entry.Operation == OP_STOP {
			break
		}

		idx++
	}
}

func showEventObjects(idx int, listAll bool) {
	objects := game.Globals.G.EventObjects

	count := len(objects)

	if idx >= 0 {
		if idx >= count {
			fmt.Printf("错误: 事件对象索引 %d 超出范围（最大索引: %d）\n", idx, count-1)
			return
		}
		obj := objects[idx]
		fmt.Printf("\n=== 事件对象 #%d ===\n", idx)
		printEventObject(obj)
	} else if listAll {
		fmt.Println("\n=== 事件对象列表 ===")
		for i := 0; i < count; i++ {
			obj := objects[i]
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

func showScripts() {
	scriptEntries := game.Globals.G.ScriptEntries

	fmt.Println("\n=== 脚本条目列表 ===")
	count := len(scriptEntries)
	for i := 0; i < count; i++ {
		entry := scriptEntries[i]
		if entry.Operation == 0 && entry.Operand[0] == 0 && entry.Operand[1] == 0 && entry.Operand[2] == 0 {
			continue
		}
		fmt.Printf("[%5d] 操作码: 0x%04X 操作数: [%d, %d, %d]  %s\n",
			i, entry.Operation, entry.Operand[0], entry.Operand[1], entry.Operand[2],
			opcodeToString(entry.Operation))
	}
}
