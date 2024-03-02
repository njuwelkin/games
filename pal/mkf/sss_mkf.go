package mkf

import "unsafe"

const (
	MAX_SCENES  = 300
	MAX_OBJECTS = 600
)

type SSSMkf struct {
	Mkf
}

func NewSSSMkf(file string) (SSSMkf, error) {
	ret := SSSMkf{Mkf{}}
	return ret, ret.Open(file)
}

func (sm *SSSMkf) GetScriptEntryChunk() (*ScriptEntryChunk, error) {
	buf, err := sm.ReadChunk(4)
	if err != nil {
		return nil, err
	}
	return NewScriptEntryChunk(buf), nil
}

func (sm *SSSMkf) GetMsgOffsetChunk() (*MsgOffsetChunk, error) {
	buf, err := sm.ReadChunk(3)
	if err != nil {
		return nil, err
	}
	var dw DWORD
	return &MsgOffsetChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(dw)))}, nil
}

func (sm *SSSMkf) GetObjects() ([MAX_OBJECTS]Object, error) {
	buf, err := sm.ReadChunk(2)
	if err != nil {
		return [MAX_OBJECTS]Object{}, err
	}
	pc := NewScriptEntryChunk(buf)
	p := pc.Get(0)
	dosObjs := *(*[MAX_OBJECTS]Object)(p)
	var objs [MAX_OBJECTS]Object
	for i := range dosObjs {
		objs[i] = *(*Object)(unsafe.Pointer(&dosObjs[i]))
		objs[i].Data[6] = dosObjs[i].Data[5] // wFlags
		objs[i].Data[5] = 0                  // wScriptDesc or wReserved2
	}
	return objs, nil
}

func (sm *SSSMkf) GetSceneChunk() (*SceneChunk, error) {
	buf, err := sm.ReadChunk(1)
	if err != nil {
		return nil, err
	}
	return &SceneChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(Scene{})))}, nil
}

func (sm *SSSMkf) GetEventObjectChunk() (*EventObjectChunk, error) {
	buf, err := sm.ReadChunk(0)
	if err != nil {
		return nil, err
	}
	return &EventObjectChunk{NewPlaneChunk(buf, int(unsafe.Sizeof(EventObject{})))}, nil
}

type ScriptEntry struct {
	Operation WORD
	Operand   [3]WORD
}

type MsgOffsetChunk struct {
	PlaneChunk
}

// chunk 4 in sss.mkf
type ScriptEntryChunk struct {
	PlaneChunk
}

func NewScriptEntryChunk(data []byte) *ScriptEntryChunk {
	return &ScriptEntryChunk{NewPlaneChunk(data, int(unsafe.Sizeof(ScriptEntry{})))}
}

func (sc *ScriptEntryChunk) GetScriptEntry(idx int) *ScriptEntry {
	return (*ScriptEntry)(sc.Get(idx))
}

type Scene struct {
	MapNum           WORD // number of the map
	ScriptOnEnter    WORD // when entering this scene, execute script from here
	ScriptOnTeleport WORD // when teleporting out of this scene, execute script from here
	EventObjectIndex WORD // event objects in this scene begins from number wEventObjectIndex + 1
}

type SceneChunk struct {
	PlaneChunk
}

func (sc *SceneChunk) GetScene(idx int) Scene {
	return *(*Scene)(sc.Get(idx))
}

type EventObject struct {
	VanishTime               SHORT  // vanish time (?)
	X                        WORD   // X coordinate on the map
	Y                        WORD   // Y coordinate on the map
	Layer                    SHORT  // layer value
	TriggerScript            WORD   // Trigger script entry
	AutoScript               WORD   // Auto script entry
	State                    SHORT  // state of this object
	TriggerMode              WORD   // trigger mode
	SpriteNum                WORD   // number of the sprite
	SpriteFrames             USHORT // total number of frames of the sprite
	Direction                WORD   // direction
	CurrentFrameNum          WORD   // current frame number
	ScriptIdleFrame          USHORT // count of idle frames, used by trigger script
	SpritePtrOffset          WORD   // FIXME: ???
	SpriteFramesAuto         USHORT // total number of frames of the sprite, used by auto script
	ScriptIdleFrameCountAuto WORD   // count of idle frames, used by auto script
}

// chunk 0 in sss.mkf
type EventObjectChunk struct {
	PlaneChunk
}

func (ec *EventObjectChunk) GetEventObject(idx int) EventObject {
	return *(*EventObject)(ec.Get(idx))
}

type ObjectPlayer struct {
	wReserved            [2]WORD // always zero
	wScriptOnFriendDeath WORD    // when friends in party dies, execute script from here
	wScriptOnDying       WORD    // when dying, execute script from here
}

type ObjectItemDos struct {
	Bitmap        WORD // bitmap number in BALL.MKF
	Price         WORD // price
	ScriptOnUse   WORD // script executed when using this item
	ScriptOnEquip WORD // script executed when equipping this item
	ScriptOnThrow WORD // script executed when throwing this item to enemy
	Flags         WORD // flags
}

type ObjectItem struct {
	Bitmap        WORD // bitmap number in BALL.MKF
	Price         WORD // price
	ScriptOnEquip WORD // script executed when equipping this item
	ScriptOnThrow WORD // script executed when throwing this item to enemy
	ScriptDesc    WORD // description script
	Flags         WORD // flags
}

type ObjectMagicDos struct {
	MagicNumber     WORD // magic number, according to DATA.MKF #3
	Reserved1       WORD // always zero
	ScriptOnSuccess WORD // when magic succeed, execute script from here
	ScriptOnUse     WORD // when use this magic, execute script from here
	Reserved2       WORD // always zero
	Flags           WORD // flags
}

type ObjectMagic struct {
	MagicNumber     WORD // magic number, according to DATA.MKF #3
	Reserved1       WORD // always zero
	ScriptOnSuccess WORD // when magic succeed, execute script from here
	ScriptOnUse     WORD // when use this magic, execute script from here
	ScriptDesc      WORD // description script
	Reserved2       WORD // always zero
	Flags           WORD // flags
}

type ObjectEnemy struct {
	EnemyID WORD // ID of the enemy, according to DATA.MKF #1.
	// Also indicates the bitmap number in ABC.MKF.
	ResistanceToSorcery WORD // resistance to sorcery and poison (0 min, 10 max)
	ScriptOnTurnStart   WORD // script executed when turn starts
	ScriptOnBattleEnd   WORD // script executed when battle ends
	ScriptOnReady       WORD // script executed when the enemy is ready
}

type ObjectPosition struct {
	PoisonLevel  WORD // level of the poison
	Color        WORD // color of avatars
	PlayerScript WORD // script executed when player has this poison (per round)
	Reserved     WORD // always zero
	EnemyScript  WORD // script executed when enemy has this poison (per round)
}

type ObjectDos struct {
	Data   [6]WORD
	Player ObjectPlayer
	Item   ObjectItemDos
	Magic  ObjectMagicDos
	Enemy  ObjectEnemy
	Poison ObjectPosition
}

type Object struct {
	Data   [7]WORD
	Player ObjectPlayer
	Item   ObjectItemDos
	Magic  ObjectMagicDos
	Enemy  ObjectEnemy
	Poison ObjectPosition
}
