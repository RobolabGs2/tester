package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/RobolabGs2/tester/exec"
	"github.com/RobolabGs2/tester/marslander"
	"github.com/RobolabGs2/tester/tester"
)

type TemplateData struct {
	Styles string
	Report []tester.Report
}

//go:embed report.tmpl
var htmlReportTemplate string

//go:embed styles.css
var defaultStyles string

func main() {
	openReport := flag.Bool("open", false, "Open report in default browser.")
	timeout := flag.Duration("timeout", time.Minute, "Timeout for one test.")
	customStylesPath := flag.String("styles", "", "Path to file with custom styles.")
	hour, min, sec := time.Now().Clock()
	filename := flag.String("output", fmt.Sprintf("mars-lander-report-%d_%d_%d.html", hour, min, sec), "File with report.")
	flag.Parse()
	solverPath := flag.Arg(0)
	if solverPath == "" {
		_, _ = fmt.Fprintln(os.Stderr, `Usage:
tester [flags] path/to/solution.exe`)
		flag.PrintDefaults()
		os.Exit(1)
	}
	reports := runTests(*timeout, solverPath, marslander.DefaultTestCases)
	file, err := os.Create(*filename)
	if err != nil {
		log.Fatalf("Can't create report file: %s\n", err)
	}
	tmplData := TemplateData{fmt.Sprintf("<style>%s</style>", defaultStyles), reports}
	if *customStylesPath != "" {
		tmplData.Styles = fmt.Sprintf(`<link href="%s" rel="stylesheet">`, *customStylesPath)
	}
	err = template.Must(template.New("").Parse(htmlReportTemplate)).Execute(file, tmplData)
	_ = file.Close()
	if err != nil {
		log.Fatalln("Failed to save report:", err)
	}
	fmt.Println("Report saved to", file.Name())
	if *openReport {
		if err := exec.OpenInBrowser(file.Name()); err != nil {
			log.Fatalf("Failed to open report %q: %s", file.Name(), err)
		}
	}
	return
}

func runTests(timeout time.Duration, solverPath string, tests []tester.TestCase) []tester.Report {
	reports := make([]tester.Report, len(tests))
	for i, test := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		report, err := tester.RunTestForSolverBinary(ctx, solverPath, test)
		cancel()
		if err != nil {
			log.Fatalf("Problems with test %d. %s: %s\n", i, test.Title(), err)
		}
		reports[i] = report
		fmt.Printf("%d. %s: %s\n", i, test.Title(), report.Summary)
	}
	return reports
}
