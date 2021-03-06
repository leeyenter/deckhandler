package data

import (
	"encoding/csv"
	"errors"
	"os"
)

// LoadCSVFile parses a CSV file that contains
// the list of cards that should be loaded into the database.
// The first row should be the header.
// The first column should be the card code, which will be used
// as a primary key to identify the card.
func LoadCSVFile(file string) ([]Card, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("could not open file " + file + ": " + err.Error())
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("error parsing CSV file: " + err.Error())
	}

	headers := make(map[string]int) // look-up for CSV columns
	cards := make([]Card, 0)

	for idx, record := range records {
		if idx == 0 {
			// Headers
			for headerIdx, header := range record[1:] {
				headers[header] = headerIdx + 1
			}
		} else {
			// Check record length first
			if len(record) < len(headers)+1 {
				return nil, errors.New("inconsistent number of column headings")
			}

			// Save record
			var card Card
			card.ID = record[0]
			card.Values = make(map[string]string, len(headers))
			for header, headerIdx := range headers {
				card.Values[header] = record[headerIdx]
			}

			cards = append(cards, card)
		}

	}

	return cards, nil
}
