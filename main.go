package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/fs"
	"math"
)

func tokenize(message string) []string {
	tokens := strings.Fields(message)
	for i, val := range tokens {
		tokens[i] = strings.ToUpper(val)
	}

	return tokens
}

func fileToBow(path string, bow map[string]int) (error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	tokens := tokenize(string(content))
	for _, token := range tokens {
		bow[token] += 1
	}

	return nil
}

func BagOfWords(root string) (map[string]int, error) {
	bow := map[string]int{}
	err := filepath.WalkDir(root, func (path string, d fs.DirEntry, err error) (error) {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		err = fileToBow(path, bow)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return map[string]int{}, err
	}
	return bow, nil
}

func main() {
	BowHam, err := BagOfWords("./enron1/ham/")
	if err != nil {
		panic(err)
	}
	BowSpam, err := BagOfWords("./enron1/spam/")
	if err != nil {
		panic(err)
	}

	totalHam := 0
	for _, freq := range BowHam {
		totalHam += freq
	}

	totalSpam := 0
	for _, freq := range BowSpam {
		totalSpam += freq
	}

	EmailBow := map[string]int{}

	err = fileToBow("./enron2/ham/0004.1999-12-10.kaminski.ham.txt", EmailBow)
	if err != nil {
		panic(err)
	}

	prob_of_doc_ham := 0.0
	prob_of_doc_spam := 0.0

	for token := range EmailBow {
		if BowHam[token] == 0 {
			fmt.Printf("Ignored %v\n", token)
			continue
		}
		p := math.Log(float64(BowHam[token])/float64(totalHam))
		prob_of_doc_ham += p
		fmt.Printf("%v => %v\n", token, p)
	}

	for token := range EmailBow {
		if BowSpam[token] == 0 {
			fmt.Printf("Ignored %v\n", token)
			continue
		}
		p := math.Log(float64(BowSpam[token])/float64(totalSpam))
		prob_of_doc_spam += p
		fmt.Printf("%v => %v\n", token, p)
	}

	fmt.Printf("probability of document(HAM): %v\n", prob_of_doc_ham)
	fmt.Printf("probability of document(SPAM): %v\n", prob_of_doc_spam)

	fmt.Printf("len(ham) == %v\n", len(BowHam))
	fmt.Printf("len(spam) == %v\n", len(BowSpam))
}
