package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/fs"
)

func tokenize(root string) map[string]int {
	BagOfWords := map[string]int{}
	filepath.WalkDir(root, func (path string, d fs.DirEntry, err error) (error) {
		if d.IsDir() {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			panic("couldn't open the file")
		}


		tokens := strings.Fields(string(content))
		for _, token := range tokens {
			BagOfWords[strings.ToUpper(token)] += 1
		}

		return nil

		/*
		for token, freq := range BagOfWords {
			fmt.Printf("|%s| => %f\n", token, float64(freq)/float64(totalCount));
		}
		*/
	})
	return BagOfWords
}

func main() {
	ham_tokens := tokenize("./enron1/")
	/*
	totalCount := 0
	for _, freq := range spam_tokens {
		totalCount += freq
	}
	*/

	for token, freq := range ham_tokens {
		fmt.Printf("|%s| => %d\n", token, freq);
	}
}
