package main

import (
	"unsafe"

	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/utils"
)

func runTriggerScript(scriptEntry mkf.WORD, eventObjID mkf.WORD) {
	g := &globals.G

	var lastEventObjectID mkf.WORD = 0
	nextScriptEntry := scriptEntry
	if eventObjID == 0xffff {
		eventObjID = lastEventObjectID
	}
	lastEventObjectID = eventObjID

	var pEvtObj *mkf.EventObject
	if eventObjID != 0 {
		pEvtObj = &g.eventObjects[eventObjID-1]
	}

	ended := false
	for scriptEntry != 0 && !ended {
		pScript := &g.scriptEntries[scriptEntry]
		switch pScript.Operation {
		case 0x0000:
			//
			// Stop running
			//
			ended = true
		case 0x0001:
			//
			// Stop running and replace the entry with the next line
			//
			ended = true
			nextScriptEntry = scriptEntry + 1

		case 0x0002:
			//
			// Stop running and replace the entry with the specified one
			//
			pEvtObj.ScriptIdleFrame++
			if pScript.Operand[1] == 0 ||
				pEvtObj.ScriptIdleFrame < pScript.Operand[1] {
				ended = true
				nextScriptEntry = pScript.Operand[0]
			} else {
				//
				// failed
				//
				pEvtObj.ScriptIdleFrame = 0
				scriptEntry++
			}
		default:
			// PAL_ClearDialog(TRUE);
			scriptEntry = interpretInstruction(scriptEntry, eventObjID)
		}
	}
}

func interpretInstruction(scriptEntry mkf.WORD, eventObjID mkf.WORD) mkf.WORD {
	g := &globals.G

	pScript := g.scriptEntries[scriptEntry]

	var pEvtObj *mkf.EventObject
	if eventObjID != 0 {
		pEvtObj = &g.eventObjects[eventObjID-1]
	}

	var pCurrent *mkf.EventObject
	var curEventObjectID mkf.WORD
	if pScript.Operand[0] == 0 || pScript.Operand[0] == 0xFFFF {
		pCurrent = pEvtObj
		curEventObjectID = eventObjID
	} else {
		i := pScript.Operand[0] - 1
		if i > 0x9000 {
			// HACK for Dream 2.11 to avoid crash
			i -= 0x9000
		}
		pCurrent = &(g.eventObjects[i])
		curEventObjectID = pScript.Operand[0]
	}

	var iPlayerRole mkf.WORD
	if pScript.Operand[0] < mkf.MAX_PLAYABLE_PLAYER_ROLES {
		iPlayerRole = g.parties[pScript.Operand[0]].PlayerRole
	} else {
		iPlayerRole = g.parties[0].PlayerRole
	}

	switch pScript.Operation {
	case 0x000B:
	case 0x000C:
	case 0x000D:
	case 0x000E:
		//
		// walk one step
		//
		pEvtObj.Direction = pScript.Operation - 0x000B
		npcWalkOneStep(eventObjID, 2)
	case 0x0010:
		//
		// Walk straight to the specified position
		//
		if !npcWalkTo(eventObjID, mkf.INT(pScript.Operand[0]), mkf.INT(pScript.Operand[1]), mkf.INT(pScript.Operand[2]), 3) {
			scriptEntry--
		}
	case 0x0011:
		//
		// Walk straight to the specified position, at a lower speed
		//
		if mkf.DWORD(eventObjID&1)^(g.frameNum&1) != 0 {
			if !npcWalkTo(eventObjID, mkf.INT(pScript.Operand[0]), mkf.INT(pScript.Operand[1]), mkf.INT(pScript.Operand[2]), 2) {
				scriptEntry--
			}
		} else {
			scriptEntry--
		}
	case 0x0012:
		//
		// Set the position of the event object, relative to the party
		//
		pCurrent.X = pScript.Operand[1] + uint16(g.viewport.X()) + uint16(g.partyoffset.X())
		pCurrent.Y = pScript.Operand[2] + uint16(g.viewport.Y()) + uint16(g.partyoffset.Y())
	case 0x0013:
		//
		// Set the position of the event object
		//
		pCurrent.X = pScript.Operand[1]
		pCurrent.Y = pScript.Operand[2]
	case 0x0014:
		//
		// Set the gesture of the event object
		//
		pEvtObj.CurrentFrameNum = pScript.Operand[0]
		pEvtObj.Direction = kDirSouth
	case 0x0015:
		//
		// Set the direction and gesture for a party member
		//
		g.partyDirection = pScript.Operand[0]
		g.parties[pScript.Operand[2]].Frame = g.partyDirection*3 + pScript.Operand[1]
	case 0x0016:
		//
		// Set the direction and gesture for an event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.Direction = pScript.Operand[1]
			pCurrent.CurrentFrameNum = pScript.Operand[2]
		}
	case 0x0017:
		//
		// set the player's extra attribute
		//
		i := pScript.Operand[0] - 0xB
		base := (*mkf.WORD)(unsafe.Pointer(&g.equipmentEffect[i])) // HACKHACK
		p := (*mkf.WORD)(unsafe.Pointer(uintptr(unsafe.Pointer(base)) + uintptr(pScript.Operand[1]*mkf.MAX_PLAYER_ROLES+eventObjID)))
		*p = mkf.WORD(pScript.Operand[2])
	case 0x0018:
		//
		// Equip the selected item
		//
		i = pScript.Operand[0] - 0x0B
		g.curEquipPart = i

		//
		// The wEventObjectID parameter here should indicate the player role
		//
		//PAL_RemoveEquipmentEffect(wEventObjectID, i);

		if g.playerRoles.Equipment[i][eventObjID] != pScript.Operand[1] {
			w := g.playerRoles.Equipment[i][eventObjID]
			g.playerRoles.Equipment[i][eventObjID] = pScript.Operand[1]

			i, foundI := getItemIndexToInventory(pScript.Operand[1])
			_, foundJ := getItemIndexToInventory(w)
			if foundI && i < MAX_INVENTORY && g.inventory[i].Amount == 1 && w != 0 && !foundJ {
				//
				// When the number of items you want to wear is 1
				// and the number of items you are wearing is also 1,
				// replace them directly, instead of removing items
				// and adding them at the end of the item menu
				//
				g.inventory[i].Item = w
			} else {
				addItemIndexToInventory(pScript.Operand[1], -1)

				if w != 0 {
					addItemIndexToInventory(w, 1)
				}
			}

			g.lastUnequippedItem = w
		}
	case 0x0019:
		//
		// Increase/decrease the player's attribute
		//

		// TODO: use raw
		var iPlayerRole mkf.WORD
		if pScript.Operand[2] == 0 {
			iPlayerRole = eventObjID
		} else {
			iPlayerRole = pScript.Operand[2] - 1
		}
		base := (*mkf.WORD)(unsafe.Pointer(&g.playerRoles))
		p := (*mkf.WORD)(unsafe.Pointer(uintptr(unsafe.Pointer(base)) + uintptr(pScript.Operand[0]*mkf.MAX_PLAYER_ROLES+iPlayerRole)))
		*p += mkf.WORD(pScript.Operand[1])
	case 0x001A:
		//
		// Set player's stat
		//
		p := utils.WordArray(unsafe.Pointer(&g.playerRoles), unsafe.Sizeof(g.playerRoles))
		if g.curEquipPart != -1 {
			//
			// In the progress of equipping items
			//
			p = utils.WordArray(unsafe.Pointer(&g.equipmentEffect[g_iCurEquipPart]), unsafe.Sizeof(g.equipmentEffect[g_iCurEquipPart]))
		}
		if pScript.Operand[2] == 0 {
			//
			// Apply to the current player. The wEventObjectID should
			// indicate the player role.
			//
			iPlayerRole = eventObjID
		} else {
			iPlayerRole = pScript.Operand[2] - 1
		}
		p[pScript.Operand[0]*mkf.MAX_PLAYER_ROLES+iPlayerRole] = mkf.WORD(pScript.Operand[1])
	case 0x001B:
		//
		// Increase/decrease player's HP
		//
		if pScript.Operand[0] != 0 {
			g.scriptSuccess = false
			//
			// Apply to everyone
			//
			for i := uint16(0); i <= g.maxPartyMemberIndex; i++ {
				w := g.parties[i].PlayerRole
				if increaseHPMP(w, (mkf.SHORT)(pScript.Operand[1]), 0) {
					g.scriptSuccess = true
				}
			}
		} else {
			//
			// Apply to one player. The wEventObjectID parameter should indicate the player role.
			//
			if !increaseHPMP(eventObjID, (mkf.SHORT)(pScript.Operand[1]), 0) {
				g.scriptSuccess = false
			}
		}
	case 0x001C:
		//
		// Increase/decrease player's MP
		//
	case 0x001D:
		//
		// Increase/decrease player's HP and MP
		//
	case 0x001E:
		//
		// Increase or decrease cash by the specified amount
		//
	case 0x001F:
		//
		// Add item to inventory
		//
	case 0x0020:
		//
		// Remove item from inventory
		//
	case 0x0021:
		//
		// Inflict damage to the enemy
		//
	case 0x0022:
		//
		// Revive player
		//
	case 0x0023:
		//
		// Remove equipment from the specified player
		//
	case 0x0024:
		//
		// Set the autoscript entry address for an event object
		//
	case 0x0025:
		//
		// Set the trigger script entry address for an event object
		//
	case 0x0026:
		//
		// Show the buy item menu
		//
	case 0x0027:
		//
		// Show the sell item menu
		//
	case 0x0028:
		//
		// Apply poison to enemy
		//
	case 0x0029:
		//
		// Apply poison to player
		//
	case 0x002A:
		//
		// Cure poison by object ID for enemy
		//
	case 0x002B:
		//
		// Cure poison by object ID for player
		//
	case 0x002C:
		//
		// Cure poisons by level
		//
	case 0x002D:
		//
		// Set the status for player
		//
	case 0x002E:
		//
		// Set the status for enemy
		//
	case 0x002F:
		//
		// Remove player's status
		//
	case 0x0030:
		//
		// Increase player's stat temporarily by percent
		//
	case 0x0031:
		//
		// Change battle sprite temporarily for player
		//
	case 0x0033:
		//
		// collect the enemy for items
		//
	case 0x0034:
		//
		// Transform collected enemies into items
		//
	case 0x0035:
		//
		// Shake the screen
		//
	case 0x0036:
		//
		// Set the current playing RNG animation
		//
	case 0x0037:
		//
		// Play RNG animation
		//
	case 0x0038:
		//
		// Teleport the party out of the scene
		//
	case 0x0039:
		//
		// Drain HP from enemy
		//
	case 0x003A:
		//
		// Player flee from the battle
		//
	case 0x003F:
		//
		// Ride the event object to the specified position, at a low speed
		//
	case 0x0040:
		//
		// set the trigger method for a event object
		//
	case 0x0041:
		//
		// Mark the script as failed
		//
		g.scriptSuccess = false
	case 0x0042:
		//
		// Simulate a magic for player
		//
	case 0x0043:
		//
		// Set background music
		//
	case 0x0044:
		//
		// Ride the event object to the specified position, at the normal speed
		//
	case 0x0045:
		//
		// Set battle music
		//
	case 0x0046:
		//
		// Set the party position on the map
		//
		{
			var xOffset, yOffset int
			var x, y int

			if g.partyDirection == kDirWest || g.partyDirection == kDirSouth {
				xOffset = 16
			} else {
				xOffset = -16
			}
			if g.partyDirection == kDirWest || g.partyDirection == kDirNorth {
				yOffset = 8
			} else {
				yOffset = -8
			}

			x = int(pScript.Operand[0]*32 + pScript.Operand[2]*16)
			y = int(pScript.Operand[1]*16 + pScript.Operand[2]*8)

			x -= int(g.partyoffset.X())
			y -= int(g.partyoffset.Y())

			g.viewport = PAL_XY(uint32(x), uint32(y))

			x = int(g.partyoffset.X())
			y = int(g.partyoffset.Y())

			for i := 0; i < mkf.MAX_PLAYABLE_PLAYER_ROLES; i++ {
				g.parties[i].X = mkf.SHORT(x)
				g.parties[i].Y = mkf.SHORT(y)
				g.trails[i].X = mkf.WORD(uint32(x) + g.viewport.X())
				g.trails[i].Y = mkf.WORD(uint32(y) + g.viewport.Y())
				g.trails[i].Direction = g.partyDirection

				x += xOffset
				y += yOffset
			}
		}
	case 0x0065:
		//
		// Set the player's sprite
		//
		g.playerRoles.SpriteNum[pScript.Operand[0]] = pScript.Operand[1]
		if !g.inBattle && pScript.Operand[2] != 0 {
			//PAL_SetLoadFlags(kLoadPlayerSprite);
			//PAL_LoadResources();
		}
	case 0x0075:
		//
		// Set the player party
		//
		g.maxPartyMemberIndex = 0
		for i := 0; i < 3; i++ {
			if pScript.Operand[i] != 0 {
				g.parties[g.maxPartyMemberIndex].PlayerRole =
					pScript.Operand[i] - 1

				g.maxPartyMemberIndex++
			}
		}

		if g.maxPartyMemberIndex == 0 {
			// HACK for Dream 2.11
			g.parties[0].PlayerRole = 0
			g.maxPartyMemberIndex = 1
		}

		g.maxPartyMemberIndex--

		//
		// Reload the player sprites
		//
		//PAL_SetLoadFlags(kLoadPlayerSprite);
		//PAL_LoadResources();

		//memset(gpGlobals->rgPoisonStatus, 0, sizeof(gpGlobals->rgPoisonStatus));
		//PAL_UpdateEquipments();
	}

	return 0
}

func npcWalkOneStep(evtObjID mkf.WORD, speed mkf.INT) {

}

func npcWalkTo(evtObjID mkf.WORD, x, y, h, speed mkf.INT) bool {
	return false
}

func increaseHPMP(playerRole mkf.WORD, HP mkf.SHORT, MP mkf.SHORT) bool {
	return false
}

func getItemIndexToInventory(objectID mkf.WORD) (mkf.INT, bool) {
	found := false

	var index mkf.INT = 0

	for index < MAX_INVENTORY {
		if globals.G.inventory[index].Item == objectID {
			found = true
			break
		} else if globals.G.inventory[index].Item == 0 {
			break
		}
		index++
	}

	return index, found
}

func addItemIndexToInventory(objectID mkf.WORD, num int) bool {
	return false
}
