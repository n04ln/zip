package main

import (
	"archive/zip"
	"io"
	"net/http"
)

var (
	filePath = "webiner/"
)

type Attendee struct {
	Email     string `csv:"email"`
	FirstName string `csv:"first_name"`
	LastName  string `csv:"last_name"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=webinar.zip")

	files := map[string]io.ReadWriter{}
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	for i, file := range files {
		if err := addToZip(i, file, zipWriter); err != nil {
			panic(err)
		}
	}
}

func newCSVHeader() []string {
	return []string{
		"email",
		"first_name",
		"last_name",
	}
}

func (a *Attendee) ToStrings() []string {
	return []string{a.Email, a.FirstName, a.LastName}
}

func split(total []*Attendee, unit int) [][]*Attendee {
	result := [][]*Attendee{}
	size := len(total)
	for i := 0; i < size; i += unit {
		j := i + unit
		if j > size {
			j = size
		}
		result = append(result, total[i:j])
	}
	return result
}

func addToZip(name string, file io.ReadWriter, zipWriter *zip.Writer) error {
	w, err := zipWriter.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}
