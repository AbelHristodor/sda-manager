package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("SDA Manager")

	label := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		label,
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	))
	w.ShowAndRun()

	// res := []models.Verse{}
	// conn.Find(&res, "ID IN ?", search)
	// for _, verse := range res {
	// 	log.Println("Hymn Number: " + strconv.Itoa(int(verse.HymnID)) + " Text: " + verse.Text)
	// }

}


func get_current_dir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func open_pptx_file(file_path string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", file_path).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", file_path).Start()
	default:
		log.Fatal("Unsupported platform")
	}
}