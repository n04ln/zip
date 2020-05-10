package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
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

	// Create user data for 10000 people.
	totalRecords := make([]*Attendee, 10000)
	for i := range totalRecords {
		email := fmt.Sprintf("%v@gmail.com", i+1)
		firstName := fmt.Sprintf("%v番目のユーザー", i+1)
		lastName := "でーす！"
		totalRecords[i] = &Attendee{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		}
	}

	// Divided into 100 user.
	unitRecords := split(totalRecords, 100)
	if err := os.Mkdir("webiner", 0777); err != nil {
		fmt.Printf("err: %v", err)
	}
	// defer os.RemoveAll("./webiner")

	files := map[string]io.ReadWriter{}
	// Write user information to a csv file for every 100 user.
	for i, records := range unitRecords {
		// file, err := os.Create(filePath + fmt.Sprintf("attendees%v.csv", i+1))
		// if err != nil {
		// 	fmt.Printf("err: %v", err)
		// }
		// defer file.Close()
		fileName := filePath + fmt.Sprintf("attendees%v.csv", i+1)
		files[fileName] = new(bytes.Buffer)

		writer := csv.NewWriter(files[fileName])
		defer writer.Flush()

		if err := writer.Write(newCSVHeader()); err != nil {
			fmt.Printf("err: %v", err)
		}
		for _, record := range records {
			writer.Write(record.ToStrings())
		}
	}

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// files, err := ioutil.ReadDir("./webiner")
	// if err != nil {
	// 	panic(err)
	// }
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
	// file, err := os.Open(fmt.Sprintf("./webiner/%v", name))
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// b, err := ioutil.ReadFile(fmt.Sprintf("./webiner/%v", name))
	// if err != nil {
	// 	panic(err)
	// }

	// n, err := io.Copy(os.Stdout, bytes.NewReader(b))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("name: %v, bytes: %v\n", file.Name(), n)

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
