package replacebot

import (
	"math/rand/v2"
	"strings"

	"github.com/jdkato/prose/v2"
)

type Config struct {
	CallResponses  []CallResponse        `yaml:"call_responses"`
	Replacements   map[string][]Response `yaml:"replacements"`
	ReplaceChance  float32               `yaml:"replace_chance"`
	ResponseChance float32               `yaml:"response_chance"`
	UserName       string                `yaml:"bot_user_name"`
	OAuth          string                `yaml:"oauth_token"`
	Channels       []string              `yaml:"channels"`
}

type CallResponse struct {
	Calls     []string   `yaml:"calls"`
	Responses []Response `yaml:"responses"`
}

type Response struct {
	Chance  float32 `yaml:"chance"`
	Message string  `yaml:"response"`
}

type ReplaceBot struct {
	ResponseChance float32
	ReplaceChance  float32
	CallResponses  []CallResponse
	Replacements   map[string][]Response
}

func NewReplaceBot(cfg *Config) *ReplaceBot {
	return &ReplaceBot{
		ResponseChance: cfg.ResponseChance,
		ReplaceChance:  cfg.ReplaceChance,
		CallResponses:  cfg.CallResponses,
		Replacements:   cfg.Replacements,
	}
}

func (rb *ReplaceBot) Respond(m string) *string {
	// Handle specific call and responses
	for _, cr := range rb.CallResponses {
		for _, call := range cr.Calls {
			if strings.Contains(strings.ToLower(m), call) {
				w := rand.Float32()
				for _, resp := range cr.Responses {
					if resp.Chance >= w {
						return &resp.Message
					}
				}
			}
		}
	}

	// Now do our actual string replacements
	if rand.Float32() < rb.ResponseChance {
		doc, err := prose.NewDocument(m)
		if err != nil {
			return nil
		}
		tks := doc.Tokens()
		msg := make([]string, len(tks))
		chg := false
		for i, tok := range tks {
			msg[i] = tok.Text
			if rand.Float32() < rb.ReplaceChance {
				w := rand.Float32()
				var rpls []Response
				switch tok.Tag {
				case "NN":
					rpls = rb.Replacements["single_nouns"]
				case "NNP":
					rpls = rb.Replacements["proper_nouns"]
				case "NNS":
					rpls = rb.Replacements["plural_nouns"]
				case "JJ":
					rpls = rb.Replacements["adjectives"]
				case "RB":
					rpls = rb.Replacements["adverbs"]
				case "VB":
					fallthrough
				case "VBP":
					rpls = rb.Replacements["verbs"]
				case "VBD":
					rpls = rb.Replacements["verbs_past"]
				case "VBG":
					rpls = rb.Replacements["verbs_ing"]
				}

				for _, rpl := range rpls {
					if rpl.Chance >= w {
						msg[i] = rpl.Message
						chg = true
						break
					}
				}
			}
		}
		if chg {
			r := strings.Join(msg, " ")
			return &r
		}
	}
	return nil
}
