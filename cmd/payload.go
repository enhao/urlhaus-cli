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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const payloadTempl = `{{if eq .query_status "ok"}}Malware Payload Infomation:
  Type: {{.file_type}}
  Size: {{.file_size}}
  Hash:
    - MD5:    {{.md5_hash}}
    - SHA256: {{.sha256_hash}}

  {{if .signature}}Signature:  {{.signature}}{{end}}
  First seen: {{.firstseen}}
  {{if .lastseen}}Last seen:  {{.lastseen}}{{end}}
  Number of URLs observation: {{.url_count}}
  URLhaus reference: {{.urlhaus_reference}}
  Sample download: {{.urlhaus_download}}
  {{if .virustotal}}VirusTotal: {{.virustotal.percent}}% ({{.virustotal.result}})
    Link:     {{.virustotal.link}}
  {{end}}
  List of malware URLs associated with this payload (max 100):{{range .urls}}
    * {{.url}}
      Status:       {{.url_status}}
      URLhaus:
        Reference:  {{.urlhaus_reference}}
        First seen: {{.firstseen}}{{if .lastseen}}
        Last seen:  {{.lastseen}}{{end}}
{{end}}{{else}}{{.query_status}}{{end}}
`

var hashType string

// payloadCmd represents the payload command
var payloadCmd = &cobra.Command{
	Use:   "payload",
	Short: "Get information about a payload (malware sample)",
	Long: `This command retrieves information about a payload (malware sample) that
URLhaus has retrieved.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var data io.Reader
		if hashType == "sha256" {
			data = strings.NewReader("sha256_hash=" + args[0])
		} else {
			data = strings.NewReader("md5_hash=" + args[0])
		}

		resp, err := http.Post(URL("payload"), "application/x-www-form-urlencoded", data)
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

		t := template.Must(template.New("").Parse(payloadTempl))

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
	rootCmd.AddCommand(payloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// payloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// payloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	payloadCmd.Flags().StringVarP(&hashType, "type", "t", "md5", "The hash type of the payload")
}
