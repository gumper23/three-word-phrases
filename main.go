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

// wget http://www.gutenberg.org/cache/epub/10/pg10.txt
// wget http://www.gutenberg.org/cache/epub/2009/pg2009.txt
// wget http://www.gutenberg.org/cache/epub/16328/pg16328.txt
// wget https://www.gutenberg.org/files/2600/2600-0.txt
// go build -o three-word-phrases main.go && ./three-word-phrases pg2009.txt pg16328.txt pg10.txt 2600-0.txt
func main() {
	const ranks = 10

	// Read from file names passed on the command line or stdin.
	if len(os.Args) == 1 {
		contents, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading stdin [%s]\n", err.Error())
		}
		results, err := threeWordPhrases(contents)
		printTopNResults(results, ranks)
	} else {
		ch := make(chan map[string]int)
		for _, file := range os.Args[1:] {
			go func(file string, ch chan map[string]int) {
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
				ch <- results
			}(file, ch)
		}
		for result := range ch {
			printTopNResults(result, ranks)
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
	longest := len("Three-Word Phrase")
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
	header := fmt.Sprintf("%-6s %-"+strconv.Itoa(longest+5)+"s %5s\n", "Rank", "Three-Word Phrase", "Count")
	fmt.Printf(header)
	fmt.Printf("%s\n", strings.Repeat("-", len(header)))
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

	// Replace all non-word characters with a space.
	re, err := regexp.Compile("[^\\w\\s]")
	if err != nil {
		return
	}
	s = re.ReplaceAllString(s, " ")

	// Count the three word phrases.
	re, err = regexp.Compile("\\s+")
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
