package marslander

import (
	"fmt"
	"math"
)

// Lander - марсоход
type Lander struct {
	Point
	// the horizontal speed (in m/s), can be negative.
	Ux float32
	// the vertical speed (in m/s), can be negative.
	Uy float32
	// the quantity of remaining fuel in liters.
	Fuel float32
	// the rotation angle in degrees (-90 to 90).
	Rotate float32
	// the thrust Power (0 to 4).
	Power int
}

func (l *Lander) PrettyString() string {
	return fmt.Sprintf(
		`X=%fm, Y=%fm, HSpeed=%fm/s VSpeed=%fm/s
Fuel=%fl, Angle=%f°, Power=%d (%d.0m/s2)`,
		l.X, l.Y, l.Ux, l.Uy, l.Fuel, l.Rotate, l.Power, l.Power)
}

func (l Lander) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}

func (l *Lander) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buffer string
	if err := unmarshal(&buffer); err != nil {
		return fmt.Errorf("for Lander expected string: %w", err)
	}
	_, err := fmt.Sscan(buffer, l)
	return err
}

func (l *Lander) Scan(state fmt.ScanState, _ rune) error {
	_, err := fmt.Fscan(state, &l.X, &l.Y, &l.Ux, &l.Uy, &l.Fuel, &l.Rotate, &l.Power)
	return err
}

func roundF(x float32) int {
	return int(math.Floor(float64(x)))
}

func (l Lander) String() string {
	return fmt.Sprint(roundF(l.X), roundF(l.Y), roundF(l.Ux), roundF(l.Uy), roundF(l.Fuel), roundF(l.Rotate), l.Power)
}

// Команда марсоходу
type LanderCommand struct {
	Rotate, Power int
}

func (c *LanderCommand) Scan(state fmt.ScanState, _ rune) error {
	_, err := fmt.Fscanln(state, &c.Rotate, &c.Power)
	if err != nil {
		return err
	}
	if c.Rotate < -MaxRotate || c.Rotate > MaxRotate {
		return fmt.Errorf("rotation angle should be in [-90, 90], but actial value is %d", c.Rotate)
	}
	if c.Power < 0 || c.Power > MaxPower {
		return fmt.Errorf("power should be in [0, 4], but actial value is %d", c.Power)
	}
	return nil
}

func (l *Lander) Apply(c LanderCommand) {
	dRotate := math.Round(float64(c.Rotate)) - float64(l.Rotate)
	if math.Abs(dRotate) > 15 {
		l.Rotate += float32(math.Copysign(15, dRotate))
	} else {
		l.Rotate += float32(dRotate)
	}
	if l.Rotate > MaxRotate {
		l.Rotate = MaxRotate
	} else if l.Rotate < -MaxRotate {
		l.Rotate = -MaxRotate
	}
	if c.Power > l.Power {
		l.Power++
	} else if c.Power < l.Power {
		l.Power--
	}
	if l.Power < 0 {
		l.Power = 0
	}
	if l.Power > MaxPower {
		l.Power = MaxPower
	}
}
