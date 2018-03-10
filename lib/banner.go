package lib

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// PrintBanner prints the banner
func PrintBanner(s *State) {
	var banner string
	banner = `                                                 $$\       $$\                           
                                                 $$ |      $$ |                          
 $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$$$$$\  $$$$$$$\  $$$$$$$\   $$$$$$\   $$$$$$\  
$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ \____$$\ $$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ 
$$ /  $$ |$$ /  $$ |$$ /  $$ |$$ |  \__|$$$$$$$ |$$ |  $$ |$$ |  $$ |$$$$$$$$ |$$ |  \__|
$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |     $$  __$$ |$$ |  $$ |$$ |  $$ |$$   ____|$$ |      
\$$$$$$$ |\$$$$$$  |\$$$$$$$ |$$ |     \$$$$$$$ |$$$$$$$  |$$$$$$$  |\$$$$$$$\ $$ |      
 \____$$ | \______/  \____$$ |\__|      \_______|\_______/ \_______/  \_______|\__|      
$$\   $$ |          $$\   $$ |                                                           
\$$$$$$  |          \$$$$$$  |                                                           
 \______/            \______/`

	var version = "0.1a"
	var author = "swarley7 (@swarley777)"
	if !s.Quiet {
		g := color.New(color.FgGreen, color.Bold)

		fmt.Printf("%v\n", strings.Replace(banner, "$", g.Sprintf("$"), -1))
		fmt.Printf("%v\n", LeftPad2Len(fmt.Sprintf("Author: %v", author), " ", 89))
		fmt.Printf("%v\n", LeftPad2Len(fmt.Sprintf("Version: %v", version), " ", 89))
	}

}
