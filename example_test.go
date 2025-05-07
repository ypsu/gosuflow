package gosuflow_test

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ypsu/gosuflow"
)

type renderTableWorkflow struct {
	headerLabels []string
	data         [][]string

	ExtractDimensionsSection struct{}
	rows, columns            int
	columnWidths             []int

	CenterHeaderLabelsSection struct{}
	centeredHeaderLabels      []string

	GenerateHeaderSection struct{}
	header                string

	GenerateBodySection struct{}
	body                string

	PrintSection struct{}
}

func (wf *renderTableWorkflow) ExtractDimensions(ctx context.Context) error {
	wf.rows, wf.columns = len(wf.data), len(wf.headerLabels)
	wf.columnWidths = make([]int, wf.columns)
	for c, label := range wf.headerLabels {
		wf.columnWidths[c] = utf8.RuneCountInString(label)
	}
	for r, row := range wf.data {
		if len(row) != wf.rows {
			return fmt.Errorf("rendertable.BadRowColumnsCount row=%d columns=%d want=%d", r, len(row), wf.rows)
		}
		for c, v := range row {
			wf.columnWidths[c] = max(wf.columnWidths[c], utf8.RuneCountInString(v))
		}
	}
	return nil
}

func (wf *renderTableWorkflow) CenterHeaderLabels(ctx context.Context) error {
	wf.centeredHeaderLabels = make([]string, wf.columns)
	for c, label := range wf.headerLabels {
		spaces := wf.columnWidths[c] - utf8.RuneCountInString(label) + 2
		wf.centeredHeaderLabels[c] = strings.Repeat(" ", spaces/2) + label + strings.Repeat(" ", spaces/2+spaces%2)
	}
	return nil
}

func (wf *renderTableWorkflow) writeSeparatorLine(w *strings.Builder) {
	for _, c := range wf.columnWidths {
		w.WriteByte('+')
		for range c + 2 {
			w.WriteByte('-')
		}
	}
	w.WriteByte('+')
	w.WriteByte('\n')
}

func (wf *renderTableWorkflow) GenerateHeader(ctx context.Context) error {
	w := &strings.Builder{}
	wf.writeSeparatorLine(w)
	for _, label := range wf.centeredHeaderLabels {
		w.WriteByte('|')
		w.WriteString(label)
	}
	w.WriteByte('|')
	w.WriteByte('\n')
	wf.writeSeparatorLine(w)
	wf.header = w.String()
	return nil
}

func (wf *renderTableWorkflow) GenerateBody(ctx context.Context) error {
	w := &strings.Builder{}
	for _, row := range wf.data {
		for c, label := range row {
			w.WriteByte('|')
			w.WriteByte(' ')
			for range wf.columnWidths[c] - utf8.RuneCountInString(label) {
				w.WriteByte(' ')
			}
			w.WriteString(label)
			w.WriteByte(' ')
		}
		w.WriteByte('|')
		w.WriteByte('\n')
	}
	wf.writeSeparatorLine(w)
	wf.body = w.String()
	return nil
}

func (wf *renderTableWorkflow) Print(ctx context.Context) error {
	fmt.Print(wf.header)
	fmt.Print(wf.body)
	return nil
}

func Example() {
	ctx := context.Background()
	wf := &renderTableWorkflow{
		headerLabels: []string{"City", "Population", "Area"},
		data: [][]string{
			[]string{"Dublin", "592,713", "117 km²"},
			[]string{"London", "8,866,180", "1,572 km²"},
			[]string{"Paris", "2,048,472", "105 km²"},
		},
	}
	if err := gosuflow.Run(ctx, wf); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Output:
	// +--------+------------+-----------+
	// |  City  | Population |   Area    |
	// +--------+------------+-----------+
	// | Dublin |    592,713 |   117 km² |
	// | London |  8,866,180 | 1,572 km² |
	// |  Paris |  2,048,472 |   105 km² |
	// +--------+------------+-----------+
}
