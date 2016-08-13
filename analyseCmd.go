package main

import (
	"github.com/maleck13/mtga/mtgio"
	"github.com/urfave/cli"
	"os"
	"text/template"
)

func AnalyseCmd() cli.Command {
	return cli.Command{
		Name:   "analyse",
		Action: analyse,
		Usage:  "analyse <set_code>",
	}
}

var creatureKeyWords = map[string]int{
	"flying":       2,
	"first strike": 2,
	"vigilance":2,
	"lifelink":2,
	"menace":1,
	"skulk":1,
	"hexproof":2,
	"haste":2,
	"trample":1,
}

var setAnalysisTemplate = `{{range . }} {{if gt .Score 2}}
| Name: {{.Name}} | Type: {{.Type}} | Rarity {{.Rarity}}  | Score {{.Score}}
{{end}}{{end}}
`

func analyse(cont *cli.Context) error {
	if len(cont.Args()) != 1 {
		return mtgio.NewToolError("missing set argument "+cont.Command.Usage, 1)
	}
	set := cont.Args().First()
	setData, err := GetSet(set)
	if err != nil {
		return err
	}
	if err := applyKeyWordValues(setData); err != nil{
		return err
	}
	if err := applyEnterBattleFieldAnalysis(setData); err != nil{
		return err
	}
	if err := applyWhenCastAnalysis(setData); err != nil{
		return err
	}
	if err := applyCostPowerAnalysis(setData); err != nil{
		return err
	}
	t := template.New("set")
	t, err = t.Parse(setAnalysisTemplate)
	if err != nil {
		return mtgio.NewToolError("failed to parse template "+err.Error(), 1)
	}
	if err := t.Execute(os.Stdout, setData.Cards); err != nil {
		return mtgio.NewToolError("failed to execute template "+err.Error(), 1)
	}
	return nil

}

func applyKeyWordValues(data *mtgio.SetCards)error {
	for _, card := range data.Cards {
		for keyWord,rate := range creatureKeyWords {
			if (card.IsCreature() && card.HasKeyword(keyWord)){
				card.IncrementScore(rate)
			}
		}
	}
	return nil
}

func applyEnterBattleFieldAnalysis(data *mtgio.SetCards)error{
	for _, card := range data.Cards {
		if card.HasEnterBattleFieldEffect(){
			card.IncrementScore(1)
		}
	}
	return nil
}

func applyWhenCastAnalysis(data *mtgio.SetCards)error{
	for _, card := range data.Cards {
		if card.HasWhenCastEffect(){
			card.IncrementScore(1)
		}
	}
	return nil
}

func applyCostPowerAnalysis(data *mtgio.SetCards)error{
       return nil
}
