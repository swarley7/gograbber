package lib

import (
	"github.com/benbjohnson/phantomjs"
)

type State struct {
	Cookies                string
	Expanded               bool
	Extensions             []string
	FollowRedirect         bool
	PhantomProcess         phantomjs.Process
	ScreenshotQuality      int
	ScreenshotDirectory    string
	ReportDirectory        string
	ScanOutputDirectory    string
	ProjectName            string
	DirbustOutputDirectory string
	IncludeLength          bool
	NoStatus               bool
	Hosts                  StringSet
	InputFile              string
	Debug                  bool
	ExcludeList            []string
	Password               string
	Ports                  IntSet
	Jitter                 int
	Sleep                  float64
	Quiet                  bool
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
