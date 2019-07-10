package reporting

import (
	"fmt"
)

var report Report

type Report struct {
	ReceivedPerPeriod   uint64
	NumbersWrittenTotal uint64
	UniqPerPeriod       uint64
	DuplicatesPerPeriod uint64
	ReceivedTotal       uint64
}

func Duplicated() {
	report.DuplicatesPerPeriod++
	report.ReceivedPerPeriod++
	report.ReceivedTotal++
}
func Uniq() {
	report.NumbersWrittenTotal++
	report.UniqPerPeriod++
	report.ReceivedPerPeriod++
	report.ReceivedTotal++
}

func PrintReport() {
	fmt.Printf("%+v\n", report)

	report.DuplicatesPerPeriod = 0
	report.UniqPerPeriod = 0
	report.ReceivedPerPeriod = 0
}
