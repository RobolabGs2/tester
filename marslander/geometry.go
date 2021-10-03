package marslander

import "math"

// Нам нужно только ограниченное множество значений косинусов/синусов, можем предпросчитать
var _cos = [2*MaxRotate + 2]float32{}
var _sin = [2*MaxRotate + 2]float32{}

// Функция init вызывается в самом начале программы для каждого пакета
func init() {
	for i := range _cos {
		_cos[i] = float32(math.Cos(float64(i) / 180 * math.Pi))
		_sin[i] = float32(math.Sin(float64(i) / 180 * math.Pi))
	}
}

func Cos(rotate float32) float32 {
	return _cos[int(math.Round(float64(rotate+MaxRotate)))]
}

func Sin(rotate float32) float32 {
	return _sin[int(math.Round(float64(rotate+MaxRotate)))]
}

func sqrf(x float32) float32 {
	return x * x
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (from Point) DistanceTo(to Point) float64 {
	return math.Sqrt(float64(sqrf(to.X-from.X) + sqrf(to.Y-from.Y)))
}

type Line struct {
	segment []Point
	A, B, C float32
}

func MakeLine(segment []Point) Line {
	return Line{segment: segment,
		A: segment[0].Y - segment[1].Y,
		B: segment[1].X - segment[0].X,
		C: segment[0].X*segment[1].Y - segment[1].X*segment[0].Y,
	}
}

// WhereX - x левее, правее или в пределах линии?
func (l *Line) WhereX(x float32) int {
	if l.segment[0].X <= x {
		if x < l.segment[1].X {
			return 0
		}
		return +1
	}
	return -1
}

// RawDistance - подстановка точки в уравнение прямой
// > 0 - под линией
// < 0 - над линией
// = 0 - на линии
func (l *Line) RawDistance(p Point) float32 {
	return l.A*p.X + l.B*p.Y + l.C
}

func (l *Line) Underground(p Point) bool {
	return l.RawDistance(p) <= 0
}

// Точка пересечения линий
func (l *Line) IntersectLines(l2 Line) (Point, bool) {
	c := l.A*l2.B - l2.A*l.B
	if c == 0 {
		return Point{}, false
	}
	a := l.B*l2.C - l2.B*l.C
	b := l2.A*l.C - l.A*l2.C
	return Point{a / c, b / c}, true
}

// IntersectPath - Пересёк ли марсоход линию при перемещении из точки from в to
func (l *Line) IntersectPath(from Point, to Point) bool {
	if l.Underground(to) {
		// Место назначения под землёй - точно разбились
		return true
	}
	if !l.Underground(from) {
		// Старт и финиш оказались над землёй - точно не пересекаем
		return false
	}
	intersect, ok := l.IntersectLines(MakeLine([]Point{from, to}))
	if !ok {
		return false
	}
	return l.WhereX(intersect.X) == 0
}

func (l *Line) Flat() bool {
	return l.segment[0].Y == l.segment[1].Y
}
