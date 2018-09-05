package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qlik-oss/enigma-go"
)

func flattenSettings(settings map[string]string) string {
	result := ""
	for name, value := range settings {
		if result != "" {
			result += ";"
		}
		result += name + "=" + value
	}
	return result
}
func SetupConnections(ctx context.Context, doc *enigma.Doc, path string, configFilePath string) error {

	connectionConfigEntries := make(map[string]ConnectionConfigEntry)

	if configFilePath != "" {
		config := ReadConnectionsFile(configFilePath)
		for name, configEntry := range config.Connections {
			connectionConfigEntries[name] = configEntry
		}
	}
	if path != "" {
		config := ReadConnectionsFile(path)
		for name, configEntry := range config.Connections {
			connectionConfigEntries[name] = configEntry
		}
	}

	connections, err := doc.GetConnections(ctx)
	//fmt.Println("------ Setting up connections ------")

	for name, configEntry := range connectionConfigEntries {
		var connection *enigma.Connection
		if configEntry.Path != "" {
			connection = &enigma.Connection{
				Name:             name,
				Type:             "folder",
				UserName:         "",
				Password:         "",
				ConnectionString: configEntry.Path,
			}
		} else {
			connection = &enigma.Connection{
				Name:             name,
				Type:             configEntry.Type,
				UserName:         configEntry.Username,
				Password:         configEntry.Password,
				ConnectionString: "CUSTOM CONNECT TO \"provider=" + configEntry.Type + ";" + flattenSettings(configEntry.Settings) + "\"",
			}
		}
		if strings.HasPrefix(connection.Password, "${") && strings.HasSuffix(connection.Password, "}") {
			varName := strings.TrimSuffix(strings.TrimPrefix(connection.Password, "${"), "}")
			connection.Password = os.Getenv(varName)
			if connection.Password == "" {
				fmt.Println("WARNING environment variable not found:", varName)
			} else {
				LogVerbose("Resolved password from environment variable " + varName)
			}
		}

		if existingConnectionID := findExistingConnection(connections, connection.Name); existingConnectionID != "" {
			LogVerbose("Modifying connection: " + connection.Name + " (" + existingConnectionID + ")")
			err = doc.ModifyConnection(ctx, existingConnectionID, connection, true)
		} else {
			LogVerbose("Creating new connection: " + fmt.Sprint(connection))
			_, err = doc.CreateConnection(ctx, connection)
		}

		if err != nil {
			fmt.Println("Could not create/modify connection", err)
			os.Exit(1)
		}
	}
	return err
}

func findExistingConnection(connections []*enigma.Connection, name string) string {
	for _, con := range connections {
		if con.Name == name {
			return con.Id
		}
	}
	return ""
}