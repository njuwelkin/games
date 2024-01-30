package tank

const (
	BulletBorderInPix   = 15
	bulletCollisionSize = 10
)

type Bullet struct {
	GameObject
	level int
	//owner int
}

func NewBullet(x, y, level int, dir Direction, owner *Tank) *Bullet {
	ret := Bullet{
		level: level,
	}
	ret.X, ret.Y = x, y
	ret.Dir = dir
	ret.NextDir = dir
	ret.animation = GetBulletAnimation()
	ret.Speed = ret.getSpeed()
	ret.collisionSize = bulletCollisionSize
	ret.ObjType = BulletType
	ret.Moving = true
	ret.host = &ret
	return &ret
}

func (b *Bullet) Update(count int) {
	b.GameObject.Update(count)
}

func (b *Bullet) getSpeed() int {
	return (b.level + 1) * 2
}

func GetBulletAnimation() *Animation {
	return NewAnimation(defaultInterval).appendImage(gr.GetBulletImage())
}
