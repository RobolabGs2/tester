package marslander

import (
	"bytes"
	"fmt"
)

// Surface - хранит поверхность земли
type Surface struct {
	Points []Point `json:"points"`
	Lines  []Line  `json:"-"`
}

func (s *Surface) Init() {
	s.Lines = make([]Line, len(s.Points)-1)
	for i := range s.Lines {
		s.Lines[i] = MakeLine(s.Points[i : i+2])
	}
}

// Flat Поиск плоской поверхности
func (s Surface) Flat() *Line {
	for i := 0; i < len(s.Lines); i++ {
		if s.Lines[i].Flat() {
			return &s.Lines[i]
		}
	}
	panic("Not found a flat place!")
}

// LineFrom ищет принадлежность координаты x к линии поверхности,
// считая, что марсоход не улетит далеко от текущей линии номер i
func (s Surface) LineFrom(i int, x float32) (int, *Line) {
	// За границами
	if x < 0 || x > s.Points[len(s.Points)-1].X {
		return -1, nil
	}
	// dI - смещение влево или вправо относительно текущей линии
	for dI := s.Lines[i].WhereX(x); dI != 0; dI = s.Lines[i].WhereX(x) {
		i += dI
	}
	return i, &s.Lines[i]
}

// Scan - реализуем интерфейс для чтения из потока
func (s *Surface) Scan(state fmt.ScanState, _ rune) error {
	var surfaceN int
	if _, err := fmt.Fscan(state, &surfaceN); err != nil {
		return err
	}
	s.Points = make([]Point, surfaceN)
	for i := 0; i < surfaceN; i++ {
		_, err := fmt.Fscan(state, &s.Points[i].X, &s.Points[i].Y)
		if err != nil {
			return err
		}
	}
	s.Init()
	return nil
}

func (s Surface) String() string {
	buffer := new(bytes.Buffer)
	if _, err := fmt.Fprintln(buffer, len(s.Points)); err != nil {
		panic(err)
	}
	for _, point := range s.Points {
		_, err := fmt.Fprintln(buffer, point.X, point.Y)
		if err != nil {
			panic(err)
		}
	}
	res := buffer.String()
	return res[:len(res)-1] // cut new line
}

func (s Surface) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

func (s *Surface) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buffer string
	if err := unmarshal(&buffer); err != nil {
		return fmt.Errorf("for Surface expected string: %w", err)
	}
	_, err := fmt.Sscan(buffer, s)
	return err
}
