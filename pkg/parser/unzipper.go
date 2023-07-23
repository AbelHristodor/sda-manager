package parser

import (
	"archive/zip"
	"log"
	"sort"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Slide struct {
	Number  int
	Text	string
}

func GetFileNames(filenames []*zip.File) []string {
	var names []string
	for _, f := range filenames {
		names = append(names, f.FileHeader.Name)
	}
	return names
}

func ExtractZipFile(zip_file_path string) []Slide {
	
	rd, err := zip.OpenReader(zip_file_path)
	if err != nil {
		log.Fatal("Error opening zip file: ", err, " ", zip_file_path)
	}

	defer rd.Close()

	// Find the slides dir in the zip
	const slidesDir = "ppt/slides/"
	var slides []*zip.File
	for _, f := range rd.File {
		if f.FileHeader.Name[:len(slidesDir)] == slidesDir {
			ext := strings.Split(f.FileHeader.Name, ".")
			if ext[len(ext)-1] == "xml" {
				slides = append(slides, f)
			}
		}
	}
	// Sort the slides array
	sort.Slice(slides, func(i, j int) bool {
		return slides[i].FileHeader.Name < slides[j].FileHeader.Name
	})

	var slidesContents []Slide = make([]Slide, len(slides))
	
	// Extract the slides
	for s, slideFile := range slides {
		slideReader, err := slideFile.Open()
		if err != nil {
			log.Fatal("Error opening slide file: ", err)
		}
		defer slideReader.Close()

		doc := etree.NewDocument()
		
		if _, err := doc.ReadFrom(slideReader); err != nil {
			log.Fatal("Error reading slide file: ", err)
		}

		tag := doc.FindElements("//a:t")
		text := ""
		for _, t := range tag {
			text += RomanianToASCII(t.Text()) + " "
		}

		slidesContents = append(slidesContents, Slide{Number: s + 1, Text: text})
	}

	return slidesContents
}

func RomanianToASCII(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isDiacritic), norm.NFC)
	result, _, _ := transform.String(t, input)
	result = strings.ReplaceAll(result, "ă", "a")
	result = strings.ReplaceAll(result, "â", "a")
	result = strings.ReplaceAll(result, "î", "i")
	result = strings.ReplaceAll(result, "ș", "s")
	result = strings.ReplaceAll(result, "ț", "t")
	return result
}

func isDiacritic(c rune) bool {
	return unicode.Is(unicode.Mn, c) // Mn: nonspacing marks
}
