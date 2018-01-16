package gchart

import (
	"encoding/csv"
	"fmt"
	"github.com/rvegas/gchart/resources"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"errors"
)

// HeaderTypes are the available and compatible data types for values in each column
var HeaderTypes = map[string]bool{
	"date":   true,
	"string": true,
	"number": true,
}

// Main struct to hold chart information
type gchart struct {
	headers   []dataHeader
	data      []dataRow
	style     string
	height    int
	width     int
	title     string
	extraOpts map[string]string
}

// Struct for data rows
type dataRow struct {
	Values []interface{}
}

// Struct for data headers
type dataHeader struct {
	Name string
	Type string
}

// A gChart MUST always be initiated first
func initiate() *gchart {
	newChart := &gchart{}
	newChart.width = 900
	newChart.height = 400
	newChart.title = "GChart"
	newChart.extraOpts = make(map[string]string)
	return newChart
}

// NewLineChart creates a new GChart object with a Line Chart Style
func NewLineChart() *gchart {
	newChart := initiate()
	newChart.style = "LineChart"
	return newChart
}

// NewAreaChart creates a new GChart object with an Area Chart Style
func NewAreaChart() *gchart {
	newChart := initiate()
	newChart.style = "AreaChart"
	return newChart
}

// NewColumnChart creates a new GChart object with a Column Chart Style
func NewColumnChart() *gchart {
	newChart := initiate()
	newChart.style = "ColumnChart"
	return newChart
}

// Generates the header string for the chart
func (c gchart) generateHeaders() string {
	headers := ""
	for _, header := range c.headers {
		headers += fmt.Sprintf("\t\tdata.addColumn('%s', '%s');\n", header.Type, header.Name)
	}
	return headers
}

// Generates the rows string for the chart
func (c gchart) generateRows() string {
	rows := "\t\tdata.addRows([\n"

	for _, data := range c.data {
		row := "\t\t\t["
		for _, value := range data.Values {
			switch value.(type) {
			case float64, float32:
				row += fmt.Sprintf("%f,", value)
			case int, int64:
				row += fmt.Sprintf("%d,", value)
			case time.Time:
				date, _ := value.(time.Time)
				row += fmt.Sprintf("new Date(%d, %d, %d, %d, %d, %d),", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
			case string:
				row += fmt.Sprintf("'%s',", value)
			}
		}
		row = strings.TrimRight(row, ",") + "],\n"
		rows += row
	}
	rows += "\t\t]);\n"
	return rows
}

// Executes the drawing command in the html file
func (c gchart) generateDraw() string {
	return fmt.Sprintf("%s%s%s", resources.DRAW_BEGIN, c.style, resources.DRAW_END)
}

// Generates default and useful html chart options
func (c gchart) generateOptions() string {
	options := ""
	options += fmt.Sprintf("\t\ttitle:'%s',\n", c.title)
	options += fmt.Sprintf("\t\twidth:%d,\n", c.width)
	options += fmt.Sprintf("\t\theight:%d,\n", c.height)

	for opt, value := range c.extraOpts {
		options += fmt.Sprintf("\t\t%s:'%s',\n", opt, value)
	}

	return options
}

// Generate renders the GChart into a static HTML file.
func (c gchart) Generate() (string, error) {
	err := c.Validate()
	if err != nil {
		return "", err
	}
	return "" +
			resources.HEADER_BEGIN +
			resources.DATA_BEGIN +
			c.generateHeaders() +
			c.generateRows() +
			resources.OPT_BEGIN +
			c.generateOptions() +
			resources.OPT_END +
			resources.DRAW_BEGIN +
			c.style +
			resources.DRAW_END +
			resources.DATA_END +
			resources.HEADER_END +
			resources.BODY +
			resources.FOOTER,
		nil
}

// Validate makes sure that the gChart is valid
func (c gchart) Validate() error {
	i := 0
	if c.title == "" {
		return fmt.Errorf("the chart is missing a required title")
	}
	if c.width == 0 {
		return fmt.Errorf("the chart is missing a required width")
	}
	if c.height == 0 {
		return fmt.Errorf("the chart is missing a required height")
	}
	for _, row := range c.data {
		if len(c.headers) != len(row.Values) {
			return fmt.Errorf("the number of values in row %d does not match the header count", i)
		}
		i++
	}
	return nil
}

// SetTitle sets a title for the gChart which will show in the top, centered
func (c *gchart) SetTitle(title string) error {
	c.title = title
	return nil
}

// SetWidth sets the desired with in pixels for the gChart
func (c *gchart) SetWidth(width int) error {
	c.width = width
	return nil
}

// SetCurvedLine modifies the current gChart Line Chart with a curved line style
func (c *gchart) SetCurvedLine() error {
	if c.style != "LineChart" {
		return errors.New("incompatible chart style, it MUST be a line chart")
	}
	c.extraOpts["curveType"] = "function"
	return nil
}

// SetHeight sets the desired height of the gChart
func (c *gchart) SetHeight(height int) error {
	c.height = height
	return nil
}

// AddHeader adds a column header for the gChart data
func (c *gchart) AddHeader(name string, headerType string) error {
	if !HeaderTypes[headerType] {
		return fmt.Errorf("header type %s not supported", headerType)
	}
	c.headers = append(c.headers, dataHeader{Name: name, Type: headerType})
	return nil
}

// AddRow adds a row of column values to the gChart data
func (c *gchart) AddRow(values []interface{}) error {
	dataRow := dataRow{}
	dataRow.Values = values
	c.data = append(c.data, dataRow)

	return nil
}

// LoadCSV tries to load a correctly formed CSV file into gChart data
func (c *gchart) LoadCSV(filename string, loadHeaders bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)
	i := 0

	for {
		record, err := reader.Read()

		if !loadHeaders && i == 0 {
			i++
			continue
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var values []interface{}
		for _, element := range record {
			if val, err := strconv.ParseInt(element, 10, 64); err == nil {
				values = append(values, val)
			} else if val, err := strconv.ParseFloat(element, 64); err == nil {
				values = append(values, val)
			} else if val, err := strconv.ParseBool(element); err == nil {
				values = append(values, val)
			} else if val, err := time.Parse(time.UnixDate, element); err == nil {
				values = append(values, val)
			} else if val, err := time.Parse(time.RFC3339, element); err == nil {
				values = append(values, val)
			} else {
				values = append(values, val)
			}
		}
		c.AddRow(values)
	}

	return nil
}
