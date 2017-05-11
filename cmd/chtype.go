// Copyright © 2017 Tillman Mosley III <tmosley@dermpathlab.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// chtypeCmd represents the chtype command
var chtypeCmd = &cobra.Command{
	Use:   "chtype",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			accNum       string
			origCaseType string
			caseType     string
			casePKey     int
			caseYear     int
			seqNumber    int
			seq_number   *int64
			case_type_id *string
			year         *int64
			IsLetter     = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
		)

		if flags := len(args); flags < 2 {
			fmt.Println("usage: chtype <accNum> <caseType>")
			os.Exit(1)
		}

		if IsLetter(args[1]) {
			caseType = strings.ToUpper(args[1])
		} else {
			fmt.Println("Case Type must be a string")
			os.Exit(1)
		}

		accNum = strings.ToUpper(args[0])

		timestamp := time.Now().Format(time.RFC3339Nano)

		origCaseType, caseYear, seqNumber = FindCaseSubstrings(accNum)

		fmt.Printf("%s %s %d %s %d\n", accNum, caseType, caseYear, origCaseType, seqNumber)

		cases, err := GetCase(seqNumber, origCaseType, caseYear)
		if err != nil {
			log.Fatal(err)
		}

		casePKey = cases.ID
		seq_number = cases.Seq_number
		case_type_id = cases.Case_type_id
		year = cases.Year
		fmt.Println(*case_type_id, *year, *seq_number, casePKey)

		stmt2, err := db.Prepare("UPDATE public.case set case_type_id=$1 WHERE seq_number = $2 AND case_type_id = $3 AND year = $4")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()
		newCaseType := *cases.Case_type_id

		caseTypeMatch := strings.Compare(origCaseType, newCaseType)

		if caseTypeMatch == 0 {
			result, err := stmt2.Exec(caseType, seqNumber, origCaseType, caseYear)
			if err != nil {
				log.Fatal(err)
			}
			affect, err := result.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(affect, "rows affected")

			insertStmt := `
			INSERT INTO case_note (item_id, note, created_at, application_id, case_id,
			content_type_id, user_id, show_on_report)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			_, err = db.Exec(insertStmt, casePKey, "Changed case type to "+caseType, timestamp, 15, casePKey, 43, 14, false)
			if err != nil {
				log.Fatal(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(chtypeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chtypeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chtypeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
