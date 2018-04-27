package lib

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"
)

func MarkdownReport(s *State, targets chan Host) string {
	var report bytes.Buffer
	t := time.Now()
	currTime := fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	var reportFile string
	if s.ProjectName != "" {
		reportFile = path.Join(s.ReportDirectory, fmt.Sprintf("%v_%v_Report.md", s.ProjectName, currTime))

	} else {
		reportFile = path.Join(s.ReportDirectory, fmt.Sprintf("%v_Report.md", currTime))
	}
	file, err := os.Create(reportFile)
	if err != nil {
		panic(err)
	}
	// Header
	report.WriteString(fmt.Sprintf("# Gograbber report - %v (%v)\n", s.ProjectName, currTime))
	for URLComponent := range targets {
		url := fmt.Sprintf("%v://%v:%v/%v\n", URLComponent.Protocol, URLComponent.HostAddr, URLComponent.Port, URLComponent.Path)
		report.WriteString(fmt.Sprintf("## %v\n", url))
		report.WriteString(fmt.Sprintf("![%v](../../%v)\n", URLComponent.ScreenshotFilename, URLComponent.ScreenshotFilename))

	}
	file.WriteString(report.String())
	return reportFile
}

func TextOutput(s *State) {

}

func JSONify(s *State) {

}

func SanitiseFilename(UnsanitisedFilename string) string {
	r := regexp.MustCompile("[0-9a-zA-Z-._]")
	return r.ReplaceAllString(UnsanitisedFilename, "[0-9a-zA-Z-._]")
}
