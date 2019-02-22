// Copyright © 2019 En-Hao Hu <enhao.mobile@gmail.com>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const tagTempl = `{{if eq .query_status "ok"}}Malware Tag Infomation:
  First seen: {{.firstseen}}
  {{if .lastseen}}Last seen:  {{.lastseen}}{{end}}
  Number of URLs observation: {{.url_count}}

  List of malware URLs associated with this tag (max 1000):{{range .urls}}
    * {{.url}}
      Status:       {{.url_status}}
      URLhaus:
        ID:         {{.url_id}}
        Reference:  {{.urlhaus_reference}}
        Date added: {{.dateadded}}
        Reporter:   {{.reporter}}
{{end}}{{else}}{{.query_status}}{{end}}
`

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Get information about a tag",
	Long:  `This command retrieves information about a tag.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Post(URL("tag"), "application/x-www-form-urlencoded", strings.NewReader("tag="+args[0]))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if len(b) == 0 {
			return
		}

		if rawOutput {
			fmt.Printf("%s", b)
			return
		}

		t := template.Must(template.New("").Parse(tagTempl))

		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(b), &m); err != nil {
			log.Fatal(err)
		}

		if err := t.Execute(os.Stdout, m); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
