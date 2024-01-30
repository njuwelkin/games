package gobang

type MenualKey string
type MenualItem struct {
	Level      int
	Estimation Score
}

type ChessMenual struct {
	lru      *LRUCache
	db       *GobangDB
	hitCache int
	queries  int
}

func NewChessMenual() *ChessMenual {
	ret := &ChessMenual{
		lru: NewLRUCache(10000),
	}
	return ret
}

func (cm *ChessMenual) OpenDB() {
	db, err := OpenDB("gobang.db")
	if err != nil {
		cm.db = nil
	}
	cm.db = db
	cm.db.Init()
}

func (cm *ChessMenual) Get(key MenualKey) *MenualItem {
	cm.queries++
	item := cm.lru.Get(key)
	if item == nil && cm.db != nil {
		item, _ = cm.db.Get(string(key))
	}
	if item != nil {
		cm.hitCache++
	}
	return item
}

func (cm *ChessMenual) Put(key MenualKey, item *MenualItem) {
	cm.lru.Put(key, item)
	if cm.db != nil {
		cm.db.Put(string(key), *item)
	}
}

func (cm *ChessMenual) CloseDB() {
	if cm.db != nil {
		cm.db.Close()
	}
}
