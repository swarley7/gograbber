package lib

import (
	"net/url"
	"os"
)

type State struct {
	Cookies         string
	Expanded        bool
	Extensions      []string
	FollowRedirect  bool
	IncludeLength   bool
	NoStatus        bool
	Hosts           StringSet
	InputFile       string
	Debug           bool
	ExcludeList     []string
	Password        string
	Ports           IntSet
	Jitter          int
	Sleep           float64
	ProxyURL        *url.URL
	Quiet           bool
	ShowIPs         bool
	Protocols       StringSet
	StatusCodes     IntSet
	StatusCodesIgn  IntSet
	Screenshot      bool
	Threads         int
	URLFile         string
	Dirbust         bool
	SingleURL       string
	URLComponents   []Host
	Paths           StringSet
	UseSlash        bool
	Scan            bool
	UserAgent       string
	Username        string
	Verbose         bool
	Wordlist        string
	OutputDirectory string
	SignalChan      chan os.Signal
	Terminate       bool
}
