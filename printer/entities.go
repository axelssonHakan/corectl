package printer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/enigma-go"
)

// PrintGenericEntities prints a list of the id and type of all generic entities in the app
func PrintGenericEntities(allInfos []*enigma.NxInfo, entityType string, printAsJSON bool, printAsBash bool) {
	if printAsJSON {
		specifiedEntityTypeInfos := []*enigma.NxInfo{}
		for _, info := range allInfos {
			if (entityType == "object" && info.Type != "measure" && info.Type != "dimension") || entityType == info.Type {
				specifiedEntityTypeInfos = append(specifiedEntityTypeInfos, info)
			}
		}
		buffer, err := json.Marshal(specifiedEntityTypeInfos)
		if err != nil {
			internal.FatalError(err)
		}
		fmt.Println(prettyJSON(buffer))
	} else if printAsBash {
		for _, info := range allInfos {
			if (entityType == "object" && info.Type != "measure" && info.Type != "dimension") || entityType == info.Type {
				PrintToBashComp(info.Id)
			}
		}
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Type"})

		for _, info := range allInfos {
			if (entityType == "object" && info.Type != "measure" && info.Type != "dimension") || entityType == info.Type {
				writer.Append([]string{info.Id, info.Type})
			}
		}
		writer.Render()
	}
}

// PrintGenericEntityProperties prints the properties of the generic entity defined by entityID
func PrintGenericEntityProperties(state *internal.State, entityID string, entityType string) {
	var err error
	var properties json.RawMessage
	switch entityType {
	case "object":
		genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		properties, err = genericObject.GetPropertiesRaw(state.Ctx)
	case "measure":
		genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		properties, err = genericMeasure.GetPropertiesRaw(state.Ctx)
	case "dimension":
		genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		properties, err = genericDimension.GetPropertiesRaw(state.Ctx)
	}
	if err != nil {
		internal.FatalError(err)
	}
	json := prettyJSON(properties)
	if len(json) == 0 {
		fmt.Printf("Error: No %s by id '%s'\n", entityType, entityID)
	} else {
		fmt.Println(json)
	}
}

// PrintGenericEntityLayout prints the layout of the object defined by objectID
func PrintGenericEntityLayout(state *internal.State, entityID string, entityType string) {
	var err error
	var layout json.RawMessage
	switch entityType {
	case "object":
		genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericObject.GetLayoutRaw(state.Ctx)
	case "measure":
		genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericMeasure.GetLayoutRaw(state.Ctx)
	case "dimension":
		genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericDimension.GetLayoutRaw(state.Ctx)
	}
	if err != nil {
		internal.FatalError(err)
	}
	json := prettyJSON(layout)
	if len(json) == 0 {
		fmt.Printf("Error: No %s by id '%s'\n", entityType, entityID)
	} else {
		fmt.Println(json)
	}
}

func prettyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "   ")
	return prettyJSON.String()
}

// EvalObject evalutes the data of the object identified by objectID
func EvalObject(ctx context.Context, doc *enigma.Doc, objectID string) {
	object, err := doc.GetObject(ctx, objectID)
	if err != nil {
		internal.FatalError(err)
	}

	layout, err := object.GetLayoutRaw(ctx)
	if err != nil {
		internal.FatalError(err)
	}

	layoutMap := make(map[string]interface{})
	err = json.Unmarshal(layout, &layoutMap)
	if err != nil {
		fmt.Println(err)
	}
	resultCubeMap := make(map[string]*enigma.HyperCube)
	getAllHyperCubes("", layoutMap, resultCubeMap)
	if len(resultCubeMap) == 0 {
		fmt.Printf("Object %s contains no data\n", objectID)
		return
	}
	for _, hypercube := range resultCubeMap {
		printHypercube(hypercube)
	}
}

func printHypercube(hypercube *enigma.HyperCube) {
	writer := tablewriter.NewWriter(os.Stdout)

	headers := []string{}
	for _, dim := range hypercube.DimensionInfo {
		headers = append(headers, dim.FallbackTitle)
	}
	for _, mes := range hypercube.MeasureInfo {
		headers = append(headers, mes.FallbackTitle)
	}

	writer.SetAutoFormatHeaders(false)
	writer.SetHeader(headers)

	for _, page := range hypercube.DataPages {
		for _, row := range page.Matrix {
			data := []string{}
			for _, cell := range row {
				data = append(data, cell.Text)
			}
			writer.Append(data)
		}
	}
	writer.Render()
}

func getAllHyperCubes(path string, jsonMap map[string]interface{}, resultCubeMap map[string]*enigma.HyperCube) {
	for key, value := range jsonMap {
		if key == "qHyperCube" {
			hypercube := toHypercube(value)
			resultCubeMap[key] = hypercube

		} else if !strings.HasPrefix(key, "q") {
			var customSubpath string
			if path == "" {
				customSubpath = path + "/" + key
			} else {
				customSubpath = key
			}
			switch jsonSubMap := value.(type) {
			case map[string]interface{}:
				getAllHyperCubes(customSubpath, jsonSubMap, resultCubeMap)
			}
		}
	}
}

func toHypercube(value interface{}) *enigma.HyperCube {
	subnode, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	var hypercube enigma.HyperCube
	err = json.Unmarshal(subnode, &hypercube)
	if err != nil {
		fmt.Println(err)
	}
	return &hypercube
}
