package dictionary

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Results struct {
	ID       *string   `json:"id,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	Results  []Result  `json:"results"`
	Word     *string   `json:"word,omitempty"`
}

type Metadata struct {
	Operation *string `json:"operation,omitempty"`
	Provider  *string `json:"provider,omitempty"`
	Schema    *string `json:"schema,omitempty"`
}

type Result struct {
	ID             *string        `json:"id,omitempty"`
	Language       *string        `json:"language,omitempty"`
	LexicalEntries []LexicalEntry `json:"lexicalEntries"`
	Type           *string        `json:"type,omitempty"`
	Word           *string        `json:"word,omitempty"`
}

type LexicalEntry struct {
	Entries         []Entry          `json:"entries"`
	Language        *string          `json:"language,omitempty"`
	LexicalCategory *LexicalCategory `json:"lexicalCategory,omitempty"`
	Pronunciations  []Pronunciation  `json:"pronunciations"`
	Text            *string          `json:"text,omitempty"`
}

func (entry *LexicalEntry) RenderLexicalCategory() string {
	return strings.ToUpper(*entry.LexicalCategory.Text)
}

type Entry struct {
	Etymologies     []string `json:"etymologies"`
	HomographNumber *string  `json:"homographNumber,omitempty"`
	Senses          []Sense  `json:"senses"`
}

func (e Entry) HasEtymology() bool {
	return len(e.Etymologies) > 0
}

func (e Entry) RenderEtymology() interface{} {
	var buf bytes.Buffer
	for _, ex := range e.Etymologies {
		fmt.Fprintf(&buf, "%s\n", ex)
	}
	return buf.String()
}

type Sense struct {
	Definitions      []string        `json:"definitions"`
	Examples         []Example       `json:"examples"`
	ID               *string         `json:"id,omitempty"`
	ShortDefinitions []string        `json:"shortDefinitions"`
	Subsenses        []Subsense      `json:"subsenses"`
	ThesaurusLinks   []ThesaurusLink `json:"thesaurusLinks"`
}

func (s Sense) RenderExamples() string {
	var buf bytes.Buffer
	for _, ex := range s.Examples {
		fmt.Fprintf(&buf, "%s\n", ex.Render())
	}
	return buf.String()
}

func (s Sense) HasExamples() bool {
	return len(s.Examples) > 0
}

func (ex Example) Render() string {
	return *ex.Text
}

type Example struct {
	Text *string `json:"text,omitempty"`
}

type Subsense struct {
	Definitions      []string          `json:"definitions"`
	Examples         []Example         `json:"examples"`
	ID               *string           `json:"id,omitempty"`
	Regions          []LexicalCategory `json:"regions"`
	Registers        []LexicalCategory `json:"registers"`
	ShortDefinitions []string          `json:"shortDefinitions"`
	ThesaurusLinks   []ThesaurusLink   `json:"thesaurusLinks"`
}

func (s Subsense) HasExamples() bool {
	return len(s.Examples) > 0
}

func (s Subsense) RenderTags() string {
	tags := make([]string, 0)

	for _, t := range s.Regions {
		tags = append(tags, cleanTag(*t.Text))
	}

	for _, t := range s.Registers {
		tags = append(tags, cleanTag(*t.Text))
	}

	if len(tags) == 0 {
		return ""
	}
	return color.GreenString(fmt.Sprintf("%s ", strings.Join(tags, ", ")))
}

func (s Subsense) RenderExamples() string {
	var buf bytes.Buffer
	for _, ex := range s.Examples {
		fmt.Fprintf(&buf, "%s\n", ex.Render())
	}
	return buf.String()
}

func cleanTag(s string) (output string) {
	output = s
	output = strings.Replace(output, "_", " ", -1)
	return
}

type LexicalCategory struct {
	ID   *string `json:"id,omitempty"`
	Text *string `json:"text,omitempty"`
}

type ThesaurusLink struct {
	EntryID *string `json:"entry_id,omitempty"`
	SenseID *string `json:"sense_id,omitempty"`
}

type Pronunciation struct {
	AudioFile        *string  `json:"audioFile,omitempty"`
	Dialects         []string `json:"dialects"`
	PhoneticNotation *string  `json:"phoneticNotation,omitempty"`
	PhoneticSpelling *string  `json:"phoneticSpelling,omitempty"`
}
