package dot_test

import (
	"encoding/json"
	"testing"

	"github.com/jopbrown/gobase/object/dot"
	"github.com/stretchr/testify/assert"
)

var data = []byte(`{
	"id": "1000",
	"type": "web",
	"name": "service",
	"ppm": 0.55,
	"qos": 1,
	"numarr": [1, 2, 3],
	"batters":
		{
			"batter":
				[
					{ "id": "3011", "type": "Apple" },
					{ "id": "3012", "type": "Book" },
					{ "id": "3013", "type": "Cake" },
					{ "id": "3014", "type": "Disk" }
				]
		},
	"topping":
		[
			{ "id": "8101", "type": "Egg" },
			{ "id": "8102", "type": "Folder" },
			{ "id": "8105", "type": "Gemma" },
			{ "id": "8107", "type": "Hint" },
			{ "id": "8106", "type": "Information" },
			{ "id": "8103", "type": "Jackal" },
			{ "id": "8104", "type": "King" }
		]
}`)

func TestGet_Map(t *testing.T) {
	mdata := make(map[string]any)
	err := json.Unmarshal(data, &mdata)
	assert.NoError(t, err)

	assert.Equal(t, "1000", dot.Get[string](mdata, "id"))
	assert.Equal(t, 0.55, dot.Get[float64](mdata, "ppm"))
	assert.Equal(t, 1.0, dot.Get[float64](mdata, "qos"))

	assert.Equal(t, "Apple", dot.Get[string](mdata, "batters.batter[0].type"))
	assert.Equal(t, "3014", dot.Get[string](mdata, "batters.batter[3].id"))
	assert.Panics(t, func() { dot.Get[string](mdata, "batters.batter[5].id") })

	batters := dot.Get[any](mdata, "batters")
	assert.Equal(t, "3013", dot.Get[string](batters, "batter[2].id"))

	assert.Equal(t, "Jackal", dot.Get[string](mdata, "topping[5].type"))

	topping := dot.Get[any](mdata, "topping")
	dot.Set(topping, "[1].id", "9999")
	assert.Equal(t, "9999", dot.Get[string](topping, "[1].id"))

	dot.Set(mdata, "numarr[1]", 9.9)
	assert.Equal(t, 9.9, dot.Get[float64](mdata, "numarr[1]"))
}

func TestGet_Struct(t *testing.T) {
	type Topping struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	type Batter struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	type Batters struct {
		Batter []Batter `json:"batter"`
	}

	type Model struct {
		ID      string     `json:"id"`
		Type    string     `json:"type"`
		Name    string     `json:"name"`
		Ppm     float64    `json:"ppm"`
		Qos     int        `json:"qos"`
		NumArr  []int      `json:"numarr"`
		Batters Batters    `json:"batters"`
		Topping []*Topping `json:"topping"`
	}

	mdata := &Model{}
	err := json.Unmarshal(data, &mdata)
	assert.NoError(t, err)

	assert.Equal(t, "1000", dot.Get[string](mdata, "id"))
	assert.Equal(t, 0.55, dot.Get[float64](mdata, "ppm"))
	assert.Equal(t, 1, dot.Get[int](mdata, "qos"))

	assert.Equal(t, "Apple", dot.Get[string](mdata, "batters.batter[0].type"))
	assert.Equal(t, "3014", dot.Get[string](mdata, "batters.batter[3].id"))
	assert.Panics(t, func() { dot.Get[string](mdata, "batters.batter[5].id") })

	batters := dot.Get[Batters](mdata, "batters")
	assert.Equal(t, "3013", dot.Get[string](batters, "batter[2].id"))

	assert.Equal(t, "Jackal", dot.Get[string](mdata, "topping[5].type"))

	topping := dot.Get[any](mdata, "topping")
	dot.Set(topping, "[1].id", "9999")
	assert.Equal(t, "9999", mdata.Topping[1].ID)

	dot.Set(mdata, "numarr[1]", 9)
	assert.Equal(t, 9, mdata.NumArr[1])
}
