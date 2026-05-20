package main

type global struct {
	Text   TextLib
	Config Config
	Font   Font
	G      gameData

	UpdatedInBattle bool
	ScriptSuccess   bool
	NSpriteToDraw   int

	MaxPartyMemberIndex int
	NFollower           int
}

var globals global

type Resource struct {
	LoadFlags byte
	Map       Map
	//LPSPRITE  *lppEventObjectSprites // event object sprites
	//int       nEventObject           // number of event objects

	//LPSPRITE rglpPlayerSprite[MAX_PLAYABLE_PLAYER_ROLES] // player sprites

}

func initGlobalSetting() {
	globals.Text = loadText()
	globals.Config = loadConfig()
	globals.Font = newFont()
	globals.G = loadGameData()
	globals.UpdatedInBattle = false
}
