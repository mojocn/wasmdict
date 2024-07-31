package wasmdict

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/csv"
	"io"
	"log"
	"strings"
)

//go:embed data/ecdict.csv
var ecdictCsv []byte

//go:embed data/lemma.en.txt
var lemmaEnTxt []byte

type DictEntryEC struct {
	Word        string
	Phonetic    string
	Definition  string
	Translation string
	Pos         string
	Collins     string
	Oxford      string
	Tag         string
	Bnc         string
	Frq         string
	Exchange    string
	Detail      string
	Audio       string
}

func (d *DictEntryEC) Map() map[string]interface{} {
	if d == nil {
		return nil
	}
	return map[string]interface{}{
		"word":        d.Word,
		"phonetic":    d.Phonetic,
		"definition":  d.Definition,
		"translation": d.Translation,
		"pos":         d.Pos,
		"collins":     d.Collins,
		"oxford":      d.Oxford,
		"tag":         d.Tag,
		"bnc":         d.Bnc,
		"frq":         d.Frq,
		"exchange":    d.Exchange,
		"detail":      d.Detail,
		"audio":       d.Audio,
	}
}

var dictMapSingleton = map[string]DictEntryEC{}
var lemmaMapSingleton = map[string]string{}

// PreLoadEcDict ensures that the dictionary and lemma maps are loaded into memory.
// It checks if the `dictMapSingleton` and `lemmaMapSingleton` are empty.
// If they are empty, it calls `parseDict` to load the dictionary data and `parseLemma` to load the lemma data.
func PreLoadEcDict() {
	if len(dictMapSingleton) == 0 {
		dictMapSingleton = parseDict()
	}
	if len(lemmaMapSingleton) == 0 {
		lemmaMapSingleton = parseLemma()
	}
	// Free up memory by setting the byte slices to nil.
	ecdictCsv = nil
	lemmaEnTxt = nil
}

func parseDict() map[string]DictEntryEC {
	dictMap := map[string]DictEntryEC{}

	r := csv.NewReader(bytes.NewReader(ecdictCsv))
	rowElementCount := 13
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if rowElementCount != len(record) {
			continue
		}
		word := strings.TrimSpace(record[0])
		if word == "" {
			continue
		}
		dictItem := DictEntryEC{
			Word:        word,
			Phonetic:    record[1],
			Definition:  removeBr(record[2]),
			Translation: removeBr(record[3]),
			Pos:         record[4],
			Collins:     record[5],
			Oxford:      record[6],
			Tag:         record[7],
			Bnc:         record[8],
			Frq:         record[9],
			Exchange:    record[10],
			Detail:      record[11],
			Audio:       record[12],
		}
		if dictItem.Word == "word" && dictItem.Phonetic == "phonetic" {
			//skip csv header
			continue
		}
		dictMap[word] = dictItem
		dictMap[strings.ToLower(word)] = dictItem
	}
	return dictMap
}
func removeBr(w string) string {
	return strings.ReplaceAll(w, "\\n", "\n")
}
func parseLemma() map[string]string {
	lemmaMap := map[string]string{}

	scanner := bufio.NewScanner(bytes.NewReader(lemmaEnTxt))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ";") {
			continue
		}
		parts := strings.Split(line, " -> ")
		if len(parts) != 2 {
			continue
		}
		rParts, lParts := strings.Split(parts[0], "/"), strings.Split(parts[1], ",")
		originalWord := ""
		if len(rParts) > 0 {
			originalWord = strings.TrimSpace(rParts[0])
		}
		if originalWord == "" {
			continue
		}
		for _, lemma := range lParts {
			lemma = strings.TrimSpace(lemma)
			if lemma == "" {
				continue
			}
			lemmaMap[lemma] = originalWord
		}
	}
	return lemmaMap
}

// EcLookUp searches for a given word in the dictionary and its lemma form.
// It first trims any leading or trailing spaces from the input word.
// If the word exists in the lemma map (lemmaMapSingleton), it retrieves the base form of the word.
// Otherwise, it proceeds with the original word.
// Then, it attempts to find the base word in the dictionary map (dictMapSingleton).
// If found, it returns a pointer to the corresponding DictEntryEC.
// If the word or its base form is not found in the dictionary, it returns nil.
//
// Parameters:
// - word: The word to look up in the dictionary.
//
// Returns:
// - *DictEntryEC: A pointer to the dictionary item if found; otherwise, nil.
func EcLookUp(word string) *DictEntryEC {
	PreLoadEcDict()                         // Ensure the dictionary is loaded before searching.
	word = strings.TrimSpace(word)          // Trim spaces from the input word.
	baseWord, ok := lemmaMapSingleton[word] // Check if the word has a base form in the lemma map.
	if !ok {
		baseWord = word // Use the original word if no base form is found.
	}
	dictItem, ok := dictMapSingleton[baseWord] // Look up the base word in the dictionary map.
	if !ok {
		return nil // Return nil if the word is not found in the dictionary.
	}
	return &dictItem // Return a pointer to the found dictionary item.
}

// EcQueryLike searches for words in the dictionary that start with the given prefix.
// It ensures the dictionary is loaded into memory before performing the search.
// It iterates through all words in the dictionary and collects those that start with the specified prefix.
// The search stops once the specified number of matches (cnt) is reached.
//
// Parameters:
// - w: The prefix to search for in the dictionary.
// - cnt: The maximum number of matching words to return.
//
// Returns:
// - []string: A slice of words that start with the given prefix, up to the specified count.
func EcQueryLike(w string, cnt int) []string {
	PreLoadEcDict()
	var result []string
	for word, _ := range dictMapSingleton {
		if strings.HasPrefix(word, w) {
			result = append(result, word)
		}
		if len(result) >= cnt {
			break
		}
	}
	return result
}
