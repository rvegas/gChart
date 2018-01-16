package main

import (
	"github.com/rvegas/gchart"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	testChart := gchart.NewLineChart()
	testChart.SetHeight(400)
	testChart.SetWidth(900)
	testChart.SetTitle("TESTING")

	testChart.AddHeader("date", "date")
	testChart.AddHeader("amount", "number")
	testChart.SetCurvedLine()

	pwd, _ := os.Getwd()
	//TODO: Make this friendlier. If it gives you trouble, just feed the abs path in your fs
	testChart.LoadCSV(pwd + "/github.com/rvegas/gChart/examples/resources/data.csv", false)

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
