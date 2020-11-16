package search

import (
	"strings"

	"github.com/baking-bad/bcdhub/internal/helpers"
	"github.com/pkg/errors"
)

// Scorable -
type Scorable interface {
	GetFields() []string
	GetIndex() string
	GetScores(string) []string
	Parse(highlight map[string][]string, data []byte) (interface{}, error)
}

// indices
const (
	contractIndex   = "contract"
	operationIndex  = "operation"
	bigmapdiffIndex = "bigmapdiff"
	tzipIndex       = "tzip"
	metadataIndex   = "metadata" // TODO: constants in separate package
)

// Indices - list of indices availiable to search
var Indices = []string{
	contractIndex,
	operationIndex,
	bigmapdiffIndex,
	tzipIndex,
}

var scorables = map[string]Scorable{
	contractIndex:   &Contract{},
	operationIndex:  &Operation{},
	bigmapdiffIndex: &BigMap{},
	tzipIndex:       &Token{},
	metadataIndex:   &Metadata{},
}

// ScoreInfo -
type ScoreInfo struct {
	Scores  []string
	Indices []string

	indicesMap map[string]struct{}
	fieldsMap  map[string]struct{}
}

func newScoreInfo() ScoreInfo {
	return ScoreInfo{
		Scores:  make([]string, 0),
		Indices: make([]string, 0),

		indicesMap: make(map[string]struct{}),
		fieldsMap:  make(map[string]struct{}),
	}
}

func (si *ScoreInfo) addIndex(index string) {
	if _, ok := si.indicesMap[index]; ok {
		return
	}
	si.indicesMap[index] = struct{}{}
	si.Indices = append(si.Indices, index)
}

func (si *ScoreInfo) addScore(score string) {
	val := strings.Split(score, "^")
	field := val[0]
	if _, ok := si.fieldsMap[field]; ok {
		return
	}
	si.fieldsMap[field] = struct{}{}
	si.Scores = append(si.Scores, score)
}

func (si *ScoreInfo) addScores(scores ...string) {
	for i := range scores {
		si.addScore(scores[i])
	}
}

// GetScores -
func GetScores(searchString string, fields []string, indices ...string) (ScoreInfo, error) {
	info := newScoreInfo()
	if len(indices) == 0 {
		for i := range Indices {
			info.addIndex(Indices[i])
		}
		return info, nil
	}

	for i := range indices {
		model, ok := scorables[indices[i]]
		if !ok {
			return info, errors.Errorf("[GetSearchScores] Unknown scorable model: %s", indices[i])
		}
		index := model.GetIndex()
		if helpers.StringInArray(index, Indices) {
			modelScores := model.GetScores(searchString)
			info.addIndex(index)
			if len(fields) > 0 {
				for i := range modelScores {
					for j := range fields {
						if strings.HasPrefix(modelScores[i], fields[j]) {
							info.addScore(modelScores[i])
						}
					}
				}
			} else {
				info.addScores(modelScores...)
			}
		}
	}

	return info, nil
}

// Parse -
func Parse(index string, highlight map[string][]string, data []byte) (interface{}, error) {
	fields := make([]string, 0)
	for key := range highlight {
		fields = append(fields, key)
	}

	switch index {
	case contractIndex:
		return scorables[index].Parse(highlight, data)
	case operationIndex:
		return scorables[index].Parse(highlight, data)
	case bigmapdiffIndex:
		return scorables[index].Parse(highlight, data)
	case tzipIndex:
		token := scorables[index]

		var found bool
		for _, field := range token.GetFields() {
			for _, h := range fields {
				if field == h {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			return token.Parse(highlight, data)
		}
		return scorables[metadataIndex].Parse(highlight, data)
	default:
		return nil, errors.Errorf("Unknown index: %s", index)
	}
}

// Result -
type Result struct {
	Count int64  `json:"count"`
	Time  int64  `json:"time"`
	Items []Item `json:"items"`
}

// Item -
type Item struct {
	Type       string              `json:"type"`
	Value      string              `json:"value"`
	Group      *Group              `json:"group,omitempty"`
	Body       interface{}         `json:"body"`
	Highlights map[string][]string `json:"highlights,omitempty"`

	Network string `json:"-"`
}

// Group -
type Group struct {
	Count int64 `json:"count"`
	Top   []Top `json:"top"`
}

// NewGroup -
func NewGroup(docCount int64) *Group {
	return &Group{
		Count: docCount,
		Top:   make([]Top, 0),
	}
}

// Top -
type Top struct {
	Network string `json:"network"`
	Key     string `json:"key"`
}
