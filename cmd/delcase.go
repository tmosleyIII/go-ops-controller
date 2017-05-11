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
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// delcaseCmd represents the delcase command
var delcaseCmd = &cobra.Command{
	Use:   "delcase",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		var (
			accNum       string
			origCaseType string
			caseYear     int
			seqNumber    int
			casePKey     int
		)

		if flags := len(args); flags < 1 {
			fmt.Println("usage: delcase <accNum>")
			os.Exit(1)
		}

		accNum = strings.ToUpper(args[0])
		timestamp := time.Now().Format(time.RFC3339Nano)

		origCaseType, caseYear, seqNumber = FindCaseSubstrings(accNum)
		cases, err := GetCase(seqNumber, origCaseType, caseYear)
		if err != nil {
			log.Fatal(err)
		}

		casePKey = cases.ID

		updateStmt := `
		UPDATE public.case set status=98, deleted_at=$1 WHERE seq_number = $2 AND case_type_id = $3 AND year = $4 AND deleted_at IS NULL`
		stmt, err := db.Prepare(updateStmt)
		if err != nil {
			log.Fatal(err)
		}

		result, err := stmt.Exec(timestamp, seqNumber, origCaseType, caseYear)
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
		_, err = db.Exec(insertStmt, casePKey, "Case deleted", timestamp, 15, casePKey, 43, 1, false)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(delcaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// delcaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// delcaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
