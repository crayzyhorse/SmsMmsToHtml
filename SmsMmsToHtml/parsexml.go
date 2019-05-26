package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
Our xml file does not unmarshal very well. There are
attributes that should be elements of their own,
and the mms messages are listed almost as an afterthought
at the end of the file.  We can use encoding/xml
to unmarshal the mms, however we will go about disecting
the sms portion through other means
*/

// Smses struct
type Smses struct {
	XMLName xml.Name `xml:"smses"`
	Count   string   `xml:"count,attr"`
	Mms     []Mms    `xml:"mms"`
}

// Mms struct, this contains the attributes
// and additional structs needed for our html
type Mms struct {
	XMLName      xml.Name `xml:"mms"`
	Timestamp    string   `xml:"date,attr"`
	SentBySender string   `xml:"msg_box,attr"`
	ContactNum   string   `xml:"address,attr"`
	Date         string   `xml:"readable_date,attr"`
	ContactName  string   `xml:"contact_name,attr"`
	Parts        Parts    `xml:"parts"`
}

// Parts struct, child of Mms struct
type Parts struct {
	XMLName xml.Name `xml:"parts"`
	Part    []Part   `xml:"part"`
}

// Part struct, Mms >> Parts >> Part
type Part struct {
	XMLName  xml.Name `xml:"part"`
	ImgName  string   `xml:"img src,attr"`
	Seq      string   `xml:"seq,attr"`
	Ct       string   `xml:"ct,attr"`
	ImgName2 string   `xml:"cl,attr"`
	Data     string   `xml:"data,attr"`
	Text     string   `xml:"text,attr"`
}

// SplitXML function to seperate our sms and mms
func SplitXML(filename string, outgoing string) (*bytes.Buffer, error) {

	// Prepare our html output file
	// start by opening our template
	htmlFile, err := os.Open("template.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	// Read our template into the html buffer
	html := new(bytes.Buffer)
	html.ReadFrom(htmlFile)

	// Open our xml file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Start reading from our xml file with a reader.
	reader := bufio.NewReader(file)

	// Initialize our sms and mms blocks
	smsblock := new(bytes.Buffer)
	mmsblock := new(bytes.Buffer)

	countsms := 0

	var line string
	for {

		// Read XML file line by line with ReadString('\n')
		line, err = reader.ReadString('\n')

		// Remove blank space at beginning of line
		line = strings.TrimLeft(line, " ")

		// Seperate lines based on prefix
		if strings.HasPrefix(line, "<sms protocol") {
			smsblock.WriteString(line)
			countsms++
		} else {
			mmsblock.WriteString(line)
		}

		if err != nil {
			break
		}
	}

	// We can close our xml file now
	file.Close()

	smsLines := Scanit(smsblock)
	newSmsTime := Between(smsLines[0], "date=\"", "\" type")
	senderName := Between(smsLines[0], "contact_name=\"", "\" />")
	senderNum := Between(smsLines[0], "address=\"", "\" date")

	html.WriteString("\n<title>" + senderName + "</title></head>\n<body><div class=\"parent\"><br>")

	// read our mmsblock as a byte array.
	byteValue, _ := ioutil.ReadAll(mmsblock)

	var smses Smses

	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'smses' which we defined above
	xml.Unmarshal(byteValue, &smses)

	totalMessages, err := strconv.Atoi(smses.Count)
	if err != nil {
		fmt.Println("int to string conversion error, var totalMessages")
	}

	newMmsTime := smses.Mms[0].Timestamp
	s := 0
	m := 0

	for i := 0; i < totalMessages; i++ {

		if newSmsTime <= newMmsTime {

			if s >= countsms {
				break
			}

			if Between(smsLines[s], "type=\"", "\" sub") == "1" {
				html.WriteString("\n<div class=\"inBubble\"><div><p class=\"splitP\">" + senderName + " " + senderNum)
			} else {
				html.WriteString("\n<div class=\"outBubble\"><div><p class=\"splitP\">" + outgoing)
			}
			html.WriteString(":<span>" + Between(smsLines[s], "readable_date=\"", "\" contact") + "</span></p></div>\n")
			html.WriteString("\t<div class=\"txt\"><p>\n" + "\t" + Between(smsLines[s], "body=\"", "\" toa") + "\n\t</p></div>\n")
			html.WriteString("</div><br>")
			s++
			if s < countsms {
				newSmsTime = Between(smsLines[s], "date=\"", "\" type")
			}

		} else {

			if smses.Mms[m].SentBySender == "1" {
				html.WriteString("\n<div class=\"inBubble\"><div><p class=\"splitP\">" + senderName + " " + senderNum)
			} else {
				html.WriteString("\n<div class=\"outBubble\"><div><p class=\"splitP\">" + outgoing)
			}
			html.WriteString(":<span>" + smses.Mms[m].Date + "</span></p></div>\n")

			for b := 0; b < len(smses.Mms[m].Parts.Part); b++ {
				if smses.Mms[m].Parts.Part[b].Ct == "text/plain" {
					html.WriteString("\n\t<div class=\"txt\"><p>" + smses.Mms[m].Parts.Part[b].Text + "</p></div>")
				} //<p><img src="  changed to   <p><video width="560" src="
				if (smses.Mms[m].Parts.Part[b].Seq == "0") && (smses.Mms[m].Parts.Part[b].Ct != "text/plain") {
					if strings.HasPrefix(smses.Mms[m].Parts.Part[b].Ct, "video") { // ct="video*
						html.WriteString("\t<div class=\"imgBox\"><p><video width=\"560\" src=\"data:" + smses.Mms[m].Parts.Part[b].Ct + ";base64, " + smses.Mms[m].Parts.Part[b].Data + "\" alt=\"" + smses.Mms[m].Parts.Part[b].ImgName2 + "\" controls /></p></div>")
					} else {
						html.WriteString("\t<div class=\"imgBox\"><p><img src=\"data:" + smses.Mms[m].Parts.Part[b].Ct + ";base64, " + smses.Mms[m].Parts.Part[b].Data + "\" alt=\"" + smses.Mms[m].Parts.Part[b].ImgName2 + "\" /></p></div>")
					}
				}
			}

			html.WriteString("\n</div><br>")
			m++
			if m < len(smses.Mms) {
				newMmsTime = smses.Mms[m].Timestamp
			} else {
				newMmsTime = "2000000000000"
			}
		}
	}

	html.WriteString("\n</div></body></html>")

	return html, err

}
