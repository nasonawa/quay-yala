package report

import (
	"html/template"

	"github.com/nasonawa/quay-yala/pkg"
)

var Textreport = `
* Number of requests: {{ .TotalRequests }}
* First: {{ .FirstRecord.Ip }} - - [{{ .FirstRecord.Date.Format "02/Jan/2006:15:04:05 -0700" }}] "HEAD {{ .FirstRecord.Url }} HTTP/1.0" {{ .FirstRecord.Code }} - "-" "-" UserAgent: "{{ .FirstRecord.UserAgent }}"
* Last: {{ .LastRecord.Ip }} - - [{{ .LastRecord.Date.Format "02/Jan/2006:15:04:05 -0700" }}] "HEAD {{ .LastRecord.Url }} HTTP/1.0" {{ .LastRecord.Code }} - "-" "-" UserAgent: "{{ .LastRecord.UserAgent }}"
* Number of 200s: {{ .Count200 }} - {{ printf "%.2f" .Percent200 }}%
* Number of 3XXs: {{ .Count3xx }} - {{ printf "%.2f" .Percent3xx }}%
* Number of 4XXs: {{ .Count4xx }} - {{ printf "%.2f" .Percent4xx }}%
* Number of 5XXs: {{ .Count5xx }} - {{ printf "%.2f" .Percent5xx }}%
* Highest completed request count is {{ .PeakRequestCount }} at {{ .PeakRequestTime }}
* Highest number of 4XX responses is {{ .Peak4xxCount }} at {{ .Peak4xxTime }}
* Highest number of 5XX responses is {{ .Peak5xxCount }} at {{ .Peak5xxTime }}
`

type StatusStats struct {
	Success     int
	Warning     int
	Error       int
	Redirection int
	Total       int
}

func CreateStatusStats(ar []pkg.AccessRecord) (ss StatusStats) {

	ss.Total = len(ar)

	for _, a := range ar {

		switch {

		case a.Code >= 200 && a.Code <= 299:
			ss.Success += 1
		case a.Code >= 300 && a.Code <= 399:
			ss.Redirection += 1
		case a.Code >= 400 && a.Code <= 499:
			ss.Warning += 1
		case a.Code >= 500 && a.Code <= 999:
			ss.Error += 1
		}
	}

	return
}

type HtmlTemplateData struct {
	Records  []pkg.AccessRecord
	Stats    StatusStats
	JSONData template.JS
}

type TextReportAnalysis struct {
	TotalRequests int
	FirstRecord   pkg.AccessRecord
	LastRecord    pkg.AccessRecord

	Count200   int
	Percent200 float64
	Count3xx   int
	Percent3xx float64
	Count4xx   int
	Percent4xx float64
	Count5xx   int
	Percent5xx float64

	PeakRequestCount int
	PeakRequestTime  string

	Peak4xxCount int
	Peak4xxTime  string

	Peak5xxCount int
	Peak5xxTime  string
}

func AccessLogTextReportAnalysis(records []pkg.AccessRecord) (at TextReportAnalysis) {

	totalfloat := float64(len(records))
	statusStats := CreateStatusStats(records)

	at.TotalRequests = statusStats.Total

	at.FirstRecord = records[0]
	at.LastRecord = records[len(records)-1]

	at.Count200 = statusStats.Success
	at.Percent200 = (float64(statusStats.Success) / totalfloat) * 100

	at.Count3xx = statusStats.Redirection
	at.Percent3xx = (float64(statusStats.Redirection) / totalfloat) * 100

	at.Count4xx = statusStats.Warning
	at.Percent4xx = (float64(statusStats.Warning) / totalfloat) * 100

	at.Count5xx = statusStats.Error
	at.Percent5xx = (float64(statusStats.Error) / totalfloat) * 100

	requestpermin := make(map[string]int)
	err4xxpermin := make(map[string]int)
	err5xxpermin := make(map[string]int)

	for _, v := range records {

		timekey := v.Date.Format("02:Jan:2006:15:04")

		requestpermin[timekey]++
		if v.Code >= 400 && v.Code < 500 {
			err4xxpermin[timekey]++
		}

		if v.Code >= 500 && v.Code < 600 {
			err5xxpermin[timekey]++
		}
	}

	at.PeakRequestTime, at.PeakRequestCount = getPeakTime(requestpermin)
	at.Peak4xxTime, at.Peak4xxCount = getPeakTime(err4xxpermin)
	at.Peak5xxTime, at.Peak5xxCount = getPeakTime(err5xxpermin)

	return
}

func getPeakTime(records map[string]int) (string, int) {

	max := 0
	timeStamp := ""

	for time, count := range records {
		if count > max {
			max = count
			timeStamp = time
		}
	}
	return timeStamp, max
}
