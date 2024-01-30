package tank

import (
	"log"
)

type ControlPanel struct {
	dir            Direction
	fireButtonDown bool
}

func (cp *ControlPanel) Drive(dir Direction) {
	cp.dir = dir
}

func (cp *ControlPanel) Hold() {
	cp.Drive(dirNone)
}

func (cp *ControlPanel) PressFire() {
	cp.fireButtonDown = true
}

func (cp *ControlPanel) UnPressFire() {
	cp.fireButtonDown = false
}

type tankEquipCondition int

const (
	equipNone tankEquipCondition = iota
	equipSheld
	equipShip
)

type tankStatus int

const (
	tankBorn tankStatus = iota
	tankNormal
	tankFreeze
	tankBoom
)

type Tank struct {
	GameObject
	count  int
	level  int
	status tankStatus
	equip  tankEquipCondition
	//tankType
	//speed          int // speed in pix
	fireInterval   int
	lastFire       int
	maxBulletQuant int
	flyingBullet   int
	turning        bool

	cp ControlPanel
}

func (t *Tank) Update(count int) {
	//if t.status == tankBorn && t.count > 90 {
	//	t.SetStatus(tankNormal)
	//}
	t.handleControl()
	t.GameObject.Update(count)
	t.count++
}

func (t *Tank) SetStatus(status tankStatus) {
	if status == t.status {
		return
	}
	t.status = status
	switch status {
	case tankBorn:
	case tankNormal:
	case tankBoom:
	case tankFreeze:
	default:
		log.Fatal("")
	}
	t.changeAnimation()
}

func (t *Tank) changeAnimation() {
	t.animation = GetP1TankAnimation(t.level, t.equip, t.status)
}

func (t *Tank) handleControl() {
	if t.status != tankNormal {
		return
	}
	if t.cp.dir != dirNone {
		t.animation.Resume()
		if !t.Moving {
			if t.NextDir != t.cp.dir {
				t.turning = true
				Timer.AddEvent(6, func() {
					t.GameObject.Moving = true
					t.turning = false
				})
				t.NextDir = t.cp.dir
			} else {
				if !t.turning {
					t.GameObject.Moving = true
				}
			}
		}

	} else {
		t.GameObject.Moving = false
		if t.status == tankNormal && t.equip == equipNone {
			t.animation.Pause()
		}
	}
	if t.cp.fireButtonDown {
		t.fire()
	}
}

func (t *Tank) fire() {
	//NewBullet(t.X, t.Y, 1, t.dir)
	if t.count-t.lastFire < t.fireInterval {
		return
	}
	if t.flyingBullet >= t.maxBulletQuant {
		return
	}
	t.lastFire = t.count
	t.flyingBullet++
	t.Ground.NewBullet(t.X, t.Y, 1, t.Dir, t)
}

func (t *Tank) NotifyHit() {
	t.flyingBullet--
}

func NewP1Tank(x, y int) *Tank {
	ret := Tank{
		status: tankBorn,
	}
	ret.ObjType = P1TankType
	ret.Speed = 3
	ret.X = x
	ret.Y = y
	ret.Moving = false
	ret.collisionSize = tankSizeInPix
	ret.Dir = dirUP
	ret.cp.dir = dirNone
	ret.ObjType = P1TankType
	ret.fireInterval = 10
	ret.maxBulletQuant = 1
	ret.host = &ret
	ret.animation = GetP1TankAnimation(0, equipNone, tankBorn)
	Timer.AddEvent(90, func() {
		ret.SetStatus(tankNormal)
	})
	return &ret
}

func GetP1TankAnimation(level int, equip tankEquipCondition, status tankStatus) *Animation {
	if status == tankBorn {
		return NewAnimation(15).appendImages(gr.GetBornImages())
	}
	switch level {
	case 0:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(0)).appendImage(gr.GetTankImageWithoutBorder(1))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(0)).appendImage(gr.GetTankImage(2))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(0)).appendImage(gr.GetTankImage(1))
		}
	case 1:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(3)).appendImage(gr.GetTankImageWithoutBorder(4))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(3)).appendImage(gr.GetTankImage(5))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(3)).appendImage(gr.GetTankImage(4))
		}
	case 2:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(6)).appendImage(gr.GetTankImageWithoutBorder(7))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(6)).appendImage(gr.GetTankImage(8))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(6)).appendImage(gr.GetTankImage(7))
		}
	}
	log.Fatal("")
	return nil
}

func GetP2TankAnimation(level int, equip tankEquipCondition, status tankStatus) *Animation {
	if status == tankBorn {
		return NewAnimation(15).appendImages(gr.GetBornImages())
	}
	switch level {
	case 0:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(10)).appendImage(gr.GetTankImageWithoutBorder(11))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(10)).appendImage(gr.GetTankImage(12))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(10)).appendImage(gr.GetTankImage(11))
		}
	case 1:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(13)).appendImage(gr.GetTankImageWithoutBorder(14))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(13)).appendImage(gr.GetTankImage(15))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(13)).appendImage(gr.GetTankImage(14))
		}
	case 2:
		switch equip {
		case equipNone:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(16)).appendImage(gr.GetTankImageWithoutBorder(17))
		case equipSheld:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(16)).appendImage(gr.GetTankImage(18))
		case equipShip:
			return NewAnimation(defaultInterval).appendImage(gr.GetTankImage(16)).appendImage(gr.GetTankImage(17))
		}
	}
	log.Fatal("")
	return nil
}
