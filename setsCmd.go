package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"text/template"

	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/maleck13/mtga/mtgio"
	"github.com/urfave/cli"
)

var refreshSetsFlag bool

const (
	SETS_DATA = ".mtga_sets.json"
	SET_DATA  = ".mtga_set_%s.json"
)

func SetsCmd() cli.Command {
	return cli.Command{
		Name:        "get",
		Description: "get data about sets",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "sets",
				Action: sets,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:        "refresh",
						Destination: &refreshSetsFlag,
						Usage:       "forces a fresh pull of the the sets data",
					},
				},
			},
			cli.Command{
				Name:   "set",
				Action: set,
				Usage:  "set <set_code>",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:        "refresh",
						Destination: &refreshSetsFlag,
						Usage:       "forces a fresh pull of the the sets data",
					},
				},
			},
		},
	}
}

var setsTemplate = `{{range . }}| Name: {{.Name}} | Code: {{.Code}}  |
{{end}}`
var setTemplate = `{{range . }}| Name: {{.Name}} | Type: {{.Type}} | Rarity {{.Rarity}}  | Score {{.Score}}
{{end}}`

func GetSets() ([]*mtgio.Set, error) {
	home := homedir.Get()
	path := filepath.Join(home, SETS_DATA)
	var (
		f    *os.File
		sets []*mtgio.Set
	)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		refreshSetsFlag = true
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, mtgio.NewToolError("failed to stat data "+err.Error(), 1)
	}
	defer f.Close()
	if refreshSetsFlag {
		sets, err = mtgio.GetSets()
		if err != nil {
			return nil, mtgio.NewToolError(err.Error(), 1)
		}
		enc := json.NewEncoder(f)
		if err := enc.Encode(sets); err != nil {
			return nil, mtgio.NewToolError("failed to encode and save sets "+err.Error(), 1)
		}
	}

	if nil == sets {
		dec := json.NewDecoder(f)
		if err := dec.Decode(&sets); err != nil {
			return nil, mtgio.NewToolError("failed to decode sets data "+err.Error(), 1)
		}
	}
	return sets, nil
}

func GetSet(name string) (*mtgio.SetCards, error) {
	home := homedir.Get()
	path := filepath.Join(home, fmt.Sprintf(SET_DATA, name))
	var (
		f   *os.File
		set *mtgio.SetCards
	)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		refreshSetsFlag = true
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, mtgio.NewToolError("failed to stat data "+err.Error(), 1)
	}
	defer f.Close()
	if refreshSetsFlag {
		set, err = mtgio.GetSet(name)
		if err != nil {
			return nil, mtgio.NewToolError(err.Error(), 1)
		}
		enc := json.NewEncoder(f)
		if err := enc.Encode(set); err != nil {
			return nil, mtgio.NewToolError("failed to encode and save sets "+err.Error(), 1)
		}
	}
	if nil == set {
		dec := json.NewDecoder(f)
		if err := dec.Decode(&set); err != nil {
			return nil, mtgio.NewToolError("failed to decode sets data "+err.Error(), 1)
		}
	}

	return set, nil
}

func set(context *cli.Context) error {
	if len(context.Args()) != 1 {
		return mtgio.NewToolError("you need to pass a set code "+context.Command.Usage, 1)
	}
	setName := context.Args().First()
	set, err := GetSet(setName)
	if err != nil {
		return mtgio.NewToolError("failed to get set "+err.Error(), 1)
	}
	t := template.New("set")
	t, err = t.Parse(setTemplate)
	if err != nil {
		return mtgio.NewToolError("failed to parse template "+err.Error(), 1)
	}
	if err := t.Execute(os.Stdout, set.Cards); err != nil {
		return mtgio.NewToolError("failed to execute template "+err.Error(), 1)
	}
	return nil
}

func sets(context *cli.Context) error {

	sets, err := GetSets()
	if err != nil {
		return mtgio.NewToolError("failed to get sets "+err.Error(), 1)
	}
	t := template.New("sets")
	t, err = t.Parse(setsTemplate)
	if err != nil {
		return mtgio.NewToolError("failed to parse template "+err.Error(), 1)
	}
	if err := t.Execute(os.Stdout, sets); err != nil {
		return mtgio.NewToolError("failed to execute template "+err.Error(), 1)
	}

	return nil
}
