/*
Copyright © 2020 Christopher Maahs <cmaahs@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/cmaahs/no-more-lateness/calendar"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display a list of upcoming calendar entries",
	Long:  `Display a list of upcoming online meeting items.`,
	Run: func(cmd *cobra.Command, args []string) {
		attendeeAddress, _ = cmd.Flags().GetString("attendee-address")
		displayMeetings()
	},
}

func displayMeetings() {

	cal, err := calendar.GetProvider("google")
	if err != nil {
		fmt.Println("bad")
		os.Exit(1)
	}

	_, cerr := cal.GetClient()
	if cerr != nil {
		fmt.Println("bad")
		os.Exit(1)
	}

	out, eerr := cal.GetEvents(5, attendeeAddress)
	if eerr != nil {
		fmt.Println("no events")
	}

	if len(out) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		// if !noHeaders {
		table.SetHeader([]string{"START", "SOON", "EVENT", "GOING", "LINK"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		// }
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t") // pad with tabs
		table.SetNoWhiteSpace(true)
		table.SetColumnColor(tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor})

		nextMeeting := 0
		for _, evt := range out {
			var ml int
			if ml = len(evt.MeetingLink.String()); ml > 80 {
				ml = 80
			}
			row := []string{evt.Start.Format("2006-01-02 15:04"), fmt.Sprintf("%t", evt.IsMeetingSoon), evt.Description, evt.MeetingResponse, evt.MeetingLink.String()[0:ml]}
			minutesUntilStart := time.Until(evt.Start).Minutes()
			if minutesUntilStart >= 0 {
				nextMeeting++
			}
			if evt.IsMeetingSoon || nextMeeting == 1 {
				table.Rich(row, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}})
			} else {
				table.Append(row)
			}
		}
		table.Render()
	}

}

func init() {
	viewCmd.Flags().StringP("attendee-address", "a", "", "Specify your attendee email address")
	viewCmd.MarkFlagRequired("attendee-address")

	rootCmd.AddCommand(viewCmd)

}
