package robot

import (
	"github.com/pintobikez/bitstamp/api"
	"strconv"
	"strings"
	"time"
)

type Rules struct {
	Selldrop   int8     `yaml:"selldrop,omitempty"`
	Sellup     int8     `yaml:"sellup,omitempty"`
	Buydrop    int8     `yaml:"buydrop,omitempty"`
	Currencies []string `yaml:"currencies,omitempty"`
}

type Engine struct {
	r   *Rules
	api *api.API
}

func New(ru *Rules, ap *api.API) *Engine {
	return &Engine{r: ru, api: ap}
}

func (e *Engine) RunRules() {
	bl, err := e.api.AccountBalance()
	if err != nil {

	}
}
