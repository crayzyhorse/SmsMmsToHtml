package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func myUsage() {
	fmt.Printf("\n Convert xml to html with SMS and MMS media\n")
	fmt.Printf(" Made to work with the xml file produced\n")
	fmt.Printf(" by the SMS Backup & Restore app by Carbonite(TM)\n")
	fmt.Printf(" Under the app preferences:\n")
	fmt.Printf(" -BE SURE TO ENABLE Add Readable Date\n")
	fmt.Printf(" -BE SURE TO ENABLE Add Contact Name\n")
	fmt.Printf(" -BE SURE TO ENABLE Include MMS Messages\n")
	fmt.Printf(" -BE SURE TO ENABLE Include Emoji/Special Characters\n")
	fmt.Printf(" SMS Backup & Restore can be found in the\n")
	fmt.Printf(" Google Play Store, tested on version 9.74.1\n\n")
	fmt.Printf(" The template.html file can be edited to \n")
	fmt.Printf(" adjust color, layout, font, etc. of the messages,\n")
	fmt.Printf(" message bubbles, and background.\n\n")
	fmt.Printf(" View Readme file for more information.\n\n")
	fmt.Printf(" Usage: %s [file] [arg] ...\n", os.Args[0])
	fmt.Printf("  -n	: Outgoing Name or Number (default: local sender)\n\n")
	fmt.Printf(" example: %s 2019-06-30.xml\n", os.Args[0])
	fmt.Printf(" example: %s -n=Jennifer 2019-06-30.xml\n", os.Args[0])
	fmt.Printf(" example: %s -n=432-555-8765 2019-06-30.xml\n", os.Args[0])
	// Add Readable Date  Add Contact Name  Include Emoji/Special Characters
}

func main() {

	flag.Usage = myUsage

	outgoing := flag.String("n", "local sender", "Outgoing Name or Number")
	flag.Parse()

	fmt.Println("Working...")

	// Check for .xml as argument
	//if len(os.Args) < 2 || strings.ToLower(filepath.Ext(os.Args[1])) != ".xml" {
	if len(os.Args) < 2 || strings.ToLower(filepath.Ext(flag.Arg(0))) != ".xml" {
		fmt.Println("Please enter a valid .xml filename")
		return
	}

	// Get full (absoulute) path with filename
	filename, _ := filepath.Abs(flag.Arg(0))

	// Get directory only
	dir, err := filepath.Abs(filepath.Dir(flag.Arg(0)))
	if err != nil {
		fmt.Println("Fatal error, cannot resolve path from argument")
	}

	// Check path exists, possible redundant error
	// checking since the routine above should return
	// error if dir doesn't exist
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Path not valid")
			return
		}
	}

	// Check that the .xml file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File not found, Please check your path and filename", filename)
		return
	}

	fmt.Println("Outgoing message sender set to", *outgoing)
	fmt.Println("Found file ", flag.Arg(0))
	fmt.Println("Begin parsing...")

	// Parse the xml into html
	html, err := SplitXML(filename, *outgoing)
	if err != nil {
		fmt.Println("Cannot parse XML. Check you are using xml file from smsbackup")
		return
	}

	fmt.Println("Finished parsing...")

	// Create html subdirectory if it doesn't already exist
	htmldir := dir + string(filepath.Separator) + "html"
	_, err = os.Stat(htmldir)
	if os.IsNotExist(err) {
		errDir := os.Mkdir(htmldir, 0777)
		if errDir != nil {
			fmt.Println("Cannot create html directory.")
			return
		}
	}

	// Create the new file and write our html to it
	htmlname := filepath.Base(filename)
	htmlname = strings.TrimSuffix(htmlname, filepath.Ext(htmlname)) + ".html"

	fmt.Println("Attempting to write html to ./html/" + htmlname)

	newHtmlfile, err := os.Create(htmldir + string(filepath.Separator) + htmlname)
	if err != nil {
		fmt.Println("Cannot create new html file.")
		return
	}

	defer newHtmlfile.Close()

	fmt.Fprintf(newHtmlfile, html.String())

	fmt.Println("Success!")

}
