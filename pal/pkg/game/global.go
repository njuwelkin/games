package game

type Global struct {
	Text                TextLib
	Config              Config
	Font                Font
	G                   GameData
	UpdatedInBattle     bool
	ScriptSuccess       bool
	NSpriteToDraw       int
	MaxPartyMemberIndex int
	NFollower           int
}

var Globals Global

type Resource struct {
	LoadFlags byte
	Map       Map
	//LPSPRITE  *lppEventObjectSprites // event object sprites
	//int       nEventObject           // number of event objects

	//LPSPRITE rglpPlayerSprite[MAX_PLAYABLE_PLAYER_ROLES] // player sprites

}

func InitGlobalSetting() {
	Globals.Config = loadConfig()
	Globals.Text = loadText()
	Globals.Font = newFont()
	Globals.G = LoadGameData()
	Globals.UpdatedInBattle = false
}
