package quit

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

type Quit struct{}

func (Quit) Name() string        { return "Quit jShell" }
func (Quit) Description() string { return "Quit jShell, exiting from this server 🙁" }
func (Quit) Cleanup() error      { return nil }
func (Quit) Run() error {
	prompt := promptui.Prompt{
		Label:     "Are you sure?",
		IsConfirm: true,
	}

	for {
		result, err := prompt.Run()
		switch result {
		case "y", "Y":
			fmt.Print("\033[H\033[2J")
			fmt.Println("bye-bye")

			os.Exit(0)

		case "n", "N", "":
			return nil

		default:
			if err != nil {
				break
			}
		}
	}

	return nil
}
