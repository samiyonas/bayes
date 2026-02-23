package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const THRESOLD = 200

func counter(bow map[string]int) int {
	total := 0
	for _, freq := range bow {
		if freq < THRESOLD {
			continue
		}
		total += freq
	}

	return total
}

func tokenize(message string) []string {
	tokens := strings.Fields(message)
	for i, val := range tokens {
		tokens[i] = strings.ToUpper(val)
	}

	return tokens
}

func fileToBow(path string, bow map[string]int) error {
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

func docProbabilityOverClass(ham, spam map[string]int, path string, totalHam, totalSpam, vocabSize int) (float64, float64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}

	tokens := tokenize(string(content))
	hamProb := 0.0
	spamProb := 0.0

	for _, token := range tokens {
		hamProb += math.Log(float64(ham[token]+1) / float64(totalHam+vocabSize))
		spamProb += math.Log(float64(spam[token]+1) / float64(totalSpam+vocabSize))
	}

	return hamProb, spamProb, nil
}

func BagOfWords(root string, bow map[string]int) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
		return err
	}
	return nil
}

func getVocabSize(ham, spam map[string]int) int {
	vocab := map[string]struct{}{}
	for word := range ham {
		vocab[word] = struct{}{}
	}

	for word := range spam {
		vocab[word] = struct{}{}
	}

	return len(vocab)
}

func main() {
	fmt.Printf("Training...\n")
	BowHam := map[string]int{}
	BowSpam := map[string]int{}

	for i := 1; i <= 5; i += 1 {
		path := fmt.Sprintf("./training_data/enron%v/%v/", i, "ham")
		err := BagOfWords(path, BowHam)
		if err != nil {
			panic(err)
		}
		path = fmt.Sprintf("./training_data/enron%v/%v/", i, "spam")
		err = BagOfWords(path, BowSpam)
		if err != nil {
			panic(err)
		}
	}

	totalHam := counter(BowHam)
	totalSpam := counter(BowSpam)
	total := totalHam + totalSpam
	ogHam := float64(totalHam) / float64(total)
	ogSpam := float64(totalSpam) / float64(total)
	hamHamCount := 0
	spamHamCount := 0
	hamSpamCount := 0
	spamSpamCount := 0

	vocabSize := getVocabSize(BowHam, BowSpam)

	fmt.Printf("classifying ham...\n")
	err := filepath.WalkDir("./training_data/enron6/ham/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return nil
		}

		hamProb, spamProb, err := docProbabilityOverClass(BowHam, BowSpam, path, totalHam, totalSpam, vocabSize)
		if err != nil {
			return err
		}

		hp := math.Log(ogHam) + hamProb
		sp := math.Log(ogSpam) + spamProb

		if hp > sp {
			hamHamCount += 1
		} else {
			spamHamCount += 1
		}

		return nil

	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("	Ham: %v\n", hamHamCount)
	fmt.Printf("	Spam: %v\n", spamHamCount)

	fmt.Printf("classifying spam...\n")
	err = filepath.WalkDir("./training_data/enron6/spam/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return nil
		}

		hamProb, spamProb, err := docProbabilityOverClass(BowHam, BowSpam, path, totalHam, totalSpam, vocabSize)
		if err != nil {
			return err
		}

		hp := math.Log(ogHam) + hamProb
		sp := math.Log(ogSpam) + spamProb

		if hp > sp {
			hamSpamCount += 1
		} else {
			spamSpamCount += 1
		}

		return nil

	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("	Ham: %v\n", hamSpamCount)
	fmt.Printf("	Spam: %v\n", spamSpamCount)
}
