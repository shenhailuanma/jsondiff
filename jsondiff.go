package jsondiff

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type JsonDiffItem struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

const (
	JsonDiffTypeAdded   = "Added"
	JsonDiffTypeChanged = "Changed"
	JsonDiffTypeRemoved = "Removed"
)

func JsonDiff(jsonA, jsonB string) ([]JsonDiffItem, error) {
	return jsonDiff(jsonA, jsonB, "")
}

func jsonDiff(jsonA, jsonB string, prefix string) ([]JsonDiffItem, error) {
	var changes = []JsonDiffItem{}

	if jsonA == jsonB {
		// directly same, return
		return changes, nil
	}

	a := gjson.Parse(jsonA)
	b := gjson.Parse(jsonB)

	var aOnlyKeys = []string{}
	var bOnlyKeys = []string{}

	if a.IsObject() && b.IsObject() {
		a.ForEach(func(aKey, aValue gjson.Result) bool {
			var onlyInA = true
			b.ForEach(func(bKey, bValue gjson.Result) bool {
				if aKey.String() == bKey.String() {
					onlyInA = false

					// check
					subChanges, _ := jsonDiff(aValue.Raw, bValue.Raw, prefix+"."+aKey.String())
					for _, subChangeOne := range subChanges {
						changes = append(changes, subChangeOne)
					}
				}
				return true
			})

			if onlyInA == true {
				aOnlyKeys = append(aOnlyKeys, aKey.String())
				changes = append(changes, JsonDiffItem{
					Type: JsonDiffTypeRemoved,
					Path: prefix + "." + aKey.String(),
				})
			}
			return true
		})

		// keys only in b
		b.ForEach(func(bKey, bValue gjson.Result) bool {
			var onlyInB = true
			a.ForEach(func(aKey, avalue gjson.Result) bool {
				if aKey.String() == bKey.String() {
					onlyInB = false
				}
				return true
			})

			if onlyInB == true {
				bOnlyKeys = append(bOnlyKeys, bKey.String())
				changes = append(changes, JsonDiffItem{
					Type: JsonDiffTypeAdded,
					Path: prefix + "." + bKey.String(),
				})
			}
			return true
		})
	} else if a.IsArray() && b.IsArray() {

		aArray := a.Array()
		bArray := b.Array()

		if len(aArray) == len(bArray) {
			var arrayExistChanges = false
			for i := 0; i < len(aArray); i++ {
				if (aArray[i].IsObject() && bArray[i].IsObject()) || (aArray[i].IsArray() && bArray[i].IsArray()) {
					subChanges, _ := jsonDiff(aArray[i].Raw, bArray[i].Raw, fmt.Sprintf("%s.%d", prefix, i))
					if len(subChanges) > 0 {
						arrayExistChanges = true
					}
				} else {
					if aArray[i].Raw != bArray[i].Raw {
						arrayExistChanges = true
					}
				}
			}

			if arrayExistChanges {
				changes = append(changes, JsonDiffItem{
					Type: JsonDiffTypeChanged,
					Path: prefix,
				})
			}
		} else {
			if jsonA != jsonB {
				changes = append(changes, JsonDiffItem{
					Type: JsonDiffTypeChanged,
					Path: prefix,
				})
			}
		}
	} else {
		if jsonA != jsonB {
			changes = append(changes, JsonDiffItem{
				Type: JsonDiffTypeChanged,
				Path: prefix,
			})
		}
	}

	return changes, nil
}
