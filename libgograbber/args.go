package libgograbber

import (
	"net/http"
	"net/url"
	"os"
)

// Contains State that are read in from the command
// line when the program is invoked.
type State struct {
	Client         *http.Client
	Cookies        string
	Expanded       bool
	Extensions     []string
	FollowRedirect bool
	IncludeLength  bool
	Mode           string
	NoStatus       bool
	Password       string
	Printer        PrintResultFunc
	Processor      ProcessorFunc
	ProxyUrl       *url.URL
	Quiet          bool
	Setup          SetupFunc
	ShowIPs        bool
	ShowCNAME      bool
	StatusCodes    IntSet
	Threads        int
	Url            string
	UseSlash       bool
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
