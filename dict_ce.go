package wasmdict

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"io"
	"log"
	"strings"
)

//go:embed data/cedict_1_0_ts_utf-8_mdbg.zip
var ceDictZipData []byte

type DictEntryCE struct {
	Traditional string
	Simplified  string
	Pinyin      string
	English     string
}

func (d *DictEntryCE) Map() map[string]interface{} {
	if d == nil {
		return nil
	}
	return map[string]interface{}{
		"traditional": d.Traditional,
		"simplified":  d.Simplified,
		"pinyin":      d.Pinyin,
		"english":     d.English,
	}
}

var ceItems []*DictEntryCE

func parseLine(line string) *DictEntryCE {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil
	}
	if strings.HasPrefix(line, "#") {
		return nil
	}

	parts := strings.Split(line, "/")
	if len(parts) <= 1 {
		return nil
	}
	english := parts[1]
	charAndPinyin := strings.Split(parts[0], "[")
	characters := strings.Fields(charAndPinyin[0])
	if len(characters) < 2 {
		return nil
	}
	traditional := characters[0]
	simplified := characters[1]
	pinyin := strings.TrimRight(charAndPinyin[1], "]")
	return &DictEntryCE{
		Traditional: traditional,
		Simplified:  simplified,
		Pinyin:      pinyin,
		English:     english,
	}
}

func removeSurnames(entries []*DictEntryCE) []*DictEntryCE {
	var result []*DictEntryCE
	for i := 0; i < len(entries); i++ {
		if strings.Contains(entries[i].English, "surname ") {
			if i+1 < len(entries) && entries[i].Traditional == entries[i+1].Traditional {
				continue
			}
		}
		result = append(result, entries[i])
	}
	return result
}

// PreLoadCeDict loads the CEDICT data from the embedded ZIP file and parses it into dictionary entries.
// It returns two slices of dictionary entries: one sorted by Traditional Chinese and the other by Simplified Chinese.
func PreLoadCeDict() []*DictEntryCE {
	if len(ceItems) > 0 {
		return ceItems
	}
	filesMap, err := extractZipBytes(ceDictZipData)
	if err != nil {
		log.Println("Error extracting ZIP data:", err)
		return nil
	}
	ceData, ok := filesMap["cedict_ts.u8"]
	if !ok {
		log.Println("Error extracting cedict_ts.u8")
		return nil
	}
	var result []*DictEntryCE
	lines := strings.Split(string(ceData), "\n")
	for _, line := range lines {
		entry := parseLine(line)
		if entry != nil {
			result = append(result, entry)
		}
	}
	ceItems = removeSurnames(result)
	//free memory
	ceDictZipData = nil
	return ceItems
}

func extractZipBytes(zipData []byte) (map[string][]byte, error) {
	filesContent := make(map[string][]byte)

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	// Iterate through each file in the zip archive
	for _, zipFile := range zipReader.File {
		// Open the file
		f, err := zipFile.Open()
		if err != nil {
			log.Println("Error opening file:", err)
			continue
		}
		defer f.Close()

		// Read the file's contents
		content, err := io.ReadAll(f)
		if err != nil {
			log.Println("Error reading file contents:", err)
			continue
		}

		// Store the contents in the map
		filesContent[zipFile.Name] = content
	}

	return filesContent, nil
}

// CeQueryLike searches for dictionary entries that start with the specified text.
// It supports both Simplified Chinese (isCnZh = true) and Traditional Chinese (isCnZh = false).
// The search is limited to a specified number of results (count).
// If the count is reached, the search stops and returns the found entries up to that point.
// Parameters:
// - text: The text to search for at the beginning of dictionary entries.
// - isCnZh: A boolean indicating whether to search in Simplified Chinese (true) or Traditional Chinese (false).
// - count: The maximum number of entries to return.
// Returns: A slice of pointers to DictEntryCE structs that match the search criteria.
func CeQueryLike(text string, isCnZh bool, count int) (result []*DictEntryCE) {
	text = strings.TrimSpace(text)
	items := PreLoadCeDict()
	if isCnZh {
		for _, w := range items {
			if strings.HasPrefix(w.Simplified, text) {
				result = append(result, w)
			}
			if len(result) >= count {
				return result
			}
		}
	} else {
		for _, w := range items {
			if strings.HasPrefix(w.Traditional, text) {
				result = append(result, w)
			}
			if len(result) >= count {
				return result
			}
		}

	}
	return nil
}

// CeLookUp searches for a single dictionary entry that exactly matches the specified text.
// It supports both Simplified Chinese (isCnZh = true) and Traditional Chinese (isCnZh = false).
// Parameters:
// - text: The text to search for in the dictionary entries.
// - isCnZh: A boolean indicating whether to search in Simplified Chinese (true) or Traditional Chinese (false).
// Returns: A pointer to a DictEntryCE struct if a match is found, or nil if no match is found.
func CeLookUp(text string, isCnZh bool) *DictEntryCE {
	text = strings.TrimSpace(text)
	items := PreLoadCeDict()
	if isCnZh {
		for _, w := range items {
			if strings.Compare(w.Simplified, text) == 0 {
				return w
			}
		}
	} else {
		for _, w := range items {
			if strings.Compare(w.Traditional, text) == 0 {
				return w
			}
		}
	}
	return nil
}
