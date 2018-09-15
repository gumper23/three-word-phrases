package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type kv struct {
	Key   string
	Value int
}

func main() {
	// Read from file names passed on the command line or stdin.
	if len(os.Args) == 1 {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading stdin [%s]\n", err.Error())
		}
		results, err := threeWordPhrases(contents)
		printResults(results)
	} else {
		for i := 1; i < len(os.Args); i++ {
			contents, err := ioutil.ReadFile(os.Args[i])
			if err != nil {
				log.Fatalf("error reading %s: %s", os.Args[i], err.Error())
			}
			results, err := threeWordPhrases(contents)
			if err != nil {
				log.Fatalf("error parsing %s: %s", os.Args[i], err.Error())
			}
			printResults(results)
		}
	}
}

func printResults(counts map[string]int) {
	// Sort counts by descending value.
	var keyValues []kv
	longest := 0
	for k, v := range counts {
		if len(k) > longest {
			longest = len(k)
		}
		keyValues = append(keyValues, kv{k, v})
	}
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Value > keyValues[j].Value
	})

	// Print the top 100 counts.
	format := "[%03d]: [%-" + strconv.Itoa(longest) + "s] => [%d]\n"
	fmt.Printf("%s", format)
	for i, kv := range keyValues {
		if i == 100 {
			break
		}
		fmt.Printf(format, i+1, kv.Key, kv.Value)
	}
}

// threeWordPhrases counts the three word phrases in contents
func threeWordPhrases(contents []byte) (counts map[string]int, err error) {
	counts = make(map[string]int)

	// Convert the contents to lowercase.
	s := strings.ToLower(string(contents))

	// Replace all punctuation characters with a space.
	re, err := regexp.Compile("\\W+")
	if err != nil {
		return
	}
	s = re.ReplaceAllString(s, " ")

	// Replace multiple whitespace characters with a single space.
	re, err = regexp.Compile("\\s+")
	if err != nil {
		return
	}
	s = re.ReplaceAllString(s, " ")

	// Count the three word phrases
	words := strings.Split(s, " ")
	for i := 0; i < len(words)-2; i++ {
		phrase := words[i] + " " + words[i+1] + " " + words[i+2]
		counts[phrase]++
	}
	return
}
