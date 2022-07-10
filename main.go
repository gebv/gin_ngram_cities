package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/kelindar/bitmap"
)

var allAliases = map[string][]string{
	"st": {"saint", "street"},
}

func main() {
	records := readCsvFile_getCities("_cities_en.csv")
	log.Println("total cities:", len(records))
	prepareIndex(records)
	log.Println("total records in GIN:", len(gramToID))

	cases := []string{
		"Saint-Petersburg",
		"St Petersburg",
	}
	for _, in := range cases {
		log.Println()
		log.Println("lookup:", in)
		res := lookup(normilizeAndSanitize(in))
		log.Printf("result: %#v\n", res)
	}

}

var seqID int64
var gramToID = map[string][]int64{}
var idToRecord = map[int64]string{}

func lookup(in string) []string {
	in = normilizeAndSanitize(in)
	words := strings.Fields(in)

	res := bitmap.Bitmap{}
	andOperands := []bitmap.Bitmap{}

	// WHERE word1 AND word2 ... AND wordN
	// в случае алиасов
	// WHERE word1 AND (alias1 OR alias2 OR ...) ... AND wordN
	for _, word := range words {
		if len(allAliases[word]) == 0 {
			for _, gram := range toNgram(word) {
				gramBitmap := bitmap.Bitmap{}
				for _, id := range gramToID[gram] {
					gramBitmap.Set(uint32(id))
					res.Set(uint32(id))
				}
				andOperands = append(andOperands, gramBitmap)
			}
		} else {
			orList := bitmap.Bitmap{}
			for _, alias := range allAliases[word] {
				for _, aliasGram := range toNgram(alias) {
					or := bitmap.Bitmap{}
					for _, id := range gramToID[aliasGram] {
						or.Set(uint32(id))
					}
					orList.Or(or)
				}
				andOperands = append(andOperands, orList)
			}
		}

	}

	// после каждого пересечения остается все меньше и меньше подходящих записей
	for _, and := range andOperands {
		res.And(and)
		// TODO: можно на этом этапе регулировать когда прекращать поиск
	}

	log.Println("matched count:", res.Count())
	// итоговый результат

	matched := []string{}
	res.Range(func(x uint32) {
		matched = append(matched, idToRecord[int64(x)])
	})

	return matched
}

func prepareIndex(records []string) {
	for _, record := range records {
		id := nextID()
		idToRecord[id] = record

		words := strings.Fields(record)
		for _, word := range words {
			if len(word) <= 2 {
				// пропускаем слова короче 2-х чимволов
				// по хорошму удалять надо все артиклы, предлоги и прочее
				continue
			}
			normal := normilizeAndSanitize(word)
			for _, gram := range toNgram(normal) {
				gramToID[gram] = append(gramToID[gram], id)
			}
		}
	}
}

func toNgram(in string) []string {
	const size = 3
	res := []string{}
	for i := 0; i < len(in); i++ {
		if i+size > len(in) {
			break
		}
		res = append(res, in[i:i+size])
	}
	return res
}

func normilizeAndSanitize(in string) string {
	in = strings.TrimSpace(in)
	in = strings.ToLower(in)
	return sanitize(in)
}

func nextID() int64 {
	defer func() {
		seqID++
	}()
	return seqID
}

func readCsvFile_getCities(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to read file %q: %w\n", filePath, err)
	}
	defer f.Close()

	res := []string{}
	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: handle errors
			continue
		}
		res = append(res, record[0])
	}

	return res
}

func sanitize(s string) string {
	return nonAlphanumericRegex.ReplaceAllString(s, " ")
}

var nonAlphanumericRegex = regexp.MustCompile("[^A-Za-z0-9]+")

// if err != nil {
// 	log.Fatal(err)
// }
// newStr := reg.ReplaceAllString("#Golang#Python$Php&Kotlin@@", "-")
