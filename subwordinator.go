package main

import (
	"fmt"
	"bufio"
	"io"
	"os"
	"sync"
)

// word per line thing.
func getWordList (in io.Reader) map[string]byte {
	reader := bufio.NewReader(in)
	var nextWord func () (string, error)
	rv := make(map[string]byte)

	
	nextWord = func () (string, error) {
		line, notDone, err  := reader.ReadLine()
		if err != nil {
			return "", err
		}

		for notDone {
			var nextChunk []byte
			nextChunk, notDone, err = reader.ReadLine()
			if err != nil {
				return "", err
			}
			line = append(line, nextChunk...)
		}
		return string(line), err
	}

	line, linerr := nextWord()
	
	for linerr == nil {
		rv[line] = 1
		line, linerr = nextWord()
	}

	return rv
}

func wordsFromFile (filename string) map[string]byte {
	file, err := os.Open(filename)
	if err != nil {
		panic("could not open the file")
	}
	defer file.Close()
	return getWordList(file)
}

func allDrop(s string) [] string {
	results := make([]string, 0)
	for i := 0; i < len(s); i++ {
		results = append(results, string(append([]byte(s[0:i]), []byte(s[i+1:])...)))
	}

	return results
}
		
func main () {
	words := wordsFromFile("/usr/share/dict/words")

	interesting := os.Args[1]

	var mainWg sync.WaitGroup
	mainWg.Add(1)

	var hahahaha func(string, *sync.WaitGroup)
	hahahaha = func(s string, inWg *sync.WaitGroup) {
		defer inWg.Done()
		if s == "" {
			return
		}
		if _,ok := words[s]; ok {
			fmt.Println(s)
			ad := allDrop(s)
			if len(ad) > 0 {
				var wg sync.WaitGroup
				wg.Add(len(ad))
				for _, each := range ad {
					go hahahaha(each, &wg)
				}
				wg.Wait()
			}
		} 
	}
	go hahahaha(interesting, &mainWg)

	mainWg.Wait()
}
	
