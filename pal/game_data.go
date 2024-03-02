package main

import (
	"github.com/njuwelkin/games/pal/mkf"
)

const ()

type gameData struct {
	eventObjects      []mkf.EventObject
	scenes            [mkf.MAX_SCENES]mkf.Scene
	objects           [mkf.MAX_OBJECTS]mkf.Object
	scriptEntries     []mkf.ScriptEntry
	stores            []mkf.Store
	enemies           []mkf.Enemy
	enemyTeams        []mkf.EnemyTeam
	playerRoles       mkf.PlayerRoles
	magics            []mkf.Magic
	battleFields      []mkf.BattleField
	leveUpMagics      []mkf.LevelUpMagicAll
	enemyPos          mkf.EnemyPos
	levelUpExp        [mkf.MAX_LEVELS + 1]mkf.WORD
	battleEffectIndex [10][2]mkf.WORD
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
	ret.scriptEntries = []mkf.ScriptEntry{}
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
	ret.stores = []mkf.Store{}
	for i := 0; i < stc.Len(); i++ {
		store := stc.GetStore(i)
		ret.stores = append(ret.stores, store)
	}
	// load enemies
	ec, err := data.GetEnemyChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < ec.Len(); i++ {
		enemy := ec.GetEnemy(i)
		ret.enemies = append(ret.enemies, enemy)
	}
	// load enemy team
	etc, err := data.GetEnemyTeamChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < etc.Len(); i++ {
		enemyTeam := etc.GetEnemyTeam(i)
		ret.enemyTeams = append(ret.enemyTeams, enemyTeam)
	}
	// load magic
	mc, err := data.GetMagicChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < mc.Len(); i++ {
		magic := mc.GetMagic(i)
		ret.magics = append(ret.magics, magic)
	}
	// load battle field
	bfc, err := data.GetBattleFieldChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < bfc.Len(); i++ {
		battleField := bfc.GetBattleField(i)
		ret.battleFields = append(ret.battleFields, battleField)
	}
	// load level up magic
	lumc, err := data.GetLevelUpMagicChunk()
	if err != nil {
		panic(err.Error())
	}
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
}
