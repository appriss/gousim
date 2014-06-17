package reports

import (
	"text/template"
	"time"
	"fmt"
	"io"
	"encoding/json"
	"bytes"
)

var funcMap = template.FuncMap{"max":max, "min":min, "average":average, "total":total}

//The default text report template string
var text = "\nSimulation Results (Seconds):\n\n    Start Time: {{printf \"%32s\" (.StartTime)}}\n      End Time: {{printf \"%32s\" (.EndTime)}}\n\n{{range $i, $v := .Samples}}    Metric: {{printf \"%-30s\" $i}} Max: {{printf \"%-10s\" (max $v)}} Min: {{printf \"%-10s\" (min $v)}} Average: {{printf \"%-10s\" (average $v)}} Total: {{printf \"%-10s\" (total $v)}}  \n{{end}}\n"

func GetFuncMap() template.FuncMap {
	return funcMap
}

//Prints out the default text report template
func TextReport(output io.Writer, data []byte) error {
	tmpl, err := template.New("TEXT REPORT").Funcs(funcMap).Parse(text)
	if err != nil {
		return err
	}

	report, err := loadData(data)
	if err != nil {
		return err
	}

	err = tmpl.Execute(output, report)
	if err != nil {
		return err
	}

	return nil
}

//Prints out the default html report template
func HtmlReport() error {
	return nil
}

//Prints out a custom report template provided by the caller
func CustomReport(output io.Writer, result []byte, tmplstr string) error {
	return nil
}

type Report struct {
        StartTime string
        EndTime   string
        Samples   map[string][]int64
}

//Function made available in the template to print the highest sample value
func max(data []int64) string {
        var value int64
        for _, n := range data {
                if n > value {
                        value = n
                }
        }
        duration, _ := time.ParseDuration(fmt.Sprintf("%dns", value))
        return fmt.Sprintf("%.3f", duration.Seconds())
}

//Function made available in the template to print the lowest sample value
func min(data []int64) string {
        value := data[0]
        for _, n := range data {
                if n < value {
                        value = n
                }
        }
        duration, _ := time.ParseDuration(fmt.Sprintf("%dns", value))
        return fmt.Sprintf("%.3f", duration.Seconds())
}

//Function made available in the template to print the average across the sample values
func average(data []int64) string {
        var total int64
        for _, n := range data {
                total += n
        }
        duration, _ := time.ParseDuration(fmt.Sprintf("%dns", (total / int64(len(data)))))
        return fmt.Sprintf("%.3f", duration.Seconds())
}

//Function made available in the template to print the total (SUM) of the sample values
func total(data []int64) string {
        var total int64
        for _, n := range data {
                total += n
        }
        duration, _ := time.ParseDuration(fmt.Sprintf("%dns", total))
        return duration.String()
}

func loadData(data []byte) (Report, error) {
        var report Report
        err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&report)
        if err != nil {
               return report, err
        }
        return report, nil
}

