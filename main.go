package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Author = "craig brookes"
	app.Name = "mtga"
	app.Usage = "mtga analyse --set=EMN"
	app.Commands = []cli.Command{
		AnalyseCmd(),
		SetsCmd(),
	}
	app.Run(os.Args)

}
