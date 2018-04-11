package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/rvegas/gchart"
)

func main() {

	testChart := gchart.NewLineChart()
	testChart.SetHeight(700)
	testChart.SetWidth(2000)

	testChart.AddHeader("run #", "number")
	testChart.AddHeader("time", "number")
	testChart.AddHeader("req/s", "number")
	testChart.AddHeader("cpu", "number")
	testChart.AddHeader("mem", "number")
	testChart.SetCurvedLine()
	//testChart.SetLogarithmicVerticalAxis()

	pwd, _ := os.Getwd()
	//TODO: Make this friendlier. If it gives you trouble, just feed the abs path in your fs
	testChart.LoadCSV(pwd+"/github.com/rvegas/gchart/examples/resources/startend.csv", false)

	result, err := testChart.Generate()
	check(err)

	file, _ := os.Create("/tmp/startend.html")
	fmt.Fprint(file, result)

	fmt.Println("Done")
	exec.Command("xdg-open", fmt.Sprintf("file:///%s", file.Name())).Start()

	fs := http.FileServer(http.Dir("/tmp"))
	http.Handle("/", fs)

	fmt.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

// fatal if there is an error
func check(err error) {
	if err != nil {
		panic(err)
	}
}
