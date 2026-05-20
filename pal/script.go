package main

import (
	"fmt"
	"math/rand"
	"unsafe"

	"github.com/njuwelkin/games/pal/mkf"
	"github.com/njuwelkin/games/pal/utils"
)

/*
++

	Purpose:
	  Runs a trigger script.
	Parameters:
	  [IN]  wScriptEntry - The script entry to execute.
	  [IN]  wEventObjectID - The event object ID which invoked the script.
	Return value:
	  The entry point of the script.

--
*/
func runTriggerScript(scriptEntry mkf.WORD, eventObjID mkf.WORD) {
	g := &globals.G

	var lastEventObjectID mkf.WORD = 0
	nextScriptEntry := scriptEntry
	globals.UpdatedInBattle = false

	if eventObjID == 0xffff {
		eventObjID = lastEventObjectID
	}

	lastEventObjectID = eventObjID
	var pEvtObj *mkf.EventObject
	if eventObjID != 0 {
		pEvtObj = &g.eventObjects[eventObjID-1]
	}

	globals.ScriptSuccess = true

	//
	// Set the default dialog speed.
	//
	//PAL_DialogSetDelayTime(3);

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
			if pScript.Operand[1] == 0 {
				ended = true
				nextScriptEntry = pScript.Operand[0]
			} else {
				pEvtObj.ScriptIdleFrame++
				if pEvtObj.ScriptIdleFrame < pScript.Operand[1] {
					ended = true
					nextScriptEntry = pScript.Operand[0]
				} else {
					//
					// failed
					//
					pEvtObj.ScriptIdleFrame = 0
					scriptEntry++
				}
			}
		case 0x0003:
			//
			// unconditional jump
			//
			if pScript.Operand[1] == 0 {
				scriptEntry = pScript.Operand[0]
			} else {
				pEvtObj.ScriptIdleFrame++
				if pEvtObj.ScriptIdleFrame < pScript.Operand[1] {
					scriptEntry = pScript.Operand[0]
				} else {
					//
					// failed
					//
					pEvtObj.ScriptIdleFrame = 0
					scriptEntry++
				}
			}
		case 0x0004:
			//
			// Call script
			//
			newEvtObjID := eventObjID
			if pScript.Operand[1] != 0 {
				newEvtObjID = pScript.Operand[1]
			}
			runTriggerScript(pScript.Operand[0], newEvtObjID)
			scriptEntry++
		case 0x0005:
			//
			// Redraw screen
			//
			/*
				PAL_ClearDialog(TRUE);

				if (PAL_DialogIsPlayingRNG())
				{
				   VIDEO_RestoreScreen(gpScreen);
				}
				else if (gpGlobals->fInBattle)
				{
				   PAL_BattleMakeScene();
				   VIDEO_CopyEntireSurface(g_Battle.lpSceneBuf, gpScreen);
				   VIDEO_UpdateScreen(NULL);
				}
				else
				{
				   if (pScript->rgwOperand[2])
				   {
					  PAL_UpdatePartyGestures(FALSE);
				   }

				   PAL_MakeScene();

				   VIDEO_UpdateScreen(NULL);
				   UTIL_Delay((pScript->rgwOperand[1] == 0) ? 60 : (pScript->rgwOperand[1] * 60));
				}
			*/
			scriptEntry++
		case 0x0006:
			//
			// Jump to the specified address by the specified rate
			//
			if rand.Intn(100)+1 >= int(pScript.Operand[0]) {
				scriptEntry = pScript.Operand[1]
				continue
			} else {
				scriptEntry++
			}
		case 0x0007:
			//
			// Start battle
			//
			/*
				i = PAL_StartBattle(pScript->rgwOperand[0], !pScript->rgwOperand[2]);

				if (i == kBattleResultLost && pScript->rgwOperand[1] != 0)
				{
				   wScriptEntry = pScript->rgwOperand[1];
				}
				else if (i == kBattleResultFleed && pScript->rgwOperand[2] != 0)
				{
				   wScriptEntry = pScript->rgwOperand[2];
				}
				else
				{
				   wScriptEntry++;
				}
				gpGlobals->fAutoBattle = FALSE;
			*/
			scriptEntry++
		case 0x0008:
			//
			// Replace the entry with the next instruction
			//
			scriptEntry++
			nextScriptEntry = scriptEntry
		case 0x0009:
			//
			// wait for the specified number of frames
			//
			/*
				{
				   DWORD        time;

				   PAL_ClearDialog(TRUE);

				   time = SDL_GetTicks() + FRAME_TIME;

				   for (i = 0; i < (pScript->rgwOperand[0] ? pScript->rgwOperand[0] : 1); i++)
				   {
					  PAL_DelayUntil(time);

					  time = SDL_GetTicks() + FRAME_TIME;

					  if (pScript->rgwOperand[2])
					  {
						 PAL_UpdatePartyGestures(FALSE);
					  }

					  PAL_GameUpdate(pScript->rgwOperand[1] ? TRUE : FALSE);
					  PAL_MakeScene();
					  VIDEO_UpdateScreen(NULL);
				   }
				}
			*/
			scriptEntry++
		case 0x000A:
			//
			// Goto the specified address if player selected no
			//
			/*
				PAL_ClearDialog(FALSE);

				if (!PAL_ConfirmMenu())
				{
				   wScriptEntry = pScript->rgwOperand[0];
				}
				else
				{
				   wScriptEntry++;
				}
			*/
			scriptEntry++
		case 0x003B:
			//
			// Show dialog in the middle part of the screen
			//
			/*
				PAL_ClearDialog(TRUE);
				PAL_StartDialog(kDialogCenter, (BYTE)pScript->rgwOperand[0], 0,
				   pScript->rgwOperand[2] ? TRUE : FALSE);
			*/
			scriptEntry++
		case 0x003C:
			//
			// Show dialog in the upper part of the screen
			//
			/*
				PAL_ClearDialog(TRUE);
				PAL_StartDialog(kDialogUpper, (BYTE)pScript->rgwOperand[1],
				   pScript->rgwOperand[0], pScript->rgwOperand[2] ? TRUE : FALSE);
			*/
			scriptEntry++
		case 0x003D:
			//
			// Show dialog in the lower part of the screen
			//
			/*
				PAL_ClearDialog(TRUE);
				PAL_StartDialog(kDialogLower, (BYTE)pScript->rgwOperand[1],
				   pScript->rgwOperand[0], pScript->rgwOperand[2] ? TRUE : FALSE);
			*/
			scriptEntry++
		case 0x003E:
			//
			// Show text in a window at the center of the screen
			//
			/*
				PAL_ClearDialog(TRUE);
				PAL_StartDialog(kDialogCenterWindow, (BYTE)pScript->rgwOperand[0], 0, FALSE);
			*/
			scriptEntry++
		case 0xFFFF:
			//
			// Show text
			//
			scriptEntry++
		default:
			// PAL_ClearDialog(TRUE);
			scriptEntry = interpretInstruction(scriptEntry, eventObjID)
		}
	}
	fmt.Print(nextScriptEntry)
}

/*
++

	Purpose:

	  Interpret and execute one instruction in the script.

	Parameters:

	  [IN]  wScriptEntry - The script entry to execute.

	  [IN]  wEventObjectID - The event object ID which invoked the script.

	Return value:

	  The address of the next script instruction to execute.

--
*/
func interpretInstruction(scriptEntry mkf.WORD, eventObjID mkf.WORD) mkf.WORD {
	g := &globals.G

	pScript := g.scriptEntries[scriptEntry]

	var pEvtObj *mkf.EventObject
	if eventObjID != 0 {
		pEvtObj = &g.eventObjects[eventObjID-1]
	}

	var pCurrent *mkf.EventObject
	//var curEventObjectID mkf.WORD
	if pScript.Operand[0] == 0 || pScript.Operand[0] == 0xFFFF {
		pCurrent = pEvtObj
		//curEventObjectID = eventObjID
	} else {
		i := pScript.Operand[0] - 1
		if i > 0x9000 {
			// HACK for Dream 2.11 to avoid crash
			i -= 0x9000
		}
		pCurrent = &(g.eventObjects[i])
		//curEventObjectID = pScript.Operand[0]
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
	case 0x000F:
		//
		// Set the direction and/or gesture for event object
		//
		if pScript.Operand[0] != 0xFFFF {
			pEvtObj.Direction = pScript.Operand[0]
		}
		if pScript.Operand[1] != 0xFFFF {
			pEvtObj.CurrentFrameNum = pScript.Operand[1]
		}
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
		i := pScript.Operand[0] - 0x0B
		g.curEquipPart = int(i)

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
			p = utils.WordArray(unsafe.Pointer(&g.equipmentEffect[g.curEquipPart]), unsafe.Sizeof(g.equipmentEffect[g.curEquipPart]))
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
		if pScript.Operand[0] != 0 {
			//
			// Apply to everyone
			//
			for i := uint16(0); i <= g.maxPartyMemberIndex; i++ {
				w := g.parties[i].PlayerRole
				increaseHPMP(w, 0, (mkf.SHORT)(pScript.Operand[1]))
			}
		} else {
			//
			// Apply to one player. The wEventObjectID parameter should indicate the player role.
			//
			increaseHPMP(eventObjID, 0, (mkf.SHORT)(pScript.Operand[1]))
		}
	case 0x001D:
		//
		// Increase/decrease player's HP and MP
		//
		if pScript.Operand[0] != 0 {
			//
			// Apply to everyone
			//
			for i := uint16(0); i <= g.maxPartyMemberIndex; i++ {
				w := g.parties[i].PlayerRole
				increaseHPMP(w, (mkf.SHORT)(pScript.Operand[1]), (mkf.SHORT)(pScript.Operand[1]))
			}
		} else {
			//
			// Apply to one player. The wEventObjectID parameter should indicate the player role.
			//
			increaseHPMP(eventObjID, (mkf.SHORT)(pScript.Operand[1]), (mkf.SHORT)(pScript.Operand[1]))
		}
	case 0x001E:
		//
		// Increase or decrease cash by the specified amount
		//
		if (mkf.SHORT)(pScript.Operand[0]) < 0 &&
			g.cash < mkf.DWORD(-(mkf.SHORT)(pScript.Operand[0])) {
			//
			// not enough cash
			//
			scriptEntry = pScript.Operand[1] - 1
		} else {
			g.cash += mkf.DWORD(mkf.SHORT(pScript.Operand[0]))
		}
	case 0x001F:
		//
		// Add item to inventory
		//
		addItemIndexToInventory(pScript.Operand[0], int(mkf.SHORT(pScript.Operand[1])))
	case 0x0020:
		//
		// Remove item from inventory
		//
		x := pScript.Operand[1]
		if x == 0 {
			x = 1
		}
		// Note: Need to implement PAL_CountItem
		// if x <= PAL_CountItem(pScript.Operand[0]) || pScript.Operand[2] == 0 {
		if pScript.Operand[2] == 0 {
			if !addItemIndexToInventory(pScript.Operand[0], -int(x)) {
				//
				// Try removing equipped item
				//
				// Note: Need to implement equipment removal logic
			}
		} else {
			scriptEntry = pScript.Operand[2] - 1
		}
	case 0x0021:
		//
		// Inflict damage to the enemy
		//
		// Note: Battle system not implemented yet
	case 0x0022:
		//
		// Revive player
		//
		if pScript.Operand[0] != 0 {
			//
			// Apply to everyone
			//
			g.scriptSuccess = false
			for i := uint16(0); i <= g.maxPartyMemberIndex; i++ {
				w := g.parties[i].PlayerRole
				if g.playerRoles.HP[w] == 0 {
					g.playerRoles.HP[w] = g.playerRoles.MaxHP[w] * pScript.Operand[1] / 10
					// Note: Need to implement PAL_CurePoisonByLevel and PAL_RemovePlayerStatus
					g.scriptSuccess = true
				}
			}
		} else {
			//
			// Apply to one player
			//
			if g.playerRoles.HP[eventObjID] == 0 {
				g.playerRoles.HP[eventObjID] = g.playerRoles.MaxHP[eventObjID] * pScript.Operand[1] / 10
				// Note: Need to implement PAL_CurePoisonByLevel and PAL_RemovePlayerStatus
			} else {
				g.scriptSuccess = false
			}
		}
	case 0x0023:
		//
		// Remove equipment from the specified player
		//
		if pScript.Operand[1] == 0 {
			//
			// Remove all equipments
			//
			for i := 0; i < mkf.MAX_PLAYER_EQUIPMENTS; i++ {
				w := g.playerRoles.Equipment[i][iPlayerRole]
				if w != 0 {
					addItemIndexToInventory(w, 1)
					g.playerRoles.Equipment[i][iPlayerRole] = 0
				}
				// Note: Need to implement PAL_RemoveEquipmentEffect
			}
		} else {
			w := g.playerRoles.Equipment[pScript.Operand[1]-1][iPlayerRole]
			if w != 0 {
				// Note: Need to implement PAL_RemoveEquipmentEffect
				addItemIndexToInventory(w, 1)
				g.playerRoles.Equipment[pScript.Operand[1]-1][iPlayerRole] = 0
			}
		}
	case 0x0024:
		//
		// Set the autoscript entry address for an event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.AutoScript = pScript.Operand[1]
		}
	case 0x0025:
		//
		// Set the trigger script entry address for an event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.TriggerScript = pScript.Operand[1]
		}
	case 0x0026:
		//
		// Show the buy item menu
		//
		// Note: Need to implement PAL_BuyMenu
	case 0x0027:
		//
		// Show the sell item menu
		//
		// Note: Need to implement PAL_SellMenu
	case 0x0028:
		//
		// Apply poison to enemy
		//
		// Note: Battle system not implemented yet
	case 0x0029:
		//
		// Apply poison to player
		//
		// Note: Poison system not implemented yet
	case 0x002A:
		//
		// Cure poison by object ID for enemy
		//
		// Note: Poison system not implemented yet
	case 0x002B:
		//
		// Cure poison by object ID for player
		//
		// Note: Poison system not implemented yet
	case 0x002C:
		//
		// Cure poisons by level
		//
		// Note: Poison system not implemented yet
	case 0x002D:
		//
		// Set the status for player
		//
		// Note: Status system not implemented yet
	case 0x002E:
		//
		// Set the status for enemy
		//
		// Note: Status system not implemented yet
	case 0x002F:
		//
		// Remove player's status
		//
		// Note: Status system not implemented yet
	case 0x0030:
		//
		// Increase player's stat temporarily by percent
		//
		// Note: Temporary stat modification not implemented yet
	case 0x0031:
		//
		// Change battle sprite temporarily for player
		//
		// Note: Battle sprite system not implemented yet
	case 0x0033:
		//
		// collect the enemy for items
		//
		// Note: Enemy collection system not implemented yet
	case 0x0034:
		//
		// Transform collected enemies into items
		//
		// Note: Enemy collection system not implemented yet
	case 0x0035:
		//
		// Shake the screen
		//
		// Note: Need to implement VIDEO_ShakeScreen
	case 0x0036:
		//
		// Set the current playing RNG animation
		//
		// Note: RNG animation system not implemented yet
	case 0x0037:
		//
		// Play RNG animation
		//
		// Note: RNG animation system not implemented yet
	case 0x0038:
		//
		// Teleport the party out of the scene
		//
		if !g.inBattle &&
			g.scenes[g.crtSceneNum-1].ScriptOnTeleport != 0 {
			runTriggerScript(g.scenes[g.crtSceneNum-1].ScriptOnTeleport, 0xFFFF)
		} else {
			//
			// failed
			//
			g.scriptSuccess = false
			scriptEntry = pScript.Operand[0] - 1
		}
	case 0x0039:
		//
		// Drain HP from enemy
		//
		// Note: Battle system not implemented yet
	case 0x003A:
		//
		// Player flee from the battle
		//
		// Note: Battle system not implemented yet
	case 0x003F:
		//
		// Ride the event object to the specified position, at a low speed
		//
		// Note: Need to implement PAL_PartyRideEventObject
	case 0x0040:
		//
		// set the trigger method for a event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.TriggerMode = pScript.Operand[1]
		}
	case 0x0041:
		//
		// Mark the script as failed
		//
		g.scriptSuccess = false
	case 0x0042:
		//
		// Simulate a magic for player
		//
		// Note: Battle magic system not implemented yet
	case 0x0043:
		//
		// Set background music
		//
		// Note: Audio system not implemented yet
	case 0x0044:
		//
		// Ride the event object to the specified position, at the normal speed
		//
		// Note: Need to implement PAL_PartyRideEventObject
	case 0x0045:
		//
		// Set battle music
		//
		// Note: Audio system not implemented yet
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
	case 0x0047:
		//
		// Play sound effect
		//
		// Note: Audio system not implemented yet
	case 0x0049:
		//
		// Set the state of event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.State = mkf.SHORT(pScript.Operand[1])
		}
	case 0x004A:
		//
		// Set the current battlefield
		//
		// Note: Battle system not implemented yet
	case 0x004B:
		//
		// Nullify the event object for a short while
		//
		pEvtObj.VanishTime = -15
	case 0x004C:
		//
		// chase the player
		//
		// Note: Need to implement PAL_MonsterChasePlayer
	case 0x004D:
		//
		// wait for any key
		//
		// Note: Need to implement PAL_WaitForKey
	case 0x004E:
		//
		// Load the last saved game
		//
		// Note: Save/load system not implemented yet
		return 0 // don't go further
	case 0x004F:
		//
		// Fade the screen to red color (game over)
		//
		// Note: Need to implement PAL_FadeToRed
	case 0x0050:
		//
		// screen fade out
		//
		// Note: Need to implement PAL_FadeOut
	case 0x0051:
		//
		// screen fade in
		//
		// Note: Need to implement PAL_FadeIn
	case 0x0052:
		//
		// hide the event object for a while, default 800 frames
		//
		pEvtObj.State *= mkf.SHORT(-1)
		if pScript.Operand[0] != 0 {
			pEvtObj.VanishTime = mkf.SHORT(pScript.Operand[0])
		} else {
			pEvtObj.VanishTime = 800
		}
	case 0x0053:
		//
		// use the day palette
		//
		g.night = false
	case 0x0054:
		//
		// use the night palette
		//
		g.night = true
	case 0x0055:
		//
		// Add magic to a player
		//
		// Note: Magic system not implemented yet
	case 0x0056:
		//
		// Remove magic from a player
		//
		// Note: Magic system not implemented yet
	case 0x0057:
		//
		// Set the base damage of magic according to MP value
		//
		// Note: Magic system not implemented yet
	case 0x0058:
		//
		// Jump if there is less than the specified number of the specified items
		// in the inventory
		//
		// Note: Need to implement PAL_GetItemAmount
	case 0x0059:
		//
		// Change to the specified scene
		//
		if pScript.Operand[0] > 0 && pScript.Operand[0] <= mkf.MAX_SCENES &&
			g.crtSceneNum != pScript.Operand[0] {
			//
			// Set data to load the scene in the next frame
			//
			g.crtSceneNum = pScript.Operand[0]
			// Note: Need to implement PAL_SetLoadFlags and PAL_LoadResources
			//PAL_SetLoadFlags(kLoadScene);
			//gpGlobals->fEnteringScene = TRUE;
			g.viewport = PAL_XY(0, 0)
		}
	case 0x005A:
		//
		// Halve the player's HP
		// The wEventObjectID parameter here should indicate the player role
		//
		g.playerRoles.HP[eventObjID] /= 2
	case 0x005B:
		//
		// Halve the enemy's HP
		//
		// Note: Battle system not implemented yet
	case 0x005C:
		//
		// Hide for a while
		//
		// Note: Battle system not implemented yet
	case 0x005D:
		//
		// Jump if player doesn't have the specified poison
		//
		// Note: Poison system not implemented yet
	case 0x005E:
		//
		// Jump if enemy doesn't have the specified poison
		//
		// Note: Battle and poison system not implemented yet
	case 0x005F:
		//
		// Kill the player immediately
		// The wEventObjectID parameter here should indicate the player role
		//
		g.playerRoles.HP[eventObjID] = 0
	case 0x0060:
		//
		// Immediate KO of the enemy
		//
		// Note: Battle system not implemented yet
	case 0x0061:
		//
		// Jump if player is not poisoned
		//
		// Note: Poison system not implemented yet
	case 0x0062:
		//
		// Pause enemy chasing for a while
		//
		// Note: Enemy chasing system not implemented yet
	case 0x0063:
		//
		// Speed up enemy chasing for a while
		//
		// Note: Enemy chasing system not implemented yet
	case 0x0064:
		//
		// Jump if enemy's HP is more than the specified percentage
		//
		// Note: Battle system not implemented yet
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
		// Note: Need to implement PAL_SetLoadFlags and PAL_LoadResources
	case 0x0076:
		//
		// Show FBP picture
		//
		// Note: Need to implement PAL_ShowFBP
	case 0x0077:
		//
		// Stop music
		//
		// Note: Audio system not implemented yet
	case 0x0078:
		//
		// NOP
		//
	case 0x0079:
		//
		// Jump if the player is in the party
		//
		found := false
		for i := uint16(0); i <= g.maxPartyMemberIndex; i++ {
			if g.parties[i].PlayerRole == pScript.Operand[0]-1 {
				found = true
				break
			}
		}
		if found {
			scriptEntry = pScript.Operand[1] - 1
		}
	case 0x007A:
		//
		// Move the party to the specified position, at a high speed
		//
		// Note: Need to implement PAL_PartyWalkTo
	case 0x007B:
		//
		// Move the party to the specified position, at the highest speed
		//
		// Note: Need to implement PAL_PartyWalkTo
	case 0x007D:
		//
		// Move the event object
		//
		pCurrent.X += uint16(mkf.SHORT(pScript.Operand[1]))
		pCurrent.Y += uint16(mkf.SHORT(pScript.Operand[2]))
	case 0x007E:
		//
		// Set the layer of the event object
		//
		if pScript.Operand[0] != 0 {
			pCurrent.Layer = mkf.SHORT(pScript.Operand[1])
		}
	case 0x007F:
		//
		// Move the viewport
		//
		// Note: Need to implement viewport movement with animation
	case 0x0080:
		//
		// Switch day/night palette
		//
		g.night = !g.night
	case 0x0081:
		//
		// Jump if the party is not facing the event object
		//
		// Note: Need to implement facing check
	case 0x0082:
		//
		// Walk to the specified position, high speed
		//
		if !npcWalkTo(eventObjID, mkf.INT(pScript.Operand[0]), mkf.INT(pScript.Operand[1]), mkf.INT(pScript.Operand[2]), 4) {
			scriptEntry--
		}
	case 0x0083:
		//
		// Jump if the event object is not in the specified zone
		//
		// Note: Need to implement zone check
	case 0x0085:
		//
		// Delay
		//
		// Note: Need to implement delay
	case 0x0086:
		//
		// Jump if the item is not equipped
		//
		// Note: Need to implement equipment check
	case 0x0087:
		//
		// Animate the event object
		//
		// Note: Need to implement animation
	case 0x0088:
		//
		// Set magic damage based on cash
		//
		// Note: Magic system not implemented yet
	case 0x0089:
		//
		// Set battle result
		//
		// Note: Battle system not implemented yet
	case 0x008A:
		//
		// Auto battle
		//
		// Note: Battle system not implemented yet
	case 0x008B:
		//
		// Set palette
		//
		// Note: Need to implement palette system
	case 0x008C:
		//
		// Color fade
		//
		// Note: Need to implement color fade
	case 0x008D:
		//
		// Level up the player
		//
		// Note: Need to implement level up
	case 0x008E:
		//
		// Restore the screen
		//
		// Note: Need to implement VIDEO_RestoreScreen
	case 0x008F:
		//
		// Halve cash
		//
		g.cash /= 2
	case 0x0090:
		//
		// Set object script
		//
		if pScript.Operand[0] != 0 {
			i := pScript.Operand[2]
			if i == 0 {
				pCurrent.AutoScript = pScript.Operand[1]
			} else if i == 1 {
				pCurrent.TriggerScript = pScript.Operand[1]
			}
		}
	case 0x0091:
		//
		// Jump if not first enemy of its type
		//
		// Note: Battle system not implemented yet
	case 0x0092:
		//
		// Show player magic animation
		//
		// Note: Battle animation system not implemented yet
	case 0x0093:
		//
		// Scene fade
		//
		// Note: Need to implement scene fade
	case 0x0094:
		//
		// Jump if event object state equals
		//
		if pScript.Operand[0] != 0 && pCurrent.State == mkf.SHORT(pScript.Operand[1]) {
			scriptEntry = pScript.Operand[2] - 1
		}
	case 0x0095:
		//
		// Jump if scene equals
		//
		if g.crtSceneNum == pScript.Operand[0] {
			scriptEntry = pScript.Operand[1] - 1
		}
	case 0x0096:
		//
		// Ending animation
		//
		// Note: Need to implement ending animation
	case 0x0097:
		//
		// Ride fast
		//
		// Note: Need to implement PAL_PartyRideEventObject
	case 0x0098:
		//
		// Set follower
		//
		// Note: Follower system not implemented yet
	case 0x0099:
		//
		// Change map
		//
		// Note: Need to implement map change
	case 0x009A:
		//
		// Set multiple event object states
		//
		startID := pScript.Operand[0]
		endID := pScript.Operand[1]
		state := mkf.SHORT(pScript.Operand[2])
		if startID != 0 {
			if endID < startID {
				endID = startID
			}
			for i := startID; i <= endID; i++ {
				if i-1 < mkf.WORD(len(g.eventObjects)) {
					g.eventObjects[i-1].State = state
				}
			}
		}
	case 0x009B:
		//
		// Fade to scene fixed
		//
		// Note: Need to implement fade
	case 0x009C:
		//
		// Enemy divide
		//
		// Note: Battle system not implemented yet
	case 0x009E:
		//
		// Enemy summons another monster
		//
		// Note: Battle system not implemented yet
	case 0x009F:
		//
		// Enemy transforms into something else
		//
		// Note: Battle system not implemented yet
	case 0x00A0:
		//
		// Quit game
		//
		// Note: Need to implement quit
	case 0x00A1:
		//
		// Set the positions of all party members to the same as the first one
		//
		for i := 0; i < mkf.MAX_PLAYABLE_PLAYER_ROLES; i++ {
			g.trails[i].Direction = g.partyDirection
			g.trails[i].X = mkf.WORD(int(g.parties[0].X) + int(g.viewport.X()))
			g.trails[i].Y = mkf.WORD(int(g.parties[0].Y) + int(g.viewport.Y()))
		}
		for i := uint16(1); i <= g.maxPartyMemberIndex; i++ {
			g.parties[i].X = g.parties[0].X
			g.parties[i].Y = g.parties[0].Y - 1
		}
		// Note: Need to implement PAL_UpdatePartyGestures
	case 0x00A2:
		//
		// Jump to one of the following instructions randomly
		//
		scriptEntry += mkf.WORD(rand.Intn(int(pScript.Operand[0])))
	case 0x00A3:
		//
		// Play CD music. Use the RIX music for fallback.
		//
		// Note: Audio system not implemented yet
	case 0x00A4:
		//
		// Scroll FBP to the screen
		//
		// Note: Need to implement PAL_ScrollFBP
	case 0x00A5:
		//
		// Show FBP picture with sprite effects
		//
		// Note: Need to implement PAL_ShowFBP
	case 0x00A6:
		//
		// backup screen
		//
		// Note: Need to implement VIDEO_BackupScreen
	default:
		// Note: Invalid instruction, should handle error
	}

	return scriptEntry + 1
}

func npcWalkOneStep(evtObjID mkf.WORD, speed mkf.INT) {
	g := &globals.G

	//
	// Check for invalid parameters
	//
	if evtObjID == 0 || evtObjID > mkf.WORD(len(g.eventObjects)) {
		return
	}

	p := &g.eventObjects[evtObjID-1]

	//
	// Move the event object by the specified direction
	//
	if p.Direction == kDirWest || p.Direction == kDirSouth {
		p.X -= uint16(2 * speed)
	} else {
		p.X += uint16(2 * speed)
	}

	if p.Direction == kDirWest || p.Direction == kDirNorth {
		p.Y -= uint16(speed)
	} else {
		p.Y += uint16(speed)
	}

	//
	// Update the gesture
	//
	if p.SpriteFrames > 0 {
		p.CurrentFrameNum++
		if p.SpriteFrames == 3 {
			p.CurrentFrameNum %= 4
		} else {
			p.CurrentFrameNum %= p.SpriteFrames
		}
	} else if p.SpriteFramesAuto > 0 {
		p.CurrentFrameNum++
		p.CurrentFrameNum %= p.SpriteFramesAuto
	}
}

func npcWalkTo(evtObjID mkf.WORD, x, y, h, speed mkf.INT) bool {
	g := &globals.G
	pEvtObj := &g.eventObjects[evtObjID-1]

	xOffset := (x*32 + h*16) - mkf.INT(pEvtObj.X)
	yOffset := (y*16 + h*8) - mkf.INT(pEvtObj.Y)

	if yOffset < 0 {
		if xOffset < 0 {
			pEvtObj.Direction = kDirWest
		} else {
			pEvtObj.Direction = kDirNorth
		}
	} else {
		if xOffset < 0 {
			pEvtObj.Direction = kDirSouth
		} else {
			pEvtObj.Direction = kDirEast
		}
	}

	if abs(int(xOffset)) < int(speed)*2 || abs(int(yOffset)) < int(speed)*2 {
		pEvtObj.X = uint16(x*32 + h*16)
		pEvtObj.Y = uint16(y*16 + h*8)
	} else {
		npcWalkOneStep(evtObjID, speed)
	}

	if pEvtObj.X == uint16(x*32+h*16) && pEvtObj.Y == uint16(y*16+h*8) {
		pEvtObj.CurrentFrameNum = 0
		return true
	}

	return false
}

func increaseHPMP(playerRole mkf.WORD, HP mkf.SHORT, MP mkf.SHORT) bool {
	g := &globals.G

	if HP != 0 {
		newHP := int(g.playerRoles.HP[playerRole]) + int(HP)
		if newHP < 0 {
			newHP = 0
		}
		if newHP > int(g.playerRoles.MaxHP[playerRole]) {
			newHP = int(g.playerRoles.MaxHP[playerRole])
		}
		g.playerRoles.HP[playerRole] = mkf.WORD(newHP)
	}

	if MP != 0 {
		newMP := int(g.playerRoles.MP[playerRole]) + int(MP)
		if newMP < 0 {
			newMP = 0
		}
		if newMP > int(g.playerRoles.MaxMP[playerRole]) {
			newMP = int(g.playerRoles.MaxMP[playerRole])
		}
		g.playerRoles.MP[playerRole] = mkf.WORD(newMP)
	}

	return g.playerRoles.HP[playerRole] > 0
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
	// Note: Need to implement proper inventory management
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
