package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/ui"
)

// ScriptExecutor 脚本执行器，支持分步执行脚本
type ScriptExecutor struct {
	owner *sceneScreen

	scriptEntry       mkf.WORD         // 当前脚本入口
	eventObjID        mkf.WORD         // 事件对象ID
	lastEventObjectID mkf.WORD         // 上一个事件对象ID
	pEvtObj           *mkf.EventObject // 事件对象指针
	nextScriptEntry   mkf.WORD         // 下一个脚本入口
	ended             bool             // 是否执行结束
	suspended         bool             // 是否暂停（遇到0x0005或0x004D）
	dialog            *ui.Dialog       // 当前显示的对话框
}

// NewScriptExecutor 创建新的脚本执行器
func NewScriptExecutor(owner *sceneScreen) *ScriptExecutor {
	return &ScriptExecutor{
		owner: owner,

		scriptEntry:       0,
		eventObjID:        0,
		lastEventObjectID: 0,
		pEvtObj:           nil,
		nextScriptEntry:   0,
		ended:             false,
		suspended:         false,
		dialog:            nil,
	}
}

// RunTriggerScript 开始执行触发脚本（分步执行版本）
// 返回 true 表示脚本执行完成，false 表示暂停需要继续调用 ContinueRun
func (se *ScriptExecutor) RunTriggerScript(scriptEntry mkf.WORD, eventObjID mkf.WORD) bool {
	g := &globals.G

	// 初始化状态
	se.scriptEntry = scriptEntry
	se.eventObjID = eventObjID
	se.nextScriptEntry = scriptEntry
	se.ended = false
	se.suspended = false

	globals.UpdatedInBattle = false

	if eventObjID == 0xffff {
		se.eventObjID = se.lastEventObjectID
	}

	se.lastEventObjectID = se.eventObjID

	if se.eventObjID != 0 {
		se.pEvtObj = &g.eventObjects[se.eventObjID-1]
	}

	globals.ScriptSuccess = true

	// 执行脚本（分步执行）
	return se.executeStep()
}

// ContinueRun 继续执行脚本
// 返回 true 表示脚本执行完成，false 表示需要继续调用
func (se *ScriptExecutor) ContinueRun() bool {
	if se.ended {
		return true
	}
	se.suspended = false
	return se.executeStep()
}

// IsSuspended 检查是否处于暂停状态
func (se *ScriptExecutor) IsSuspended() bool {
	return se.suspended
}

// IsEnded 检查是否执行结束
func (se *ScriptExecutor) IsEnded() bool {
	return se.ended
}

// executeStep 执行脚本步骤
func (se *ScriptExecutor) executeStep() bool {
	g := &globals.G

	for se.scriptEntry != 0 && !se.ended && !se.suspended {
		pScript := &g.scriptEntries[se.scriptEntry]

		switch pScript.Operation {
		case 0x0000:
			// 停止执行
			se.ended = true

		case 0x0001:
			// 停止执行并替换为下一行
			se.ended = true
			se.nextScriptEntry = se.scriptEntry + 1

		case 0x0002:
			// 停止执行并替换为指定行
			if pScript.Operand[1] == 0 {
				se.ended = true
				se.nextScriptEntry = pScript.Operand[0]
			} else {
				se.pEvtObj.ScriptIdleFrame++
				if se.pEvtObj.ScriptIdleFrame < pScript.Operand[1] {
					se.ended = true
					se.nextScriptEntry = pScript.Operand[0]
				} else {
					se.pEvtObj.ScriptIdleFrame = 0
					se.scriptEntry++
				}
			}

		case 0x0003:
			// 无条件跳转
			if pScript.Operand[1] == 0 {
				se.scriptEntry = pScript.Operand[0]
			} else {
				se.pEvtObj.ScriptIdleFrame++
				if se.pEvtObj.ScriptIdleFrame < pScript.Operand[1] {
					se.scriptEntry = pScript.Operand[0]
				} else {
					se.pEvtObj.ScriptIdleFrame = 0
					se.scriptEntry++
				}
			}

		case 0x0004:
			// 调用子脚本
			newEvtObjID := se.eventObjID
			if pScript.Operand[1] != 0 {
				newEvtObjID = pScript.Operand[1]
			}
			runTriggerScript(pScript.Operand[0], newEvtObjID)
			se.scriptEntry++

		case 0x0005:
			// 重绘屏幕 - 暂停等待屏幕更新
			se.owner.Notify(-1, ui.WaitForKey, nil)
			se.suspended = true
			se.scriptEntry++

		case 0x0006:
			// 按概率跳转到指定地址
			if rand.Intn(100)+1 >= int(pScript.Operand[0]) {
				se.scriptEntry = pScript.Operand[1]
				continue
			} else {
				se.scriptEntry++
			}

		case 0x0007:
			// 开始战斗（未实现）
			se.scriptEntry++

		case 0x0008:
			// 替换为下一条指令
			se.scriptEntry++
			se.nextScriptEntry = se.scriptEntry

		case 0x0009:
			// wait for the specified number of frames
			// Note: Need to implement
			se.scriptEntry++

		case 0x000A:
			// Goto the specified address if player selected no
			se.scriptEntry++

		case 0x003B:
			//
			// Show dialog in the middle part of the screen
			//
			dialog := ui.NewDialog(0, 0, 100, 300, se.owner, nil, globals.Font.NormalFont)
			se.owner.AddComponent(&dialog)
			se.dialog = &dialog
			se.scriptEntry++
		case 0x003C:
			//
			// Show dialog in the upper part of the screen
			//
			fontColor := pScript.Operand[1]
			numCharFace := pScript.Operand[0]
			playingRNG := pScript.Operand[2] != 0
			_, _, _ = fontColor, numCharFace, playingRNG
			dialog := createDialog(se.owner, fontColor, numCharFace, playingRNG)
			se.owner.AddComponent(dialog)
			se.dialog = dialog
			se.scriptEntry++
		case 0x003D:
			//
			// Show dialog in the lower part of the screen
			//
			dialog := ui.NewDialog(0, 0, 200, 300, se.owner, nil, globals.Font.NormalFont)
			se.owner.AddComponent(&dialog)
			se.dialog = &dialog
			se.scriptEntry++
		case 0x003E:
			//
			// Show text in a window at the center of the screen
			//
			/*
				PAL_ClearDialog(TRUE);
				PAL_StartDialog(kDialogCenterWindow, (BYTE)pScript->rgwOperand[0], 0, FALSE);
			*/
			se.scriptEntry++
		case 0xFFFF:
			//
			// Show text
			//
			textNum := pScript.Operand[0]
			text := []rune{}
			if textNum <= mkf.WORD(globals.Text.CountMsgs) {
				text = globals.Text.MsgBuf[textNum]
			}
			se.dialog.AppendLine(text)
			se.scriptEntry++

		case 0x004D:
			// 等待按键 - 暂停等待用户输入
			se.owner.Notify(-1, ui.WaitForKey, nil)
			se.suspended = true

		default:
			// 未实现或未知操作码，继续执行下一条
			se.scriptEntry = interpretInstruction(se.scriptEntry, se.eventObjID)
		}
	}

	return se.ended
}

func createDialog(owner ui.ParentCom, fontColor, numCharFace mkf.WORD, playingRNG bool) *ui.Dialog {
	var avatarImg *ebiten.Image = nil
	if numCharFace > 0 {
		rgm, err := mkf.NewRgmMkf("../../RGM.MKF")
		if err != nil {
			panic(err.Error())
		}
		defer rgm.Close()
		bmp, err := rgm.GetFaceBmp(mkf.INT(numCharFace))
		if err != nil || bmp == nil {
			panic(err.Error())
		}
		plt, err := mkf.GetPalette(mkf.INT(globals.G.crtPaletteNum), false)
		if err != nil {
			panic(err.Error())
		}
		avatarImg = bmp.ToImageWithPalette(plt)
	}
	dialog := ui.NewDialog(0, 0, 200, 300, owner, avatarImg, globals.Font.NormalFont)

	return &dialog
}
