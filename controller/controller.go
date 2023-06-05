package controller

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Message string `json:"message"`
}

// ROUTING METHODS
func ConvertMessagetoTXT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST") // will only accept post requests
	// if the data received is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
		return
	}

	var received_req Message
	err := json.NewDecoder(r.Body).Decode(&received_req)
	CheckError(err)
	fmt.Println(received_req)
	InsertValuesIntoTxt(received_req.Message)

	json.NewEncoder(w).Encode("Written all emails in txt files.")
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Working")
}

// CONVERTING MESSAGE TO TXT FILE
// open CSV file
func InsertValuesIntoTxt(sentence string) {
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
	// sentence := "My name is {{.Name}}. I am {{.Age}} years old. I live in {{.Country}}.\n"
	// creating a template string
	Create := func(name, t string) *template.Template {
		return template.Must(template.New(name).Parse(t))
	}
	t2 := Create("t2", sentence)

	// for every entry of the csv file, converting into txt file.
	for _, row := range rows {
		// t2.Execute(os.Stdout, row)
		f, err := os.Create("./files/" + row["Name"] + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		t2.Execute(w, row)
		w.Flush()
	}
}

// CONVERTING CSV TO MAPS
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

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
