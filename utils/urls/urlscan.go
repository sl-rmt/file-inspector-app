package urls

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/existentiality/urlscan"
)

const (
	urlscan_api_key_env = "URLSCAN_API_KEY"
)

type SearchResult struct {
	Results []struct {
		Task struct {
			Visibility string    `json:"visibility"`
			Method     string    `json:"method"`
			Domain     string    `json:"domain"`
			ApexDomain string    `json:"apexDomain"`
			Time       time.Time `json:"time"`
			UUID       string    `json:"uuid"`
			URL        string    `json:"url"`
			Tags       []string  `json:"tags,omitempty"`
		} `json:"task,omitempty"`
		Stats struct {
			UniqIPs           int `json:"uniqIPs"`
			UniqCountries     int `json:"uniqCountries"`
			DataLength        int `json:"dataLength"`
			EncodedDataLength int `json:"encodedDataLength"`
			Requests          int `json:"requests"`
		} `json:"stats"`
		Page struct {
			Country      string    `json:"country"`
			Server       string    `json:"server"`
			IP           string    `json:"ip"`
			MimeType     string    `json:"mimeType"`
			Title        string    `json:"title"`
			URL          string    `json:"url"`
			TLSValidDays int       `json:"tlsValidDays"`
			TLSAgeDays   int       `json:"tlsAgeDays"`
			Ptr          string    `json:"ptr,omitempty"`
			TLSValidFrom time.Time `json:"tlsValidFrom"`
			Domain       string    `json:"domain"`
			UmbrellaRank int       `json:"umbrellaRank"`
			ApexDomain   string    `json:"apexDomain"`
			Asnname      string    `json:"asnname"`
			Asn          string    `json:"asn"`
			TLSIssuer    string    `json:"tlsIssuer"`
			Status       string    `json:"status"`
		} `json:"page,omitempty"`
		ID         string        `json:"_id"`
		Sort       []interface{} `json:"sort"`
		Result     string        `json:"result"`
		Screenshot string        `json:"screenshot"`
	} `json:"results"`
	Total   int  `json:"total"`
	Took    int  `json:"took"`
	HasMore bool `json:"has_more"`
}

func SubmitToURLScan(host string) error {
	apiKey, err := getURLScanAPIKeyFromEnv()

	if err != nil {
		return err
	}

	// https://urlscan.io/docs/api/
	client := urlscan.NewClient(apiKey)

	// search first
	query := fmt.Sprintf("domain:%s", host)
	result, _ := client.Search(query, 100)

	log.Printf("Got %d results from URLScan. Result 1:\n", len(result.Results))

	// Print first set of results
	if len(result.Results) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 1, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, "Unique IPs\t%d\n", result.Results[0].Stats.UniqIPs)
		fmt.Fprintf(w, "IP\t%s\n", result.Results[0].Page.IP)
		fmt.Fprintf(w, "Unique Countries\t%d\n", result.Results[0].Stats.UniqCountries)
		fmt.Fprintf(w, "Country\t%s\n", result.Results[0].Page.Country)
		fmt.Fprintf(w, "Server\t%s\n", result.Results[0].Page.Server)
		fmt.Fprintf(w, "ASN Name\t%s\n", result.Results[0].Page.Asnname)
		fmt.Fprintf(w, "TLS Valid Days\t%d\n", result.Results[0].Page.TLSValidDays)
		fmt.Fprintf(w, "TLS Age\t%d\n", result.Results[0].Page.TLSAgeDays)

		//fullResults := result.Results[0].Result
	} //else {
		// // if no results submit scan
		// resp, err := client.Scan(host, urlscan.ScanOptions{
		// 	Country: "gb",
		// })

		// if err != nil {
		// 	return err
		// }

		// if resp.Message != "Submission successful" {
		// 	return fmt.Errorf("unexpected response message: %s", resp.Message)
		// }

		// result, err := client.GetResult(resp.UUID)

		// if err != nil {
		// 	return err
		// }

	//}

	return nil
}

func getURLScanAPIKeyFromEnv() (string, error) {
	key := os.Getenv(urlscan_api_key_env)

	if key == "" {
		return "", fmt.Errorf("failed. We need a URLScan API key set as the environment variable %q", urlscan_api_key_env)
	}

	log.Println("Got URLScan API key from env")

	return key, nil
}
