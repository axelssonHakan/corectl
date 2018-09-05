package internal

import (
	"context"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"strings"

	"github.com/qlik-oss/enigma-go"
)

type State struct {
	Doc     *enigma.Doc
	Ctx     context.Context
	Global  *enigma.Global
	AppID   string
	MetaURL string
	Verbose bool
}

func PrepareEngineState(ctx context.Context, engine string, sessionID string, appID string, ttl string) *State {
	LogVerbose("---------- Connecting to app ----------")

	engineURL := buildWebSocketURL(engine, ttl)
	var doc *enigma.Doc

	LogVerbose("Engine: " + engineURL)

	LogVerbose("SessionId: " + sessionID)
	headers := make(http.Header, 1)
	if sessionID != "" {
		headers.Set("X-Qlik-Session", sessionID)
	}
	global, err := enigma.Dialer{}.Dial(ctx, engineURL, headers)
	if err != nil {
		fmt.Println("Could not connect to engine:"+engine, err)
		os.Exit(-1)
	}

	go func() {
		for x := range global.SessionMessageChannel() {
			if x.Topic != "OnConnected" {
				fmt.Println(x.Topic, string(x.Content))
			}
		}
	}()
	if sessionID != "" {
		doc, err = global.GetActiveDoc(ctx)
		if doc != nil {
			appID := "SessionDoc"
			LogVerbose("Document: " + appID + "(reconnected)")
		} else {
			LogVerbose("No active doc: " + appID + ", session=" + sessionID)
		}
	}
	if doc == nil {
		if appID == "" {
			doc, err = global.CreateSessionApp(ctx)
			if doc != nil {
				appID = doc.GenericId
				LogVerbose("Document: " + appID + "(new session app)")
			} else {
				fmt.Println(err)
			}
		} else {
			if doc == nil {
				doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
				if doc != nil {
					LogVerbose("Document: " + appID + "(opened)")
				}
			}
			if doc == nil {
				_, _, err = global.CreateApp(ctx, appID, "")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if doc != nil {
					LogVerbose("Document: " + appID + "(new)")
				}
			}
		}
	}

	metaURL := buildMetadataURL(engine, appID)
	LogVerbose("Meta: " + metaURL)

	return &State{
		Doc:     doc,
		Global:  global,
		AppID:   appID,
		Ctx:     ctx,
		MetaURL: metaURL,
	}
}

func tidyUpEngine(engine string) string {
	var url string
	if strings.HasPrefix(engine, "wss://") {
		url = engine
	} else if strings.HasPrefix(engine, "ws://") {
		url = engine
	} else {
		url = "ws://" + engine
	}
	if len(strings.Split(url, ":")) == 2 {
		url += ":9076"
	}
	return url
}

func buildWebSocketURL(engine string, ttl string) string {
	engine = tidyUpEngine(engine)
	return engine + "/app/engineData/ttl/" + ttl
}

func buildMetadataURL(engine string, appID string) string {
	engine = tidyUpEngine(engine)
	engine = strings.Replace(engine, "wss://", "https://", -1)
	engine = strings.Replace(engine, "ws://", "http://", -1)
	url := fmt.Sprintf("%s/v1/apps/%s/data/metadata", engine, neturl.QueryEscape(appID))
	return url
}