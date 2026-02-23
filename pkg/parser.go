package pkg

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	clearprefix    = regexp.MustCompile(`^[^|]*\|\s*`)
	tracebackStart = regexp.MustCompile(`^Traceback \(most recent call last\):`)
	errorRe        = regexp.MustCompile(`\[ERROR\]`)
	hasher         = sha1.New()
	gunicornLogRe  = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3}`)
)

func ParseLogFile(name string) (LogData, error) {

	file, err := os.Open(name)
	scanner := bufio.NewScanner(file)
	//USING 10 MB BUFFER FOR LARGE TOKEN
	buffer := make([]byte, 10*1024*1024)
	scanner.Buffer(buffer, 10*1024*1024)

	newReport := LogData{}
	accessRecords := []AccessRecord{}
	errordb := make(ErrorDB)

	var lasterror string
	var inTrace bool
	var tblines []string

	for scanner.Scan() {
		line := scanner.Text()
		content := clearprefix.ReplaceAllString(line, "")

		if inTrace {
			if gunicornLogRe.MatchString(content) {
				errordb.AddTraceErrorRecord(lasterror, tblines)
				tblines = []string{}
				inTrace = false
			} else {
				tblines = append(tblines, content)
			}
		}

		if strings.Contains(line, "[gunicorn.access]") {
			ar := parseAccessLogs(content)
			accessRecords = append(accessRecords, ar)
			continue
		}

		if errorRe.MatchString(line) {
			cleanerror := regexp.MustCompile(`^.*?\[ERROR\]\s*`).ReplaceAllString(line, "")
			lasterror = getErrorKey(cleanerror)
			errorrecord := ErrorRecord{Message: cleanerror}
			errordb.AddErrorRecord(lasterror, errorrecord)
		}

		if tracebackStart.MatchString(content) {
			tblines = append(tblines, content)
			inTrace = true
			continue
		}

	}

	defer file.Close()
	newReport.AccessRecords = accessRecords
	newReport.Errors = errordb.getErrorRecord()
	return newReport, err
}
func parseAccessLogs(line string) (a AccessRecord) {
	words := strings.Fields(line)

	parsedDate, err := time.Parse("2006-01-02 15:04:05", words[0]+" "+words[1])
	if err != nil {
		log.Println(err.Error())
		return
	}

	a.Date = parsedDate
	a.Ip = words[5]
	a.Url = words[11]

	a.Code, _ = strconv.Atoi(words[13])
	a.Size, _ = strconv.Atoi(words[14])

	for _, s := range words[16:] {

		a.UserAgent = a.UserAgent + s
	}

	return
}

func getErrorKey(e string) string {
	sum := hasher.Sum([]byte(e))
	return hex.EncodeToString(sum[:11])
}
