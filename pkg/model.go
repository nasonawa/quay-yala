package pkg

import (
	"maps"
	"slices"
	"strings"
	"time"
)

type LogData struct {
	AccessRecords []AccessRecord
	Errors        []ErrorRecord
}

type AccessRecord struct {
	Ip        string
	Date      time.Time
	Url       string
	Code      int
	Size      int
	UserAgent string
}

type ErrorRecord struct {
	Message   string
	Traceback string
	Count     int
}

type ErrorDB map[string]ErrorRecord

func (edb ErrorDB) AddTraceErrorRecord(key string, trace []string) {

	t := strings.Join(trace, "\n")
	if v, ok := edb[key]; ok {

		if len(v.Traceback) < len(t) {
			v.Traceback = t
			edb[key] = v
		}
	}
}

func (edb ErrorDB) AddErrorRecord(key string, e ErrorRecord) {

	if v, ok := edb[key]; ok {
		v.Count += 1
		edb[key] = v
	} else {
		e.Count += 1
		edb[key] = e
	}
}

func (edb ErrorDB) getErrorRecord() []ErrorRecord {
	return slices.Collect(maps.Values(edb))
}
