package lib

import (
	"net/http"
	"net/url"
	"os"
)

type State struct {
	Client         *http.Client
	Cookies        string
	Expanded       bool
	Extensions     []string
	FollowRedirect bool
	IncludeLength  bool
	Mode           string
	NoStatus       bool
	Debug          bool
	Password       string
	Ports          IntSet
	Printer        PrintResultFunc
	Processor      ProcessorFunc
	ProxyURL       *url.URL
	Quiet          bool
	Setup          SetupFunc
	ShowIPs        bool
	Protocol       string
	ShowCNAME      bool
	StatusCodes    IntSet
	StatusCodesIgn IntSet
	Screenshot     bool
	Threads        int
	URL            string
	UseSlash       bool
	Scan           bool
	UserAgent      string
	Username       string
	Verbose        bool
	Wordlist       string
	OutputFileName string
	OutputFile     *os.File
	SignalChan     chan os.Signal
	Terminate      bool
	StdIn          bool
	InsecureSSL    bool
}
