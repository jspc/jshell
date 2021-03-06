package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	header = color.MagentaString(`
	   8 8888   d888888o.   8 8888        8 8 8888888888   8 8888         8 8888
	   8 8888 .'8888:' '88. 8 8888        8 8 8888         8 8888         8 8888
	   8 8888 8.'8888.   Y8 8 8888        8 8 8888         8 8888         8 8888
	   8 8888 '8.'8888.     8 8888        8 8 8888         8 8888         8 8888
	   8 8888  '8.'8888.    8 8888        8 8 888888888888 8 8888         8 8888
	   8 8888   '8.'8888.   8 8888        8 8 8888         8 8888         8 8888
88.        8 8888    '8.'8888.  8 8888888888888 8 8888         8 8888         8 8888
'88.       8 888'8b   '8.'8888. 8 8888        8 8 8888         8 8888         8 8888
  '88o.    8 88' '8b.  ;8.'8888 8 8888        8 8 8888         8 8888         8 8888
    'Y888888 '    'Y8888P ,88P' 8 8888        8 8 888888888888 8 888888888888 8 888888888888
		`)
)

func main() {
	for {
		fmt.Print("\033[H\033[2J")
		fmt.Println(header)
		a, err := AppMenu()
		if err != nil {
			panic(err)
		}

		fmt.Print("\033[H\033[2J")

		for _, f := range []func() error{
			a.Run,
			a.Cleanup,
		} {
			err = f()
			if err != nil {
				panic(err)
			}
		}
	}
}
