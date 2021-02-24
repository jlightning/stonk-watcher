package services

func GetStockInformation(ticker string) error {
	//finvizStockInfo, err := GetDataFromFinviz(ticker)
	//if err != nil {
	//	return err
	//}

	err := GetDataFromMarketWatch(ticker)
	if err != nil {
		return err
	}

	//fmt.Printf("%#v", finvizStockInfo)
	return nil
}
