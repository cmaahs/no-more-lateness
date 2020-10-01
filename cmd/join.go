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

	"github.com/cmaahs/no-more-lateness/calendar"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a meeting if it is within 5 minutes +/- of the start time.",
	Long: `This command will find the next meeting on your calendar, determine
	if it is within the join time, and launch the online meeting tool that is
	associated.

	EXAMPLE
	  #> no-more-lateness join`,
	Run: func(cmd *cobra.Command, args []string) {
		joinMeetings()
	},
}

func joinMeetings() {

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

	out, eerr := cal.GetEvents(5)
	if eerr != nil {
		fmt.Println("no events")
	}

	if len(out) > 0 {
		for _, evt := range out {
			if evt.IsMeetingSoon {
				_ = open.Run(evt.MeetingLink.String())
				fmt.Printf("%v,(%v),<%s>\n", evt.Description, evt.Start, evt.MeetingLink.String())
			}
			// fmt.Println(fmt.Sprintf("Soon: %t, Event: %s, Link: %s", evt.IsMeetingSoon, evt.Description, evt.MeetingLink.String()))
		}
	}

}

func init() {
	rootCmd.AddCommand(viewCmd)

}

func init() {
	rootCmd.AddCommand(joinCmd)

}
