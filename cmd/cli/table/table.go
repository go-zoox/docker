package table

import (
	"fmt"
	"os"

	"github.com/InVisionApp/tabular"
	"github.com/go-zoox/core-utils/object"
)

type Option struct {
	Clear bool
}

type Column struct {
	Key             string `json:"key"`
	Title           string `json:"title"`
	Width           int    `json:"width"`
	DisableEllipsis bool   `json:"disable_ellipsis"`
}

type Item = map[string]string

func Table(columns []Column, dataSource []Item, opts ...func(*Option)) error {
	opt := &Option{}
	for _, o := range opts {
		o(opt)
	}

	if opt.Clear {
		os.Stdout.Write([]byte("\033[2J\033[H"))
	}

	tab := tabular.New()

	for _, column := range columns {
		if column.Key == "" {
			column.Key = column.Title
		}

		tab.Col(column.Key, column.Title, column.Width)
	}

	format := tab.Print("*")
	for _, item := range dataSource {
		values := []any{}
		for _, column := range columns {
			value := object.Get(item, column.Key)
			if len(value) > column.Width {
				if !column.DisableEllipsis {
					value = value[:column.Width-3] + "..."
				} else {
					value = value[:column.Width]
				}
			}

			values = append(values, value)
		}

		fmt.Printf(format, values...)
	}

	return nil
}

func WithClearOption() func(*Option) {
	return func(opt *Option) {
		opt.Clear = true
	}
}
