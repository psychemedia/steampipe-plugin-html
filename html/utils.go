package html

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHtmlPage(webPage string) (string, error) {

	resp, err := http.Get(webPage)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func createTables(name string, url string, path string) ([]string, []string) {
	data, err := getHtmlPage(url)
	if err != nil {
		panic(err)
	}

	tables := []string{}
	columns := []string{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		panic(err)
	}

	doc.Find("table").Each(func(iTbl int, tablehtml *goquery.Selection) {
		rows := ""
		cols := ""
		tableName := fmt.Sprintf("%v_%v\n", name, iTbl)
		tables = append(tables, tableName)
		tablehtml.Find("tr").Each(func(iTr int, rowhtml *goquery.Selection) {
			row := ""
			rowhtml.Find("th").Each(func(iTh int, tableheading *goquery.Selection) {
				colName := normalize(tableheading.Text())
				if colName != "" {
					row += fmt.Sprintf("\"%s\",", colName)
					cols += fmt.Sprintf("\"%s\",", colName)
				}
			})
			rowhtml.Find("td").Each(func(iTd int, tablecell *goquery.Selection) {
				colVal := tablecell.Text()
				if colVal != "" {
					row += fmt.Sprintf("\"%s\",", colVal)
				}
			})
			rows += fmt.Sprintf("%s\n", row[:len(row)-1])
		})
		columns = append(columns, cols)
		writeRowsToCsv(path, fmt.Sprintf("%s_%d.csv", name, iTbl), rows)
	})

	return tables, columns
}

func writeRowsToCsv(path string, name string, rows string) {
	f, _ := os.Create(fmt.Sprintf("%s/%s", path, name))
	w := bufio.NewWriter(f)
	w.WriteString(rows)
	w.Flush()
}

func forceASCII(s string) string {
	rs := make([]rune, 0, len(s))
	for _, r := range s {
	  if r <= 127 {
		rs = append(rs, r)
	  }
	}
	return string(rs)
  }

func normalize(s string) string {
	rs := forceASCII(s)
	return rs
}