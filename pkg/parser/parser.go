package parser

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sda-manager/pkg/db/models"
	"strings"

	"gorm.io/gorm"
)

const PPTX = "pptx"
const PPT = "ppt"

func ParseFolder(path string, db *gorm.DB) {
	// Check if path is a directory
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		log.Fatal("Path is not a directory")
	}

	// Get all files in directory recursively
	files := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error walking through directory: ", err)
	}

	// Parse all pptx files
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal("Error opening file: ", err)
		}

		if get_extension(f) == PPTX || get_extension(f) == PPT {
			ParsePPTXToDB(file, db)
		}
	}

}

func ParsePPTXToDB(path string, db *gorm.DB) bool {

	// Check if file exists
	fr, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	
	// Check if file is pptx file
	if get_extension(fr) != PPTX && get_extension(fr) != PPT {
		log.Fatal("File is not a pptx file: ", path)
		return false
	}

	// Extract the text from file
	slides := ExtractZipFile(path)

	// Create the hymn and verses dict
	hymn := models.Hymn{}
	verses := []models.Verse{}
	
	for _, slide := range slides {
		if slide.Text != "" {
			text := slide.Text
			
			if slide.Number == 1 {
				hymn.Title = get_title(text)
				hymn.Number = getHymnNumberFromFileName(fr.Name())
			} else {
				verses = append(verses, models.Verse{Text: text})
			}
		}
	}

	hymn.Verses = verses
	hymn.Path = path
	hymn.Category = models.Standard
	if id := db.Create(&hymn); id.Error != nil {
		log.Fatal("Error creating hymn: ", id.Error)
		return false
	}
	return true
}

func get_title(text string) string {
	split := strings.Split(text, "Imnul")
	return removeDuplicateSpaces(split[0])
}

func get_extension(file *os.File) string {
	ext := strings.Split(file.Name(), ".")
	return ext[len(ext)-1]
}

func removeDuplicateSpaces(input string) string {
    // Create a regular expression pattern to match one or more consecutive spaces
    pattern := regexp.MustCompile(`\s+`)
    
    // Replace all occurrences of consecutive spaces with a single space
    cleanedText := pattern.ReplaceAllString(input, " ")
    
    // Trim any leading or trailing spaces
    cleanedText = strings.TrimSpace(cleanedText)
    
    return cleanedText
}

func getHymnNumberFromFileName(filename string) string {
	// Get the number from the filename
	re := regexp.MustCompile("\\/(\\d+)([A-z]?)\\.")
	numberString := re.FindString(filename)
	if numberString == "" {
		log.Fatal("Error getting number from filename: ", filename)
	}
	return numberString[1:len(numberString)-1]
}