package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/qlik-oss/enigma-go"
)

func flattenSettings(settings map[string]string) string {
	result := ""
	for name, value := range settings {
		result += name + "=" + value + ";"
	}
	return result
}

// SetupConnections reads all connections from both the project file path and the config file path and updates
// the list of connections in the app.
func SetupConnections(ctx context.Context, doc *enigma.Doc, separateConnectionsFile string) error {

	connectionConfigEntries := map[string]ConnectionConfigEntry{}

	if separateConnectionsFile != "" {
		config := ReadConnectionsFile(separateConnectionsFile)
		for name, configEntry := range config.Connections {
			connectionConfigEntries[name] = configEntry
		}
	} else if ConfigDir != "" {
		config := GetConnectionsConfig()
		for name, configEntry := range config.Connections {
			connectionConfigEntries[name] = configEntry
		}
	}

	connections, err := doc.GetConnections(ctx)

	for name, configEntry := range connectionConfigEntries {
		var connection = &enigma.Connection{
			Name:     name,
			Type:     configEntry.Type,
			UserName: configEntry.Username,
			Password: configEntry.Password,
		}

		if configEntry.ConnectionString != "" {
			connection.ConnectionString = configEntry.ConnectionString
		} else {
			connection.ConnectionString = "CUSTOM CONNECT TO \"provider=" + configEntry.Type + ";" + flattenSettings(configEntry.Settings) + "\""
		}

		if existingConnectionID := findExistingConnection(connections, connection.Name); existingConnectionID != "" {
			LogVerbose("Modifying connection: " + connection.Name + " (" + existingConnectionID + ")")
			err = doc.ModifyConnection(ctx, existingConnectionID, connection, true)
		} else {
			LogVerbose("Creating new connection: " + fmt.Sprint(connection))
			var id string
			id, err = doc.CreateConnection(ctx, connection)
			if err == nil {
				fmt.Println("New connection created with id: ", id)
			}
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
