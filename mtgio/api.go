package mtgio

import (
	"encoding/json"
	"fmt"
	"runtime"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"github.com/urfave/cli"
	"strings"
)

const (
	HOST     = "https://api.magicthegathering.io"
	SETS_URL = "/v1/sets"
	SET_URL  = "/v1/cards?set=%s"
)

type ToolError struct {
	*cli.ExitError
	Message string
	Line    int
	File    string
	Code    int
}

func (ae *ToolError) Error() string {
	return fmt.Sprintf("api error: %s file %s line %d", ae.Message, ae.File, ae.Line)
}

func NewToolError(message string, code int) *ToolError {
	_, f, n, _ := runtime.Caller(1)
	cliErr := cli.NewExitError(message, code)
	return &ToolError{
		cliErr,
		message,
		n,
		f,
		code,
	}
}

type Set struct {
	Block              string        `json:"block"`
	Booster            []interface{} `json:"booster"`
	Border             string        `json:"border"`
	Code               string        `json:"code"`
	MagicCardsInfoCode string        `json:"magicCardsInfoCode"`
	MkmID              int           `json:"mkm_id"`
	MkmName            string        `json:"mkm_name"`
	Name               string        `json:"name"`
	ReleaseDate        string        `json:"releaseDate"`
	Type               string        `json:"type"`
}

type Sets struct {
	Sets []*Set `json:"sets"`
}

func GetSets() ([]*Set, error) {
	result := &Sets{}

	response, err := http.Get(HOST + SETS_URL)
	if err != nil {
		return nil, NewToolError("failed to get sets "+err.Error(), 1)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, NewToolError("failed to decode response "+err.Error(), 1)
	}
	return result.Sets, nil
}

type SetCards struct {
	Cards []*Card `json:"cards"`
}

type Card struct {
	Artist       string   `json:"artist"`
	ID           string   `json:"id"`
	ImageURL     string   `json:"imageUrl"`
	Layout       string   `json:"layout"`
	Multiverseid int      `json:"multiverseid"`
	Name         string   `json:"name"`
	Number       string   `json:"number"`
	Printings    []string `json:"printings"`
	Rarity       string   `json:"rarity"`
	Set          string   `json:"set"`
	SetName      string   `json:"setName"`
	Text         string   `json:"text"`
	Type         string   `json:"type"`
	Types        []string `json:"types"`
	ManaCost     string `json:"manaCost"`
	CMC 	     string `json:"cmc"`
	Score        int      `json:"score"`
	Cost         float64  `json:"cost"`
}

func (c *Card)HasKeyword(word string)bool{
	return strings.Contains(strings.ToLower(c.Text),strings.ToLower(word))
}

func (c *Card)IncrementScore(amount int)  {
	c.Score +=amount
}

func (c *Card)IsCreature()bool{
	for _, v := range c.Types{
		if v == "Creature"{
			return true
		}
	}
	return false
}

func (c *Card)HasEnterBattleFieldEffect()bool{
	toFind := strings.ToLower("when " + c.Name + " enters the battlefield")
	return strings.Contains(toFind,strings.ToLower(c.Text))
}

func (c *Card)HasWhenCastEffect()bool{
	toFind := strings.ToLower("when you cast " + c.Name)
	return strings.Contains(toFind,strings.ToLower(c.Text))
}

func GetSet(name string) (*SetCards, error) {
	url := fmt.Sprintf(HOST+SET_URL, name)
	fmt.Println("getting set "+name, url)
	result := &SetCards{}
	response, err := http.Get(url)
	if err != nil {
		return nil, NewToolError("failed to get sets "+err.Error(), 1)
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, NewToolError("failed to decode response "+err.Error(), 1)
	}
	fmt.Println("cards", len(result.Cards))
	return result, nil
}
