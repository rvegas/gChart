package main

import (
	"github.com/rvegas/gchart"
	"fmt"
	"os"
)

func main() {

	testChart := gchart.NewLineChart()
	testChart.SetHeight(400)
	testChart.SetWidth(900)
	testChart.SetTitle("TESTING")

	testChart.AddHeader("date", "date")
	testChart.AddHeader("amount", "number")
	testChart.SetCurvedLine()

	testChart.LoadCSV("./resources/data.csv", false)

	result, err := testChart.Generate()
	check(err)

	file, _ := os.Create("/tmp/index.html")
	fmt.Fprint(file, result)

	fmt.Println("Done")
	exec.Command("xdg-open", fmt.Sprintf("file:///%s", file.Name())).Start()
}

// fatal if there is an error
func check(err error) {
	if err != nil {
		panic(err)
	}
}
