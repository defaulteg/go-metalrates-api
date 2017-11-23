package core

import (
	"errors"
	"os/exec"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func (s *Source) GetRatesFromStaticSource() error {

	// If source doesn't have any selectors for rate searching in source
	if len(s.Selectors) == 0 {
		return nil
	}
/*
	reader, _ := os.Open("./src/gitlab.com/defaulteg/api/scripts/temp/temp_data.html")
	defer reader.Close()

	_, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return errors.New("Cannot create new document from reader")
	}
*/
	// Download static html document from url
	doc, err := goquery.NewDocument(s.Site) //"http://www.kitco.com/market/"
	if err != nil {
		return err
	}

	// For each selector(html id name) search for this id field in downloaded html
	for i, selector := range s.Selectors {
		var text string
			doc.Find("#"+selector.Name).Each(func(i int, s *goquery.Selection) {
				text = s.Text()			//if element found -> grab it's value(rate)
		})

		if text == "" {
			s.Selectors[i].Rate = "-1"	//if element not found rate value will be -1
		} else {
			s.Selectors[i].Rate = text	//if element found assign it's rate to selector rate field
		}
	}

	return nil
}


// Get rates sequentially by running phantomjs each time for each source
func (s *Source) GetRatesFromDynamicSource() error {

	// If source doesn't have any selectors for rate searching in source
	if len(s.Selectors) == 0 {
		return nil
	}

	// Make string vararg for phantomjs cmd command parameters
	cmdParams := make([]string, 0)
	// Add first two parameters like: path to script, source url
	cmdParams = append(cmdParams, PathToPageElementFetcher, s.Site)

	// Add each selector name to parameter slice
	for _, selector := range s.Selectors {
		cmdParams = append(cmdParams, selector.Name)
	}

	// Execute cmd command; It will return slice of rates
	cmd := exec.Command("phantomjs", cmdParams...)

	if res, err := cmd.Output(); err != nil {
		return errors.New("Cannot execute phantomjs command. " + err.Error());
	} else {
		pjsOutputTemp := strings.TrimSpace(string(res)) 					//Trim all spaces from phantomjs response
		pjsOutput := strings.Split(pjsOutputTemp, "\n")						//Split response string to slice with '\n' divider
		if pjsOutput[0] == "Error." {										//If zero element is Error then [1] element contains error description
			return errors.New("Invalid phantomjs output\n" + pjsOutput[1])
			//fmt.Print(pjsOutput[1])
		} else {
			for i, _ := range s.Selectors {									//Assign each rate from phantomjs output to selector element
				//fmt.Println(s.Selectors[i].metalId)
				s.Selectors[i].Rate = pjsOutput[i]
			}
		}
	}
	return nil
}


