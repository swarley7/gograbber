package lib

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"
)

func buildResponseHeader(header *http.Response) string {
	// respHeader = fmt.Sprintf("%v %v", header.Proto, header.Status)
	// for k, v := range header.Header {
	// 	respHeader = fmt.Sprintf("%v\n%v: %v", respHeader, k, v[0])
	// }
	buf := new(bytes.Buffer)
	header.Write(buf)
	return buf.String()
}

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
		report.WriteString("### Screenshot\n")
		report.WriteString(fmt.Sprintf("![%v](../../%v)\n", URLComponent.ScreenshotFilename, URLComponent.ScreenshotFilename))
		if URLComponent.HTTPResp != nil {
			report.WriteString("### Response Headers\n")

			report.WriteString(fmt.Sprintf("```\n%v\n```\n", buildResponseHeader(URLComponent.HTTPResp)))
		}

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
