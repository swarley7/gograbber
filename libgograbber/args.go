package libgograbber

import (
	"net/http"
	"time"

	"github.com/swarley7/phantomjs"
)

type State struct {
	Canary                 string
	Cookies                string
	Expanded               bool
	Extensions             []string
	FollowRedirect         bool
	PhantomProcesses       []phantomjs.Process
	ScreenshotQuality      int
	ScreenshotDirectory    string
	ScreenshotFileType     string
	FollowRedirects        bool
	ReportDirectory        string
	HostHeaders            StringSet
	ScanOutputDirectory    string
	ProjectName            string
	DirbustOutputDirectory string
	IncludeLength          bool
	NumPhantomProcs        int
	HttpClient             *http.Client
	HTTPResponseDirectory  string
	Ratio                  float64
	Soft404Detection       bool
	Soft404Method          int
	PrefetchedHosts        map[string]bool
	Soft404edHosts         map[string]bool
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
	IgnoreSSLErrors        bool
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
	Targets                chan Host
	Paths                  StringSet
	UseSlash               bool
	Scan                   bool
	UserAgent              string
	Username               string
	OutputChan             chan string
	Verbose                bool
	Wordlist               string
	Version                string
	OutputDirectory        string
	StartTime              time.Time
}
