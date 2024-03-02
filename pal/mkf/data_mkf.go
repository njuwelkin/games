package mkf

import "unsafe"

const (
	MAX_STORE_ITEM            = 9
	NUM_MAGIC_ELEMENTAL       = 5
	MAX_ENEMIES_IN_TEAM       = 5
	MAX_PLAYABLE_PLAYER_ROLES = 5
	MAX_LEVELS                = 99
	MAX_PLAYER_ROLES          = 6
	MAX_PLAYER_EQUIPMENTS     = 6
	MAX_PLAYER_MAGICS         = 32
)

type DataMkf struct {
	Mkf
}

func NewDataMkf(path string) (DataMkf, error) {
	ret := DataMkf{Mkf{}}
	return ret, ret.Open(path)
}

func (dm *DataMkf) GetStoreChunk() (*StoreChunk, error) {
	buf, err := dm.ReadChunk(0)
	if err != nil {
		return nil, err
	}
	var w WORD
	return &StoreChunk{NewPlaneChunk(buf, int(MAX_STORE_ITEM*unsafe.Sizeof(w)))}, nil
}

func (dm *DataMkf) GetEnemyChunk() (*EnemyChunk, error) {
	buf, err := dm.ReadChunk(1)
	if err != nil {
		return nil, err
	}
	return &EnemyChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(Enemy{})))}, nil
}

func (dm *DataMkf) GetEnemyTeamChunk() (*EnemyTeamChunk, error) {
	buf, err := dm.ReadChunk(2)
	if err != nil {
		return nil, err
	}
	return &EnemyTeamChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(EnemyTeam{})))}, nil
}

func (dm *DataMkf) GetPlayerRoles() (*PlayerRoles, error) {
	buf, err := dm.ReadChunk(3)
	if err != nil {
		return nil, err
	}
	pc := NewPlaneChunk(buf, int(unsafe.Sizeof(PlayerRoles{})))
	p := pc.Get(0)
	return (*PlayerRoles)(p), nil
}

func (dm *DataMkf) GetMagicChunk() (*MagicChunk, error) {
	buf, err := dm.ReadChunk(4)
	if err != nil {
		return nil, err
	}
	return &MagicChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(Magic{})))}, nil
}

func (dm *DataMkf) GetBattleFieldChunk() (*BattleFieldChunk, error) {
	buf, err := dm.ReadChunk(5)
	if err != nil {
		return nil, err
	}
	return &BattleFieldChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(BattleField{})))}, nil
}

func (dm *DataMkf) GetLevelUpMagicChunk() (*LevelUpMagicChunk, error) {
	buf, err := dm.ReadChunk(6)
	if err != nil {
		return nil, err
	}
	return &LevelUpMagicChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(LevelUpMagicAll{})))}, nil
}

func (dm *DataMkf) GetBattleEffectIndex() ([10][2]WORD, error) {
	buf, err := dm.ReadChunk(11)
	if err != nil {
		return [10][2]WORD{}, err
	}
	var w WORD
	pc := NewPlaneChunk(buf, int(unsafe.Sizeof(w*10*2)))
	p := pc.Get(0)
	return *(*[10][2]WORD)(p), nil
}

func (dm *DataMkf) GetEnemyPos() (EnemyPos, error) {
	buf, err := dm.ReadChunk(13)
	if err != nil {
		return EnemyPos{}, err
	}
	pc := NewPlaneChunk(buf, int(unsafe.Sizeof(EnemyPos{})))
	p := pc.Get(0)
	return *(*EnemyPos)(p), nil
}

func (dm *DataMkf) GetLevelUpExp() ([MAX_LEVELS + 1]WORD, error) {
	buf, err := dm.ReadChunk(14)
	if err != nil {
		return [MAX_LEVELS + 1]WORD{}, err
	}
	var w WORD
	pc := NewPlaneChunk(buf, int((MAX_LEVELS+1)*unsafe.Sizeof(w)))
	p := pc.Get(0)
	return *(*[MAX_LEVELS + 1]WORD)(p), nil
}

type Store struct {
	items [MAX_STORE_ITEM]WORD
}

type StoreChunk struct{ PlaneChunk }

func (sc *StoreChunk) GetStore(idx int) Store {
	ret := Store{}
	p := sc.Get(idx)
	for i := 0; i < MAX_STORE_ITEM; i++ {
		var w WORD
		ret.items[i] = *((*WORD)(unsafe.Pointer(uintptr(p) + uintptr(i)*unsafe.Sizeof(w))))
	}
	return ret
}

type Enemy struct {
	IdleFrames          WORD // total number of frames when idle
	MagicFrames         WORD // total number of frames when using magics
	AttackFrames        WORD // total number of frames when doing normal attack
	IdleAnimSpeed       WORD // speed of the animation when idle
	ActWaitFrames       WORD // FIXME: ???
	YPosOffset          WORD
	AttackSound         SHORT                     // sound played when this enemy uses normal attack
	ActionSound         SHORT                     // FIXME: ???
	MagicSound          SHORT                     // sound played when this enemy uses magic
	DeathSound          SHORT                     // sound played when this enemy dies
	CallSound           SHORT                     // sound played when entering the battle
	Health              WORD                      // total HP of the enemy
	Exp                 WORD                      // How many EXPs we'll get for beating this enemy
	Cash                WORD                      // how many cashes we'll get for beating this enemy
	Level               WORD                      // this enemy's level
	Magic               WORD                      // this enemy's magic number
	MagicRate           WORD                      // chance for this enemy to use magic
	AttackEquivItem     WORD                      // equivalence item of this enemy's normal attack
	AttackEquivItemRate WORD                      // chance for equivalence item
	StealItem           WORD                      // which item we'll get when stealing from this enemy
	TotaleStealItem     WORD                      // total amount of the items which can be stolen
	AttackStrength      WORD                      // normal attack strength
	MagicStrength       WORD                      // magical attack strength
	Defense             WORD                      // resistance to all kinds of attacking
	Dexterity           WORD                      // dexterity
	FleeRate            WORD                      // chance for successful fleeing
	PoisonResistance    WORD                      // resistance to poison
	ElemResistance      [NUM_MAGIC_ELEMENTAL]WORD // resistance to elemental magics
	PhysicalResistance  WORD                      // resistance to physical attack
	DualMove            WORD                      // whether this enemy can do dual move or not
	CollectValue        WORD                      // value for collecting this enemy for items
}

type EnemyChunk struct{ PlaneChunk }

func (ec *EnemyChunk) GetEnemy(idx int) Enemy {
	p := ec.Get(idx)
	ret := *((*Enemy)(p))
	return ret
}

type EnemyTeam struct {
	enemy [MAX_ENEMIES_IN_TEAM]WORD
}

type EnemyTeamChunk struct{ PlaneChunk }

func (etc *EnemyTeamChunk) GetEnemyTeam(idx int) EnemyTeam {
	p := etc.Get(idx)
	ret := *((*EnemyTeam)(p))
	return ret
}

type MAGIC_SPECIAL WORD

func (ms MAGIC_SPECIAL) GetSummonEffect() WORD {
	return WORD(ms)
}

func (ms MAGIC_SPECIAL) GetLayerOffset() SHORT {
	return SHORT(ms)
}

type Magic struct {
	Effect      WORD // effect sprite
	Type        WORD // type of this magic
	XOffset     WORD
	YOffset     WORD
	gSpecific   MAGIC_SPECIAL // have multiple meanings
	Speed       SHORT         // speed of the effect
	KeepEffect  WORD          // FIXME: ???
	FireDelay   WORD          // start frame of the magic fire stage
	EffectTimes WORD          // total times of effect
	Shake       WORD          // shake screen
	Wave        WORD          // wave screen
	Unknown     WORD          // FIXME: ???
	CostMP      WORD          // MP cost
	BaseDamage  WORD          // base damage
	Elemental   WORD          // elemental (0 = No Elemental, last = poison)
	Sound       SHORT         // sound played when using this magic
}

type MagicChunk struct{ PlaneChunk }

func (mc *MagicChunk) GetMagic(idx int) Magic {
	p := mc.Get(idx)
	ret := *((*Magic)(p))
	return ret
}

type BattleField struct {
	ScreenWave  WORD                       // level of screen waving
	MagicEffect [NUM_MAGIC_ELEMENTAL]SHORT // effect of attributed magics
}

type BattleFieldChunk struct{ PlaneChunk }

func (mc *BattleFieldChunk) GetBattleField(idx int) BattleField {
	p := mc.Get(idx)
	ret := *((*BattleField)(p))
	return ret
}

type LevelUpMagic struct {
	Level WORD // level reached
	Magic WORD // magic learned
}

type LevelUpMagicAll struct {
	Magics [MAX_PLAYABLE_PLAYER_ROLES]LevelUpMagic
}

type LevelUpMagicChunk struct{ PlaneChunk }

func (mc *LevelUpMagicChunk) GetLevelUpMagic(idx int) LevelUpMagicAll {
	p := mc.Get(idx)
	ret := *((*LevelUpMagicAll)(p))
	return ret
}

type EnemyPos struct {
	Pos [MAX_ENEMIES_IN_TEAM][MAX_ENEMIES_IN_TEAM]PALPOS
}

//type LeveUpExp WORD

type Players [MAX_PLAYER_ROLES]WORD

type PlayerRoles struct {
	rgwAvatar              Players                                       // avatar (shown in status view)
	rgwSpriteNumInBattle   Players                                       // sprite displayed in battle (in F.MKF)
	rgwSpriteNum           Players                                       // sprite displayed in normal scene (in MGO.MKF)
	rgwName                Players                                       // name of player class (in WORD.DAT)
	rgwAttackAll           Players                                       // whether player can attack everyone in a bulk or not
	rgwUnknown1            Players                                       // FIXME: ???
	rgwLevel               Players                                       // level
	rgwMaxHP               Players                                       // maximum HP
	rgwMaxMP               Players                                       // maximum MP
	rgwHP                  Players                                       // current HP
	rgwMP                  Players                                       // current MP
	rgwEquipment           [MAX_PLAYER_EQUIPMENTS][MAX_PLAYER_ROLES]WORD // equipments
	rgwAttackStrength      Players                                       // normal attack strength
	rgwMagicStrength       Players                                       // magical attack strength
	rgwDefense             Players                                       // resistance to all kinds of attacking
	rgwDexterity           Players                                       // dexterity
	rgwFleeRate            Players                                       // chance of successful fleeing
	rgwPoisonResistance    Players                                       // resistance to poison
	rgwElementalResistance [NUM_MAGIC_ELEMENTAL][MAX_PLAYER_ROLES]WORD   // resistance to elemental magics
	rgwUnknown2            Players                                       // FIXME: ???
	rgwUnknown3            Players                                       // FIXME: ???
	rgwUnknown4            Players                                       // FIXME: ???
	rgwCoveredBy           Players                                       // who will cover me when I am low of HP or not sane
	rgwMagic               [MAX_PLAYER_MAGICS][MAX_PLAYER_ROLES]WORD     // magics
	rgwWalkFrames          Players                                       // walk frame (???)
	rgwCooperativeMagic    Players                                       // cooperative magic
	rgwUnknown5            Players                                       // FIXME: ???
	rgwUnknown6            Players                                       // FIXME: ???
	rgwDeathSound          Players                                       // sound played when player dies
	rgwAttackSound         Players                                       // sound played when player attacks
	rgwWeaponSound         Players                                       // weapon sound (???)
	rgwCriticalSound       Players                                       // sound played when player make critical hits
	rgwMagicSound          Players                                       // sound played when player is casting a magic
	rgwCoverSound          Players                                       // sound played when player cover others
	rgwDyingSound          Players                                       // sound played when player is dying
}
