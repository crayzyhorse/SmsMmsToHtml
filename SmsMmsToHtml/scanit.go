package main

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

// Scanit func takes the smsblock and reads each line
// checking for html escaped characters, converting them if found,
// and then adding each line to the array string
// our xml file may have emojis that are represented
// as escaped html such as "💋" represented as &#55357;&#56459;
// if the emoji doesn't show in your editor try searching the
// escaped html. I have not had good success converting these
// with go, for now we're simply converting &#10; to <br>
func Scanit(smsblock *bytes.Buffer) []string {
	smsScanner := bufio.NewScanner(smsblock)
	escape1 := regexp.MustCompile("&#10;")
	var smsLines []string
	var rawLine string
	var rawString []string

	for smsScanner.Scan() {
		rawLine = smsScanner.Text()

		if escape1.FindAllString(rawLine, -1) != nil {
			rawString = escape1.FindAllString(rawLine, -1)
			for r := 0; r < len(rawString); r++ {
				rawLine = strings.Replace(rawLine, rawString[r], strings.Replace((rawString[r]), "&#10;", "<br>", -1), 1)
			}
		}
		smsLines = append(smsLines, rawLine)
	}
	return smsLines
}

// Between func for returning a string between a
// specified startng point and an ending point
// we'll be using it to find specific attributes of
// our sms messages
func Between(originalString string, startPoint string, endPoint string) string {
	posFirst := strings.Index(originalString, startPoint)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(originalString, endPoint)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(startPoint)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return originalString[posFirstAdjusted:posLast]
}
