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

func main() {
	// Read from file names passed on the command line or stdin.
	if len(os.Args) == 1 {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading stdin [%s]\n", err.Error())
		}
		results, err := threeWordPhrases(contents)
		printTopNResults(results, 100)
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
			printTopNResults(results, 100)
		}
	}
}

// printTopNResults prints the top N results to stdout.
func printTopNResults(counts map[string]int, n int) {
	type kv struct {
		Key   string
		Value int
	}

	// Sort counts by descending value.
	var keyValues []kv
	longest := 6
	for k, v := range counts {
		if len(k) > longest {
			longest = len(k)
		}
		keyValues = append(keyValues, kv{k, v})
	}
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Value > keyValues[j].Value
	})

	// Print the top n counts.
	format := "[%03d]: [%-" + strconv.Itoa(longest) + "s] => [%d]\n"
	hdr := fmt.Sprintf("%-6s %-"+strconv.Itoa(longest+5)+"s %5s\n", "Rank", "Phrase", "Count")
	fmt.Printf(hdr)
	fmt.Printf("%s\n", strings.Repeat("-", len(hdr)))
	for i, kv := range keyValues {
		if i == n {
			break
		}
		fmt.Printf(format, i+1, kv.Key, kv.Value)
	}
}

// threeWordPhrases counts the three word phrases in contents.
func threeWordPhrases(contents []byte) (counts map[string]int, err error) {
	counts = make(map[string]int)

	// Convert contents to a lowercase string.
	s := strings.ToLower(string(contents))

	// Replace all non-word characters with a space.
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

	// Count the three word phrases.
	words := strings.Split(s, " ")
	for i := 0; i < len(words)-2; i++ {
		phrase := words[i] + " " + words[i+1] + " " + words[i+2]
		counts[phrase]++
	}
	return
}
