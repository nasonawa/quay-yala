package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/nasonawa/quay-yala/pkg"
	"github.com/nasonawa/quay-yala/report"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var inputFiles []string

var RootCmd = &cobra.Command{
	Use:   "quay-yala",
	Short: "quay-yala",
	Long:  "Log analyzer for the Red Hat Quay logs",
	RunE:  runRootCmd,
	Args:  cobra.NoArgs,
}

func init() {

	RootCmd.Flags().StringSliceVarP(&inputFiles, "input", "i", []string{}, "Input Files (Comma separated or repeated)")
}

func runRootCmd(cmd *cobra.Command, args []string) (err error) {
	if len(inputFiles) == 0 {
		cmd.Help()
		return nil
	}
	for _, f := range inputFiles {
		err = generateReportForFile(f)
	}

	return err

}

func generateReportForFile(name string) error {

	data, err := pkg.ParseLogFile(name)
	if err != nil {
		log.Println(err)
	}

	outfile, err := os.Create("report-" + filepath.Base(name))
	if err != nil {
		log.Println(err)
	}

	AccessReport := report.AccessLogTextReportAnalysis(data.AccessRecords)

	textTemplate, err := template.New("access-log").Parse(report.Textreport)

	textTemplate.Execute(outfile, AccessReport)
	fmt.Fprintln(outfile, "")

	header := []string{"Ip", "Date", "Url", "Code", "Size", "UserAgent"}
	accesstable := tablewriter.NewWriter(outfile)
	accesstable.SetHeader(header)

	for _, v := range data.AccessRecords {
		accesstable.Append([]string{v.Ip, v.Date.Format("2006-01-02 15:04:05"), v.Url,
			fmt.Sprintf("%d", v.Code), fmt.Sprintf("%d", v.Size), v.UserAgent})
	}
	fmt.Fprintln(outfile, "## Access log summary from "+name+"#")
	accesstable.Render()
	fmt.Fprintln(outfile, "## ERROR summary from "+name+"#")

	errorstable := tablewriter.NewWriter(outfile)
	errorstable.SetAutoWrapText(false)
	errorstable.SetRowLine(true)
	errorstable.SetHeaderLine(true)
	errorstable.SetHeader([]string{"Count", "Error", "TraceBack"})

	sort.Slice(data.Errors, func(i, j int) bool {

		return data.Errors[i].Count > data.Errors[j].Count
	})

	for _, e := range data.Errors {
		errorstable.Append([]string{fmt.Sprintf("%d", e.Count), e.Message, e.Traceback})
	}
	errorstable.Render()
	return nil
}
