package marslander

import (
	"errors"
	"math"
)

const (
	Gravity   float32 = 3.711
	MaxRotate         = 90
	MaxPower          = 4
	MaxX              = 7000
	MaxY              = 3000
)

//goland:noinspection GoErrorStringFormat
var (
	OutOfBoundsErr   = errors.New("Mars Lander is lost in outer space... and Opportunity with it.")
	NotEnoughFuelErr = errors.New("Mars Lander spent all the fuel")
)

type CollisionError struct {
	From, To  Point
	LineIndex int
	Line      *Line
	Lander    Lander
}

func (c CollisionError) Error() string {
	return "Mars Lander crashed. Opportunity has been destroyed."
}

func (c CollisionError) Save() bool {
	lander := c.Lander
	if c.Line.Flat() {
		dSafeRotate := math.Abs(float64(lander.Rotate))
		if !(dSafeRotate < 10) {
			return false
		}
		dSafeUx := float32(math.Abs(float64(lander.Ux)) - 20)
		dSafeUy := float32(math.Abs(float64(lander.Uy)) - 40)
		return dSafeUx <= 0 && dSafeUy <= 0
	}
	return false
}

// Move возвращает сдвинутый марсоход спустя один такт симуляции
// Возвращает ошибку, если двигаться больше нельзя
func Move(lander Lander, command LanderCommand, surface *Surface) (Lander, error) {
	const dt float32 = 1
	curLineI, curLine := surface.LineFrom(0, lander.Point.X)
	// Симуляция
	lander.Apply(command)
	power := float32(lander.Power)
	lander.Fuel -= power
	dUx := power * Cos(lander.Rotate) * dt
	dUy := (power*Sin(lander.Rotate) - Gravity) * dt
	dX := dUx*dt/2 + lander.Ux*dt
	dY := dUy*dt/2 + lander.Uy*dt
	next := Point{lander.X + dX, lander.Y + dY}
	curLineI, curLine = surface.LineFrom(curLineI, next.X)
	lander.Ux += dUx
	lander.Uy += dUy
	// вылетели за границы
	if curLineI == -1 || next.Y > MaxY {
		return lander, OutOfBoundsErr
	}
	// врезались в поверхность
	if curLine.IntersectPath(lander.Point, next) {
		old := lander.Point
		lander.Point = next
		return lander, CollisionError{
			From:      old,
			To:        next,
			LineIndex: curLineI,
			Line:      curLine,
			Lander:    lander,
		}
	}
	// кончилось топливо
	if lander.Fuel < 0 {
		return lander, NotEnoughFuelErr
	}
	lander.Point = next
	return lander, nil
}
