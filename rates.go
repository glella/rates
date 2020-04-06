// Get exchange rates YTD & PY and save to CSV file
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {

	// Get the YTD dates to download current year and PY - Last business day
	dates := datesToDownload()

	// Create table and add headers
	ratesTable := [][]string{
		[]string{"Rates", "AED", "ARS", "AUD", "BRL", "CAD", "CLP", "CNY", "COP", "EUR", "GBP", "IDR", "INR", "IQD", "JOD", "JPY", "KES", "KRW", "LKR", "MXN", "MYR", "NGN", "NZD", "RUB", "SAR", "SGD", "ZAR"},
	}

	// Download the rates for each day and append to results table
	for _, date := range dates {
		rates := downloadData(date)
		ratesTable = append(ratesTable, rates)
	}

	// Transpose columns to rows
	finalData := transpose(ratesTable)

	// Save data to CSV file
	// Get current executable path to use same dir
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Name the file with today's date and create filepath
	now := time.Now()
	strToday := now.Format("2006-01-02")
	filename := "ytd_rates_downloaded_" + strToday + ".csv"
	exPath := filepath.Dir(ex)
	filepath := path.Join(exPath, filename)

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write everything and we are done
	w := csv.NewWriter(file)
	w.WriteAll(finalData) // calls Flush internally
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

}

func downloadData(aDate string) []string {
	url := "https://openexchangerates.org/api/historical/" + aDate + ".json"
	// You need your own key for an app_id. Get it at openexchangerates.org
	// Without a key it wont work.
	url += "?app_id=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Error fetch: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var data historicalJSON
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		os.Exit(1)
	}

	return convertStructToArray(data, aDate)
}

func convertStructToArray(data historicalJSON, date string) []string {
	var result []string

	result = append(result, date)
	result = append(result, fmt.Sprintf("%.4f", data.Rates.AED))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.ARS))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.AUD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.BRL))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.CAD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.CLP))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.CNY))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.COP))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.EUR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.GBP))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.IDR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.INR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.IQD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.JOD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.JPY))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.KES))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.KRW))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.LKR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.MXN))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.MYR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.NGN))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.NZD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.RUB))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.SAR))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.SGD))
	result = append(result, fmt.Sprintf("%.4f", data.Rates.ZAR))

	return result

}

func calculateLastBusinessDay(year int, month int) string {
	now := time.Now()
	currentLocation := now.Location()
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	if lastOfMonth.Weekday() == 0 {
		lastOfMonth = lastOfMonth.AddDate(0, 0, -2)
	}
	if lastOfMonth.Weekday() == 6 {
		lastOfMonth = lastOfMonth.AddDate(0, 0, -1)
	}

	return lastOfMonth.Format("2006-01-02")
}

func elapsedMonths() int {
	now := time.Now()
	currentLocation := now.Location()
	since := time.Date(now.Year(), time.Month(1), 1, 0, 0, 0, 0, currentLocation)

	months := 0
	month := since.Month()
	for since.Before(now) {
		since = since.Add(time.Hour * 24)
		nextMonth := since.Month()
		if nextMonth != month {
			months++
		}
		month = nextMonth
	}

	return months
}

func datesToDownload() []string {
	now := time.Now()
	currentYear := now.Year()
	previousYear := currentYear - 1
	elapsedMonthsYTD := elapsedMonths()

	var result []string
	for i := 1; i < elapsedMonthsYTD+1; i++ {
		date := calculateLastBusinessDay(currentYear, i)
		result = append(result, date)
	}
	for i := 1; i < elapsedMonthsYTD+1; i++ {
		date := calculateLastBusinessDay(previousYear, i)
		result = append(result, date)
	}

	return result
}

func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

type historicalJSON struct {
	Disclaimer string `json:"disclaimer"`
	License    string `json:"license"`
	Timestamp  int    `json:"timestamp"`
	Base       string `json:"base"`
	Rates      struct {
		AED float64 `json:"AED"`
		AFN float64 `json:"AFN"`
		ALL float64 `json:"ALL"`
		AMD float64 `json:"AMD"`
		ANG float64 `json:"ANG"`
		AOA float64 `json:"AOA"`
		ARS float64 `json:"ARS"`
		AUD float64 `json:"AUD"`
		AWG float64 `json:"AWG"`
		AZN float64 `json:"AZN"`
		BAM float64 `json:"BAM"`
		BBD float64 `json:"BBD"`
		BDT float64 `json:"BDT"`
		BGN float64 `json:"BGN"`
		BHD float64 `json:"BHD"`
		BIF float64 `json:"BIF"`
		BMD float64 `json:"BMD"`
		BND float64 `json:"BND"`
		BOB float64 `json:"BOB"`
		BRL float64 `json:"BRL"`
		BSD float64 `json:"BSD"`
		BTC float64 `json:"BTC"`
		BTN float64 `json:"BTN"`
		BWP float64 `json:"BWP"`
		BYN float64 `json:"BYN"`
		BZD float64 `json:"BZD"`
		CAD float64 `json:"CAD"`
		CDF float64 `json:"CDF"`
		CHF float64 `json:"CHF"`
		CLF float64 `json:"CLF"`
		CLP float64 `json:"CLP"`
		CNH float64 `json:"CNH"`
		CNY float64 `json:"CNY"`
		COP float64 `json:"COP"`
		CRC float64 `json:"CRC"`
		CUC float64 `json:"CUC"`
		CUP float64 `json:"CUP"`
		CVE float64 `json:"CVE"`
		CZK float64 `json:"CZK"`
		DJF float64 `json:"DJF"`
		DKK float64 `json:"DKK"`
		DOP float64 `json:"DOP"`
		DZD float64 `json:"DZD"`
		EGP float64 `json:"EGP"`
		ERN float64 `json:"ERN"`
		ETB float64 `json:"ETB"`
		EUR float64 `json:"EUR"`
		FJD float64 `json:"FJD"`
		FKP float64 `json:"FKP"`
		GBP float64 `json:"GBP"`
		GEL float64 `json:"GEL"`
		GGP float64 `json:"GGP"`
		GHS float64 `json:"GHS"`
		GIP float64 `json:"GIP"`
		GMD float64 `json:"GMD"`
		GNF float64 `json:"GNF"`
		GTQ float64 `json:"GTQ"`
		GYD float64 `json:"GYD"`
		HKD float64 `json:"HKD"`
		HNL float64 `json:"HNL"`
		HRK float64 `json:"HRK"`
		HTG float64 `json:"HTG"`
		HUF float64 `json:"HUF"`
		IDR float64 `json:"IDR"`
		ILS float64 `json:"ILS"`
		IMP float64 `json:"IMP"`
		INR float64 `json:"INR"`
		IQD float64 `json:"IQD"`
		IRR float64 `json:"IRR"`
		ISK float64 `json:"ISK"`
		JEP float64 `json:"JEP"`
		JMD float64 `json:"JMD"`
		JOD float64 `json:"JOD"`
		JPY float64 `json:"JPY"`
		KES float64 `json:"KES"`
		KGS float64 `json:"KGS"`
		KHR float64 `json:"KHR"`
		KMF float64 `json:"KMF"`
		KPW float64 `json:"KPW"`
		KRW float64 `json:"KRW"`
		KWD float64 `json:"KWD"`
		KYD float64 `json:"KYD"`
		KZT float64 `json:"KZT"`
		LAK float64 `json:"LAK"`
		LBP float64 `json:"LBP"`
		LKR float64 `json:"LKR"`
		LRD float64 `json:"LRD"`
		LSL float64 `json:"LSL"`
		LYD float64 `json:"LYD"`
		MAD float64 `json:"MAD"`
		MDL float64 `json:"MDL"`
		MGA float64 `json:"MGA"`
		MKD float64 `json:"MKD"`
		MMK float64 `json:"MMK"`
		MNT float64 `json:"MNT"`
		MOP float64 `json:"MOP"`
		MRO float64 `json:"MRO"`
		MRU float64 `json:"MRU"`
		MUR float64 `json:"MUR"`
		MVR float64 `json:"MVR"`
		MWK float64 `json:"MWK"`
		MXN float64 `json:"MXN"`
		MYR float64 `json:"MYR"`
		MZN float64 `json:"MZN"`
		NAD float64 `json:"NAD"`
		NGN float64 `json:"NGN"`
		NIO float64 `json:"NIO"`
		NOK float64 `json:"NOK"`
		NPR float64 `json:"NPR"`
		NZD float64 `json:"NZD"`
		OMR float64 `json:"OMR"`
		PAB float64 `json:"PAB"`
		PEN float64 `json:"PEN"`
		PGK float64 `json:"PGK"`
		PHP float64 `json:"PHP"`
		PKR float64 `json:"PKR"`
		PLN float64 `json:"PLN"`
		PYG float64 `json:"PYG"`
		QAR float64 `json:"QAR"`
		RON float64 `json:"RON"`
		RSD float64 `json:"RSD"`
		RUB float64 `json:"RUB"`
		RWF float64 `json:"RWF"`
		SAR float64 `json:"SAR"`
		SBD float64 `json:"SBD"`
		SCR float64 `json:"SCR"`
		SDG float64 `json:"SDG"`
		SEK float64 `json:"SEK"`
		SGD float64 `json:"SGD"`
		SHP float64 `json:"SHP"`
		SLL float64 `json:"SLL"`
		SOS float64 `json:"SOS"`
		SRD float64 `json:"SRD"`
		SSP float64 `json:"SSP"`
		STD float64 `json:"STD"`
		STN float64 `json:"STN"`
		SVC float64 `json:"SVC"`
		SYP float64 `json:"SYP"`
		SZL float64 `json:"SZL"`
		THB float64 `json:"THB"`
		TJS float64 `json:"TJS"`
		TMT float64 `json:"TMT"`
		TND float64 `json:"TND"`
		TOP float64 `json:"TOP"`
		TRY float64 `json:"TRY"`
		TTD float64 `json:"TTD"`
		TWD float64 `json:"TWD"`
		TZS float64 `json:"TZS"`
		UAH float64 `json:"UAH"`
		UGX float64 `json:"UGX"`
		USD float64 `json:"USD"`
		UYU float64 `json:"UYU"`
		UZS float64 `json:"UZS"`
		VEF float64 `json:"VEF"`
		VES float64 `json:"VES"`
		VND float64 `json:"VND"`
		VUV float64 `json:"VUV"`
		WST float64 `json:"WST"`
		XAF float64 `json:"XAF"`
		XAG float64 `json:"XAG"`
		XAU float64 `json:"XAU"`
		XCD float64 `json:"XCD"`
		XDR float64 `json:"XDR"`
		XOF float64 `json:"XOF"`
		XPD float64 `json:"XPD"`
		XPF float64 `json:"XPF"`
		XPT float64 `json:"XPT"`
		YER float64 `json:"YER"`
		ZAR float64 `json:"ZAR"`
		ZMW float64 `json:"ZMW"`
		ZWL float64 `json:"ZWL"`
	} `json:"rates"`
}
