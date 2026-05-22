package main

import (
	"path/filepath"

	"github.com/njuwelkin/games/pal/mkf"
)

const (
	MAX_INVENTORY = 256
)

type GameData struct {
	// might change, load from mkf or overwrite by save file
	EventObjects []mkf.EventObject
	Scenes       [mkf.MAX_SCENES]mkf.Scene
	Objects      [mkf.MAX_OBJECTS]mkf.Object
	PlayerRoles  mkf.PlayerRoles

	// global data, load from mkf, never change
	ScriptEntries     []mkf.ScriptEntry
	Stores            []mkf.Store
	Enemies           []mkf.Enemy
	EnemyTeams        []mkf.EnemyTeam
	Magics            []mkf.Magic
	BattleFields      []mkf.BattleField
	LevelUpMagics     []mkf.LevelUpMagicAll
	EnemyPos          mkf.EnemyPos
	LevelUpExp        [mkf.MAX_LEVELS + 1]mkf.WORD
	BattleEffectIndex [10][2]mkf.WORD
	Avatars           []*mkf.BitMap

	// dynamic data
	CrtSceneNum   mkf.WORD // sdlpal-1
	PartyOffset   PAL_POS
	Cash          mkf.DWORD
	CrtPaletteNum mkf.WORD
	Night         bool
	//axPartyMemberIndex mkf.WORD

	// others
	Viewport            PAL_POS
	Parties             [mkf.MAX_PLAYABLE_PLAYER_ROLES]Party
	Trails              [mkf.MAX_PLAYABLE_PLAYER_ROLES]Trail
	FrameNum            mkf.DWORD
	PartyDirection      mkf.WORD
	EquipmentEffect     [mkf.MAX_PLAYER_EQUIPMENTS + 1]mkf.PlayerRoles
	CurEquipPart        int
	ScriptSuccess       bool
	MaxPartyMemberIndex mkf.WORD
	InBattle            bool
	Inventory           [MAX_INVENTORY]Inventory
	LastUnequippedItem  mkf.WORD
}

func LoadGameData() GameData {
	ret := GameData{}

	// load from sss.mkf
	sss, err := mkf.NewSSSMkf(filepath.Join(Globals.Config.GamePath, "SSS.MKF"))
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
	ret.ScriptEntries = make([]mkf.ScriptEntry, 0, sc.Len()) //[]mkf.ScriptEntry{}
	for i := 0; i < sc.Len(); i++ {
		entry := sc.GetScriptEntry(i)
		ret.ScriptEntries = append(ret.ScriptEntries, *entry)
	}

	// load from data.mkf
	data, err := mkf.NewDataMkf(filepath.Join(Globals.Config.GamePath, "DATA.MKF"))
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
	ret.Stores = make([]mkf.Store, 0, stc.Len()) //[]mkf.Store{}
	for i := 0; i < stc.Len(); i++ {
		store := stc.GetStore(i)
		ret.Stores = append(ret.Stores, store)
	}
	// load enemies
	ec, err := data.GetEnemyChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.Enemies = make([]mkf.Enemy, 0, ec.Len())
	for i := 0; i < ec.Len(); i++ {
		enemy := ec.GetEnemy(i)
		ret.Enemies = append(ret.Enemies, enemy)
	}
	// load enemy team
	etc, err := data.GetEnemyTeamChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.EnemyTeams = make([]mkf.EnemyTeam, 0, etc.Len())
	for i := 0; i < etc.Len(); i++ {
		enemyTeam := etc.GetEnemyTeam(i)
		ret.EnemyTeams = append(ret.EnemyTeams, enemyTeam)
	}
	// load magic
	mc, err := data.GetMagicChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.Magics = make([]mkf.Magic, 0, mc.Len())
	for i := 0; i < mc.Len(); i++ {
		magic := mc.GetMagic(i)
		ret.Magics = append(ret.Magics, magic)
	}
	// load battle field
	bfc, err := data.GetBattleFieldChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.BattleFields = make([]mkf.BattleField, 0, bfc.Len())
	for i := 0; i < bfc.Len(); i++ {
		battleField := bfc.GetBattleField(i)
		ret.BattleFields = append(ret.BattleFields, battleField)
	}
	// load level up magic
	lumc, err := data.GetLevelUpMagicChunk()
	if err != nil {
		panic(err.Error())
	}
	ret.LevelUpMagics = make([]mkf.LevelUpMagicAll, 0, lumc.Len())
	for i := 0; i < lumc.Len(); i++ {
		lum := lumc.GetLevelUpMagic(i)
		ret.LevelUpMagics = append(ret.LevelUpMagics, lum)
	}
	// load battle effect idx
	beIdx, err := data.GetBattleEffectIndex()
	if err != nil {
		panic(err.Error())
	}
	ret.BattleEffectIndex = beIdx
	// load enemy pos
	enemyPos, err := data.GetEnemyPos()
	if err != nil {
		panic(err.Error())
	}
	ret.EnemyPos = enemyPos
	// load levelup exp
	levelUpExp, err := data.GetLevelUpExp()
	if err != nil {
		panic(err.Error())
	}
	ret.LevelUpExp = levelUpExp

	// load from rgm.mkf (character face bitmaps)
	ret.LoadAvatar()

	ret.LoadDefault()
	return ret
}

func (gd *GameData) LoadAvatar() {
	rgm, err := mkf.NewRgmMkf(filepath.Join(Globals.Config.GamePath, "RGM.MKF"))
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		rgm.Close()
	}()

	size, err := rgm.GetChunkCount()
	if err != nil {
		panic(err.Error())
	}
	// 加载所有脸部位图（预设最大200个）
	gd.Avatars = make([]*mkf.BitMap, size)
	for i := 0; i < int(size); i++ {
		bmp, err := rgm.GetFaceBmp(mkf.INT(i))
		if err != nil || bmp == nil {
			gd.Avatars[i] = nil
			continue
		}
		gd.Avatars[i] = bmp
	}
}

func (gd *GameData) LoadDefault() {
	// load from sss.mkf
	sss, err := mkf.NewSSSMkf(filepath.Join(Globals.Config.GamePath, "SSS.MKF"))
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
		gd.EventObjects = append(gd.EventObjects, eo)
	}
	// load scenes
	sc, err := sss.GetSceneChunk()
	if err != nil {
		panic(err.Error())
	}
	for i := 0; i < sc.Len(); i++ {
		scene := sc.GetScene(i)
		gd.Scenes[i] = scene
	}
	// load objects
	objs, err := sss.GetObjects()
	if err != nil {
		panic(err.Error())
	}
	gd.Objects = objs

	// load from data.mkf
	data, err := mkf.NewDataMkf(filepath.Join(Globals.Config.GamePath, "DATA.MKF"))
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
	gd.PlayerRoles = *pr

	// others
	gd.CrtSceneNum = 0
	gd.Cash = 0
	gd.CrtPaletteNum = 0
	gd.Night = false
	gd.Viewport = PAL_XY(0, 0)
}

func (gd *GameData) LoadSaved(idx int) {

}

func (gd *GameData) UpdateEquipment() {

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
