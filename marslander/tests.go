package marslander

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/RobolabGs2/tester/tester"
	svg "github.com/ajstarks/svgo/float"
	"gopkg.in/yaml.v2"
)

//go:embed tests.yaml
var defaultTestCases []byte
var DefaultTestCases []tester.TestCase

func init() {
	var testCases []TestCaseL
	if err := yaml.Unmarshal(defaultTestCases, &testCases); err != nil {
		log.Fatalln("Can't deserialize default test cases:", err)
	}
	DefaultTestCases = make([]tester.TestCase, len(testCases))
	for i := range testCases {
		DefaultTestCases[i] = testCases[i]
	}
}

type TestCaseL struct {
	Name    string
	Surface Surface `yaml:"surface"`
	Lander  Lander  `yaml:"lander"`
}

func (test TestCaseL) Title() string {
	return test.Name
}

type TestReport struct {
	Path   []Lander
	Rounds []tester.Turn
}

func (test TestCaseL) Run(solver *tester.Solver) (tester.Report, error) {
	report, err := test.run(solver)
	image := new(bytes.Buffer)
	DrawMarsLanderLaunch(image, test.Surface, report.Path)
	summary := "Success"
	if err != nil {
		summary = "Failure: " + err.Error()
	}
	return tester.Report{
		test.Name,
		err == nil,
		image.String(),
		report.Rounds,
		summary,
	}, nil
}

func (test TestCaseL) run(solver *tester.Solver) (TestReport, error) {
	lander := test.Lander
	surface := test.Surface
	report := TestReport{Path: []Lander{lander}}
	round, err := solver.Turn(surface.String(), nil)
	report.Rounds = append(report.Rounds, round)
	if err != nil {
		return report, err
	}
	for i := 0; true; i++ {
		var command LanderCommand
		round, err := solver.Turn(lander.String(), &command)
		round.Info = lander.PrettyString()
		report.Rounds = append(report.Rounds, round)
		if err != nil {
			return report, err
		}
		lander, err = Move(lander, command, &surface)
		report.Path = append(report.Path, lander)
		if err != nil {
			if collision := new(CollisionError); errors.As(err, collision) && collision.Save() {
				return report, nil
			}
			return report, err
		}
	}
	panic("unreachable")
}

func svgCircleWithTitle(s *svg.SVG, x float64, y float64, r float64, title string, styles string) {
	d := s.Decimals
	_, _ = fmt.Fprintf(s.Writer, `<circle cx="%.*f" cy="%.*f" r="%.*f" style=%q>`, d, x, d, y, d, r, styles)
	s.Title(title)
	_, _ = fmt.Fprintf(s.Writer, `</circle>`)
}

func DrawMarsLanderLaunch(file io.Writer, surface Surface, path []Lander) {
	s := svg.New(file)
	s.StartviewUnit(100, 100, "%", 0, 0, MaxX, MaxY)
	for i := 1; i < len(surface.Points); i++ {
		p1 := surface.Points[i-1]
		p2 := surface.Points[i]
		s.Line(float64(p1.X), float64(MaxY-p1.Y), float64(p2.X), float64(MaxY-p2.Y), fmt.Sprintf("stroke:%s; stroke-width:5;", "black"))
	}
	for i := 1; i < len(path); i++ {
		l1, l2 := path[i-1], path[i]
		s.Line(float64(l1.X), float64(MaxY-l1.Y), float64(l2.X), float64(MaxY-l2.Y), fmt.Sprintf("stroke:%s; stroke-width:5;", "red"))
		svgCircleWithTitle(s, float64(l1.X), float64(MaxY-l1.Y), 10, fmt.Sprintf("%d: %s", i, l1.PrettyString()), "stroke: blue; fill: green;")
	}
	s.End()
}

func TestCasesToHTMLTableBody(tests []TestCaseL, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "<tbody>"); err != nil {
		return err
	}
	for i, test := range tests {
		_, err := fmt.Fprintf(writer,
			"<tr><td>%d. %s</td><td><pre>%s\n%s</pre></td></tr>\n",
			i+1, test.Name, test.Surface, test.Lander)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(writer, "</tbody>")
	return err
}
