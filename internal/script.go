package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/qlik-oss/enigma-go"
)

// SetScript loads the script file and sets it in the app.
func SetScript(ctx context.Context, doc *enigma.Doc, scriptFilePath string) {
	loadScript, err := ioutil.ReadFile(scriptFilePath)
	if err != nil {
		fmt.Printf("Could not find load script: %s", scriptFilePath)
		os.Exit(1)
	}

	err = doc.SetScript(ctx, string(loadScript))
}
