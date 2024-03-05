package main

import (
	"github.com/njuwelkin/games/pal/mkf"
)

const (
	MAX_INVENTORY = 256
)

type gameData struct {
	// might change, load from mkf or overwrite by save file
	eventObjects []mkf.EventObject
	scenes       [mkf.MAX_SCENES]mkf.Scene
	objects      [mkf.MAX_OBJECTS]mkf.Object
	playerRoles  mkf.PlayerRoles

	// global data, load from mkf, never change
	scriptEntries     []mkf.ScriptEntry
	stores            []mkf.Store
	enemies           []mkf.Enemy
	enemyTeams        []mkf.EnemyTeam
	magics            []mkf.Magic
	battleFields      []mkf.BattleField
	leveUpMagics      []mkf.LevelUpMagicAll
	enemyPos          mkf.EnemyPos
	levelUpExp        [mkf.MAX_LEVELS + 1]mkf.WORD
	battleEffectIndex [10][2]mkf.WORD

	// dynamic data
	crtSceneNum   mkf.WORD // sdlpal-1
	partyoffset   PAL_POS
	cash          mkf.DWORD
	crtPaletteNum mkf.WORD
	night         bool
	//axPartyMemberIndex mkf.WORD

	// others
	viewport            PAL_POS
	parties             [mkf.MAX_PLAYABLE_PLAYER_ROLES]Party
	trails              [mkf.MAX_PLAYABLE_PLAYER_ROLES]Trail
	frameNum            mkf.DWORD
	partyDirection      mkf.WORD
	equipmentEffect     [mkf.MAX_PLAYER_EQUIPMENTS + 1]mkf.PlayerRoles
	curEquipPart        int
	scriptSuccess       bool
	maxPartyMemberIndex mkf.WORD
	inBattle            bool
	inventory           [MAX_INVENTORY]Inventory
	lastUnequippedItem  mkf.WORD
}

func loadGameData() gameData {
	ret := gameData{}

	// load from sss.mkf
	sss, err := mkf.NewSSSMkf("SSS.MKF")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		sss.Close()
	}()
	sc, err := sss.GetScriptEntryChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.scriptEntries = make([]mkf.ScriptEntry, 0, sc.Len()) //[]mkf.ScriptEntry{}
	for i := 0; i < sc.Len(); i++ {
		entry := sc.GetScriptEntry(i)
		ret.scriptEntries = append(ret.scriptEntries, *entry)
	}

	// load from data.mkf
	data, err := mkf.NewDataMkf("DATA.MKF")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		data.Close()
	}()
	// load stores
	stc, err := data.GetStoreChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.stores = make([]mkf.Store, 0, stc.Len()) //[]mkf.Store{}
	for i := 0; i < stc.Len(); i++ {
		store := stc.GetStore(i)
		ret.stores = append(ret.stores, store)
	}
	// load enemies
	ec, err := data.GetEnemyChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.enemies = make([]mkf.Enemy, 0, ec.Len())
	for i := 0; i < ec.Len(); i++ {
		enemy := ec.GetEnemy(i)
		ret.enemies = append(ret.enemies, enemy)
	}
	// load enemy team
	etc, err := data.GetEnemyTeamChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.enemyTeams = make([]mkf.EnemyTeam, 0, etc.Len())
	for i := 0; i < etc.Len(); i++ {
		enemyTeam := etc.GetEnemyTeam(i)
		ret.enemyTeams = append(ret.enemyTeams, enemyTeam)
	}
	// load magic
	mc, err := data.GetMagicChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.magics = make([]mkf.Magic, 0, mc.Len())
	for i := 0; i < mc.Len(); i++ {
		magic := mc.GetMagic(i)
		ret.magics = append(ret.magics, magic)
	}
	// load battle field
	bfc, err := data.GetBattleFieldChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.battleFields = make([]mkf.BattleField, 0, bfc.Len())
	for i := 0; i < bfc.Len(); i++ {
		battleField := bfc.GetBattleField(i)
		ret.battleFields = append(ret.battleFields, battleField)
	}
	// load level up magic
	lumc, err := data.GetLevelUpMagicChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.leveUpMagics = make([]mkf.LevelUpMagicAll, 0, lumc.Len())
	for i := 0; i < lumc.Len(); i++ {
		lum := lumc.GetLevelUpMagic(i)
		ret.leveUpMagics = append(ret.leveUpMagics, lum)
	}
	// load battle effect idx
	beIdx, err := data.GetBattleEffectIndex()
	if err != nil {
		panic(err.Error())
	}
	ret.battleEffectIndex = beIdx
	// load enemy pos
	enemyPos, err := data.GetEnemyPos()
	if err != nil {
		panic(err.Error())
	}
	ret.enemyPos = enemyPos
	// load levelup exp
	levelUpExp, err := data.GetLevelUpExp()
	if err != nil {
		panic(err.Error())
	}
	ret.levelUpExp = levelUpExp

	ret.loadDefault()
	return ret
}

func (gd *gameData) loadDefault() {
	// load from sss.mkf
	sss, err := mkf.NewSSSMkf("SSS.MKF")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		sss.Close()
	}()
	// load event objects
	ec, err := sss.GetEventObjectChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < ec.Len(); i++ {
		eo := ec.GetEventObject(i)
		gd.eventObjects = append(gd.eventObjects, eo)
	}
	// load scenes
	sc, err := sss.GetSceneChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < sc.Len(); i++ {
		scene := sc.GetScene(i)
		gd.scenes[i] = scene
	}
	// load objects
	objs, err := sss.GetObjects()
	if err != nil {
		panic(err.Error())
	}
	gd.objects = objs

	// load from data.mkf
	data, err := mkf.NewDataMkf("DATA.MKF")
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		data.Close()
	}()
	// load PlayerRoles
	pr, err := data.GetPlayerRoles()
	if err != nil {
		panic(err.Error())
	}
	gd.playerRoles = *pr

	// others
	gd.crtSceneNum = 0
	gd.cash = 0
	gd.crtPaletteNum = 0
	gd.night = false
	gd.viewport = PAL_XY(0, 0)
}

func (gd *gameData) loadSaved(idx int) {

}

func (gd *gameData) updateEquipment() {

}

type Party struct {
	PlayerRole  mkf.WORD  // player role
	X, Y        mkf.SHORT // position
	Frame       mkf.WORD  // current frame number
	ImageOffset mkf.WORD  // FIXME: ???
}

type Trail struct {
	X, Y      mkf.WORD // position
	Direction mkf.WORD // direction
}

type Inventory struct {
	Item        mkf.WORD   // item object code
	Amount      mkf.USHORT // amount of this item
	AmountInUse mkf.USHORT // in-use amount of this item
}
