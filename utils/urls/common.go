package urls

import (
	"embed"
	"fmt"
	"strings"

	"file-inspector/utils"
)

//go:embed alexa-top-100000.txt
var alexaData embed.FS

const (
	alexafilePath = "alexa-top-100000.txt"
)

type URLChecker interface {
	Check(urlString string) (bool, error)
	CountKnownDomains() int
}

type CommonChecker struct {
	URLChecker
	data       map[string]bool
	dataLoaded bool
}

// Need a builder method so we don't have to load the data for each check
// might be unnecessary as the embedded data doesn't need dynamic loading
// but interesting to try it anyway 
func GetCommonURLChecker() (URLChecker, error) {
	return newChecker(), nil
}

func newChecker() URLChecker {
	c := &CommonChecker{}
	c.data = loadAlexTop1000()
	c.dataLoaded = true

	return c
}

func (c *CommonChecker) CountKnownDomains() int {
	return len(c.data)
}

// Check the domain if it's in the Alexa list of top 100k hosts
func (c *CommonChecker) Check(urlString string) (bool, error) {
	
	// should have been loaded when created, but let's check
	if !c.dataLoaded {
		return false, fmt.Errorf("no domain data loaded")
	}

	hostname := utils.GetHostFromURL(urlString)

	if _, ok := c.data[hostname]; ok {
		return true, nil
	}

	return false, nil
}

func loadAlexTop1000() map[string]bool {
	bytes, _ := alexaData.ReadFile(alexafilePath)
	lines := strings.Split(string(bytes), "\n")

	domains := make(map[string]bool, len(lines))

	for _, line := range lines {
		domains[line] = true
	}

	return domains
}
