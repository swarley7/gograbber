package lib

import (
	"time"

	"github.com/benbjohnson/phantomjs"
)

type State struct {
	Cookies                string
	Expanded               bool
	Extensions             []string
	FollowRedirect         bool
	PhantomProcesses       []phantomjs.Process
	ScreenshotQuality      int
	ScreenshotDirectory    string
	ReportDirectory        string
	ScanOutputDirectory    string
	ProjectName            string
	DirbustOutputDirectory string
	IncludeLength          bool
	DisplayEnvVar          string
	Ratio                  float64
	Soft404Detection       bool
	Soft404Method          int
	NoStatus               bool
	Hosts                  StringSet
	InputFile              string
	Debug                  bool
	ExcludeList            []string
	Password               string
	Ports                  IntSet
	Jitter                 int
	Sleep                  float64
	Timeout                time.Duration
	VerbosityLevel         int
	ShowIPs                bool
	Protocols              StringSet
	StatusCodes            IntSet
	StatusCodesIgn         IntSet
	ImgX                   int
	ImgY                   int
	Screenshot             bool
	Threads                int
	URLFile                string
	URLProvided            bool
	Dirbust                bool
	SingleURL              string
	PhantomJSPath          string
	URLComponents          []Host
	Paths                  StringSet
	UseSlash               bool
	Scan                   bool
	UserAgent              string
	Username               string
	Verbose                bool
	Wordlist               string
	OutputDirectory        string
}
