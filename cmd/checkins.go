// Copyright © 2019 Juha Ristolainen <juha.ristolainen@iki.fi>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twpayne/go-kml"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var checkinsCmd = &cobra.Command{
	Use:   "checkins",
	Short: "get your foursquare checkins as KML",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := viper.GetString("accessToken")
		filename := viper.GetString("output")

		fmt.Printf("Output file: %s\n", filename)
		f, err := os.Create(filename)
		check(err)
		w := bufio.NewWriter(f)

		getAllCheckins(w, accessToken)

		defer f.Close()
	},
}

func init() {
	RootCmd.AddCommand(checkinsCmd)
	checkinsCmd.Flags().StringP("accessToken", "", "", "access token (required)")
	checkinsCmd.MarkFlagRequired("accessToken")

	checkinsCmd.Flags().StringP("output", "o", "", "output file")
	checkinsCmd.MarkFlagRequired("output")

	viper.BindPFlag("accessToken", checkinsCmd.Flags().Lookup("accessToken"))
	viper.BindPFlag("output", checkinsCmd.Flags().Lookup("output"))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	pageSize = 250
)

func getAllCheckins(outputFile *bufio.Writer, accessToken string) {
	checkinsURI := getPaginatedURI(1, 0, accessToken)

	resp, err := http.Get(checkinsURI)
	check(err)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		checkinsTotalCount, err := jsonparser.GetInt(body, "response", "checkins", "count")
		pageCount := int(math.Floor(float64(checkinsTotalCount)/float64(pageSize)) + 1)
		fmt.Printf("- Total number of check-ins: %d\n", checkinsTotalCount)
		fmt.Printf("- %d pages of %d check-ins\n", pageCount, pageSize)

		// Get all pages and add to a KML-document
		kDoc := kml.Document()
		k := kml.KML(kDoc)

		for i := 0; i < pageCount; i++ {
			fmt.Printf("Getting check-ins %d to %d (page %d of %d)\n", i*pageSize, i*pageSize+pageSize, i, pageCount)
			getCheckins(pageSize, i*pageSize, accessToken, kDoc)
		}

		if err := k.WriteIndent(outputFile, "", "  "); err != nil {
			log.Fatal(err)
		}
	}
}

func getCheckins(pageSize int, offset int, accessToken string, kDoc *kml.CompoundElement) {
	// First get the 250 latest checkins. API is paginated and max page size is 250.
	checkinsURI := getPaginatedURI(pageSize, offset, accessToken)
	//fmt.Println(checkinsURI)

	resp, err := http.Get(checkinsURI)
	check(err)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			check(err)
			createdAtEpoch, err := jsonparser.GetInt(value, "createdAt")
			createdAtTime := time.Unix(createdAtEpoch, 0)
			checkinVenue, dataType, offset, err := jsonparser.Get(value, "venue")
			venueName, err := jsonparser.GetString(checkinVenue, "name")
			checkinVenueLocation, dataType, offset, err := jsonparser.Get(checkinVenue, "location")
			venueAddress, err := jsonparser.GetString(checkinVenueLocation, "address")
			venueCity, err := jsonparser.GetString(checkinVenueLocation, "city")
			venuePostCode, err := jsonparser.GetString(checkinVenueLocation, "postalCode")
			venueState, err := jsonparser.GetString(checkinVenueLocation, "state")
			venueCountry, err := jsonparser.GetString(checkinVenueLocation, "country")

			addressLines := []string{venueAddress, venueCity, venuePostCode, venueState, venueCountry}
			address := strings.Join(addressLines, ", ")
			venueLat, err := jsonparser.GetFloat(checkinVenueLocation, "lat")
			venueLng, err := jsonparser.GetFloat(checkinVenueLocation, "lng")
			kDoc.Add(kml.Placemark(
				kml.TimeStamp(kml.When(createdAtTime)),
				kml.Name(venueName),
				kml.Address(address),
				kml.Point(kml.Coordinates(kml.Coordinate{Lon: venueLng, Lat: venueLat})),
			),
			)
		}, "response", "checkins", "items")

		defer resp.Body.Close()
	} else {
		fmt.Printf("HTTP Status %d", resp.StatusCode)
	}
}

func getPaginatedURI(limit int, offset int, accessToken string) string {
	return fmt.Sprintf("https://api.foursquare.com/v2/users/self/checkins?limit=%d&offset=%d&sort=newestfirst&v=20190715&oauth_token=%s", limit, offset, accessToken)
}
