//Compare TIBCO designtimelib
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	//Read 'new' file
	flag.Parse() //Parse command line arguments ie. the file names
	fp1 := strings.Replace(flag.Arg(0), "\\", "/", -1)
	fi1, err := os.Open(fp1) //Open file 1 which is interpeded as the new/modified file
	if err != nil {
		panic(err)
	}
	defer fi1.Close()
	r1 := bufio.NewReader(fi1)

	//Put 'new' file lines in slice
	s1 := []string{} //Create an empty slice
	l1, err := r1.ReadString('\n')
	for err == nil {
		s1 = append(s1, l1)
		l1, err = r1.ReadString('\n')
	}
	if err != nil && err != io.EOF {
		fmt.Print(err)
	}

	//Read 'old' file
	fp2 := strings.Replace(flag.Arg(1), "\\", "/", -1)
	fi2, err := os.Open(fp2) //Open file 2 which is interpeded as the old/base file file
	if err != nil {
		panic(err)
	}
	defer fi2.Close()
	r2 := bufio.NewReader(fi2)

	//Put 'old' file lines in slice
	s2 := []string{} //Create an empty slice
	l2, err := r2.ReadString('\n')
	for err == nil {
		s2 = append(s2, l2)
		l2, err = r2.ReadString('\n')
	}
	if err != nil && err != io.EOF {
		fmt.Print(err)
	}
	
	//Compile all regular expressions
	rgxComment, _ := regexp.Compile("^#")
	rxSeqNoTrim, _ := regexp.Compile("^[0-9]*?=")
	rxVerTrim, _ := regexp.Compile("_[0-9].[0-9].[0-9].projlib\\\\=$")
	rxVer, _ := regexp.Compile("[0-9].[0-9].[0-9]")

	for _, v1 := range s1 { //Loop over all entries in file 1, 'v' stands for value
		e1 := strings.TrimSpace(v1)                  //'e' stands for entry
		ver1 := rxVer.FindString(e1)                 //'ver' stands for version
		compV1 := strings.Replace(ver1, ".", "", -1) //'compVer' stands for compareble version
		if !rgxComment.MatchString(e1) {             //Only go comparing if entry is no comment
			isFound := false
			for _, v2 := range s2 { //Loop over all entries in file 2 in order to find a match to the entry of file 1
				e2 := strings.TrimSpace(v2)
				ver2 := rxVer.FindString(e2)                 //'ver' stands for version
				compV2 := strings.Replace(ver2, ".", "", -1) //'compVer' stands for compareble version
				trimedSeq := rxSeqNoTrim.ReplaceAllString(e2, "")
				trimedVer := rxVerTrim.ReplaceAllString(e2, "")
				trimedSeqVer := rxVerTrim.ReplaceAllString(trimedSeq, "")
				switch {
				case strings.Contains(e1, e2):
					isFound = true
				case strings.Contains(e1, trimedSeq): //When only the sequence has changed
					//fmt.Println("M " + e1) //M stands for moved
					isFound = true
				case strings.Contains(e1, trimedVer) && compV1 < compV2: //When only the version number of the lib has been decreased
					fmt.Println("D " + e1 + " (old ver.: " + ver2 + ")") //D stands for downgrade
					isFound = true
				case strings.Contains(e1, trimedVer) && compV1 > compV2: //When only the version number of the lib has been increased
					fmt.Println("U " + e1 + " (old ver.: " + ver2 + ")") //U stands for upgrade
					isFound = true
				case strings.Contains(e1, trimedSeqVer) && compV1 < compV2: //When the version number of the lib has been decreased and the sequence has changed
					fmt.Println("D " + e1 + " (old ver.: " + ver2 + ")") //Changing of sequence is ignored and only the downgrade is mentioned
					isFound = true
				case strings.Contains(e1, trimedSeqVer) && compV1 > compV2: //When the version number of the lib has been increased and the sequence has changed
					fmt.Println("U " + e1 + " (old ver.: " + ver2 + ")") //Changing of sequence is ignored and only the update is mentioned
					isFound = true
				}
			}
			if isFound == false { //If no case set the isFound to true then it's still false and thus the entry isn't found in file 2
				fmt.Println("A " + e1) //A stands for added
			}
		}
	}

	for _, v2 := range s2 {
		e2 := strings.TrimSpace(v2)
		if !rgxComment.MatchString(e2) { //Only go comparing if entry is no comment
			isFound := false
			for _, v1 := range s1 {
				trimed := rxVerTrim.ReplaceAllString(rxSeqNoTrim.ReplaceAllString(strings.TrimSpace(v1), ""), "")
				if strings.Contains(e2, trimed) {
					isFound = true
				}
			}
			if isFound == false {
				fmt.Println("R " + e2) //R stands for removed
			}
		}
	}
	
}
