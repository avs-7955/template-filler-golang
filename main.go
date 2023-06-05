package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello")

	// open CSV file
	fd, error := os.Open("./files/new.csv")
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println("Successfully opened the CSV file!")
	defer fd.Close()

	// read CSV file
	csvReader := io.Reader(fd)
	// converting CSV to map
	rows := CSVToMap(csvReader)

	// message from web
	sentence := "My name is {{.Name}}. I am {{.Age}} years old. I live in {{.Country}}.\n"
	// creating a template string
	Create := func(name, t string) *template.Template {
		return template.Must(template.New(name).Parse(t))
	}
	t2 := Create("t2", sentence)

	// for every entry of the csv file, converting into txt file.
	for _, row := range rows {
		t2.Execute(os.Stdout, row)
		f, err := os.Create("./files/" + row["Name"] + ".docx")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		t2.Execute(w, row)
		w.Flush()
	}
}

func CSVToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[strings.TrimSpace(header[i])] = strings.TrimSpace(record[i])
			}
			rows = append(rows, dict)
		}
	}
	return rows
}
