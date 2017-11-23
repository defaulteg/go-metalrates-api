package metal

import (
	"github.com/defaulteg/api/modules/core"
	"github.com/defaulteg/api/utils"
)


//getSources -> getSelectors -> getRates(from site, sequentially) -> pushToDatabase

//TODO: Run each ??? minute
//TODO: Get historical rates


func Fetch() error {

	if err := utils.FileExists(core.PathToPageFetcher); err != nil {
		return err
	}

	if sources, err := core.GetSources(core.MetalCategory); err != nil {
		return err
	} else {
		for _, source := range sources {

			//download page?
			/*
			cmd := exec.Command("phantomjs", path2, "http://www.kitco.com/market/")

			if res, err := cmd.Output(); err != nil {
				fmt.Println("err")
			} else {
				pjsOutput := strings.TrimSpace(string(res))
				fmt.Print(pjsOutput + "successsss")
			}*/

			if err := source.GetRatesFromStaticSource(); err != nil {	//source.getRates2()
				return err
			}
		}

		//fmt.Println(sources)

		// Push to database
		if err := core.PushToDatabase(sources, core.MetalTable); err != nil {
			return err
		}
	}

	return nil
}



