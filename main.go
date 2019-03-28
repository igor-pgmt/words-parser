package main

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var resultingFile string
var unknownWordsDir string
var knownWordsDir string
var help bool

func checkFlags() {
	flag.StringVar(&resultingFile, "resultingFile", "result.csv", "Please, specify path to resulting file. [string] (default: ./result.csv")
	flag.StringVar(&unknownWordsDir, "unknownWordsDir", "./unknownWords/", "Please, specify path to directory with files with unknown words. [string] (default: ./unknownWords/)")
	flag.StringVar(&knownWordsDir, "knownWordsDir", "./knownWords/", "Please, specify path to directory with files with known words. [string] (default: ./knownWords/)")
	flag.BoolVar(&help, "help", false, "Please, specify to show this message. [bool] (default: false")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(1)
	}
}

func createDir(dir *string) (err error) {

	lastChar := (*dir)[len(*dir)-1:]
	if lastChar != string(os.PathSeparator) {
		*dir += string(os.PathSeparator)
	}

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		err = os.Mkdir(*dir, os.ModePerm)
	}

	return err
}

func main() {

	checkFlags()

	err := createDir(&unknownWordsDir)
	checkErr(err)
	err = createDir(&knownWordsDir)
	checkErr(err)

	unknownWords := getWordsMap(unknownWordsDir)
	knownWords := getWordsMap(knownWordsDir)
	removeKnownWords(unknownWords, knownWords)
	sortedResult := rankByWordCount(unknownWords)
	writeCSV(sortedResult)
}

func getWordsMap(dir string) *sync.Map {
	var wordsMap sync.Map
	files, err := ioutil.ReadDir(dir)
	checkErr(err)
	wg := &sync.WaitGroup{}
	for _, f := range files {
		wg.Add(1)
		go func(f os.FileInfo, wordsMap *sync.Map, wg *sync.WaitGroup) {
			bytes, err := ioutil.ReadFile(dir + f.Name())
			checkErr(err)
			fileContent := getWords(string(bytes))
			countWords(wordsMap, fileContent)
			wg.Done()
		}(f, &wordsMap, wg)
	}
	wg.Wait()

	return &wordsMap
}

func getWords(str string) []string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)

	reg, err := regexp.Compile("[^a-z\\-\n ]+")
	checkErr(err)
	str = reg.ReplaceAllString(str, " ")

	var oldLength = 1
	var length int
	for length < oldLength {
		oldLength = len(str)
		str = strings.Replace(str, "\n", " ", -1)
		str = strings.Replace(str, "  ", " ", -1)
		length = len(str)
	}

	return strings.Split(str, " ")
}

func countWords(wordsMap *sync.Map, words []string) {
	for _, word := range words {
		v, ok := wordsMap.Load(word)
		var count int
		if ok {
			count = v.(int)
		}
		count++
		wordsMap.Store(word, count)
	}
}

func rankByWordCount(wordsMap *sync.Map) PairList {
	pl := PairList{}
	wordsMap.Range(func(k, v interface{}) bool {
		pl = append(pl, Pair{k.(string), v.(int)})
		return true
	})
	sort.Sort(sort.Reverse(pl))

	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func writeCSV(pairList PairList) {

	file, err := os.Create(resultingFile)
	checkErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, pair := range pairList {
		value := []string{pair.Key, strconv.Itoa(pair.Value)}
		err := writer.Write(value)
		checkErr(err)
	}
}

func removeKnownWords(unknownWords, knownWords *sync.Map) {
	wg := &sync.WaitGroup{}
	unknownWords.Range(func(uw, _ interface{}) bool {
		wg.Add(1)
		go func(knownWords *sync.Map, wg *sync.WaitGroup) {
			knownWords.Range(func(kw, _ interface{}) bool {
				if kw.(string) == uw.(string) {
					unknownWords.Delete(kw)
				}
				return true
			})
			wg.Done()
		}(knownWords, wg)
		return true
	})
	wg.Wait()
}
