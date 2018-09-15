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
	"sync"
)

// wget http://www.gutenberg.org/cache/epub/10/pg10.txt
// wget http://www.gutenberg.org/cache/epub/2009/pg2009.txt
// wget http://www.gutenberg.org/cache/epub/16328/pg16328.txt
// wget https://www.gutenberg.org/files/2600/2600-0.txt
// go build -o three-word-phrases main.go && ./three-word-phrases pg2009.txt pg16328.txt pg10.txt 2600-0.txt
func main() {
	// Read from file names passed on the command line or stdin.
	if len(os.Args) == 1 {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading stdin [%s]\n", err.Error())
		}
		results, err := threeWordPhrases(contents)
		printTopNResults(results, 10)
	} else {
		var wg sync.WaitGroup
		for i := 1; i < len(os.Args); i++ {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				contents, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n\n", err.Error())
					return
				}
				results, err := threeWordPhrases(contents)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n\n", err.Error())
					return
				}
				printTopNResults(results, 10)
			}(os.Args[i])
		}
		wg.Wait()
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
		length := len(k)
		if length > longest {
			longest = length
		}
		keyValues = append(keyValues, kv{k, v})
	}
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Value > keyValues[j].Value
	})

	// Print the top n counts.
	format := "[%03d]: [%-" + strconv.Itoa(longest) + "s] => [%d]\n"
	hdr := fmt.Sprintf("%-6s %-"+strconv.Itoa(longest+5)+"s %5s\n", "Rank", "Three-Word Phrase", "Count")
	fmt.Printf(hdr)
	fmt.Printf("%s\n", strings.Repeat("-", len(hdr)))
	for i, kv := range keyValues {
		if i == n {
			break
		}
		fmt.Printf(format, i+1, kv.Key, kv.Value)
	}
	fmt.Println()
}

// threeWordPhrases counts the three word phrases in contents.
func threeWordPhrases(contents []byte) (counts map[string]int, err error) {
	counts = make(map[string]int)

	// Convert contents to a lowercase string.
	s := strings.ToLower(string(contents))

	// Replace all non-word characters with a space, except for an apostrophe or hypen.
	re, err := regexp.Compile("[^[:alnum:][:space:]'-]")
	if err != nil {
		return
	}
	s = re.ReplaceAllString(s, " ")

	// Count the three word phrases.
	re, err = regexp.Compile("[[:space:]]+")
	if err != nil {
		return
	}
	words := re.Split(s, -1)
	for i := 0; i < len(words)-2; i++ {
		phrase := fmt.Sprintf("%s %s %s", words[i], words[i+1], words[i+2])
		counts[phrase]++
	}
	return
}
