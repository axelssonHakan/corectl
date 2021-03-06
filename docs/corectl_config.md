## corectl config

With `corectl` it is possible to configure values that should be passed to your command. This can be configured in a `yaml` file residing in e.g. an repo or locally on your computer.
By default `corectl` will pick up a `corectl.yml | corectl.yaml` file from your current directory. It is also possible to pass a specific configuration file using the `--config` or `-c` flag.

All properties set in a configuration file can be overriden by passing another value as a flag instead. Properties can also be specified using environment variables. 

### Configuration properties

Below is a configuration example utilizing the different properties that are available today:

```yaml
engine: localhost:9076 
app: project1.qvf
script: script.qvs
connections:
  myconnection:
    type: testconnector
    username: gwe
    password: ${MY_PASSWORD}
    settings:
      host: corectl-test-connector
  myfolderconnection:
    connectionstring: /data
    type: folder
objects:
  - ./object-*.json
measures:
  - ./*measure*.json
dimensions:
  - ./dimension-*.json
headers:
  authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0" #generated at jwt.io with the password passw0rd
```

Note that the `password` property for the connection `myconnection` is an example use of an environment variable in the config. This can done for any property in the config file.

### engine

This property sets the URL to the engine instance that you want `corectl` to connect to by default. Can be overriden with the `-e` or `--engine` flag.

```yaml
engine: localhost:9076
```

### app

This property describes the default app name that should be used for the commands. Typically used when working against a specific app on your engine instance. If you want to open a different app than configured in the file, you can also pass the `--app` or `-a` flag.

```yaml
app: project1.qvf
```

### script

An absolute or relative path to a script file. When specifiying a script in a configuration file and running e.g. the `build` command, then `corectl` will open an app, set the script and perform a reload. Can be overriden using the `--script` flag.

```yaml
script: ./dummy-script.qvs
```

### connections

The `connections` property is an array of values and can be used to create one or many connections in your app.

Below is an example of how to define a `folder` connection in a configuration file. The `folder` connector can be used to load from or store data into files that available inside the engine container.

In this example the name of the connection will be `myfolderconnection` and be of type `folder`. The `connectionstring` is the absolute or relative path inside of the engine container.

```yaml
myfolderconnection:
    connectionstring: /data
    type: folder
```

Depending on which type of connection that should be used there may be a need for additional values in your configuration. Below is another example where we configure an connection to a custom type connection e.g. a gRPC data connector. 

When not specifying a `connectionstring` in your connection config then `corectl` builds one for you.
The value for the `type` property will be set as `provider`, and the keys and values from the `settings` property will also be appended to the `connectionstring`. The example below will yield `CUSTOM CONNECT TO \"provider=testconnector;host=corectl-test-connector;\"`.

```yaml
myconnection:
    type: testconnector
    username: gwe
    settings:
      host: corectl-test-connector
```

#### Using a separate connections file 
Optionally the `connections` property in the config file can reference a separate connections file by 
specifying a path to a yaml file instead of an inline list of connections: 
           
 ``` yaml
 connections: ./connections.yml
 ```
Connections are defined using the same format in the separate connections file:

``` yaml
connections:
 myconnection:
   type: testconnector
   username: gwe
   settings:
     host: corectl-test-connector
 myfolderconnection:
   connectionstring: /data
   type: folder
```

Note that when using the `--connections` command line option only the file reference format is supported.
 
### objects, measures and dimensions

When for example generating apps it can be useful to create objects in an app from json files that are stored remotely or locally. The properties `objects`, `measures` and `dimensions` are arrays where you can set multiple files or using wildcards for paths to files. It is also possible to use nested json structures with this approach, and then multiple objects will be created in engine.

```yaml
objects:
  - ./object-*.json
measures:
  - ./*measure*.json
dimensions:
  - ./dimension-*.json
```

### headers

When connecting a e.g. a Qlik Sense Server installation or a Qlik Associative Engine running with JWT validation enabled, there will probably be a need to pass headers with your commands. This can be done either by using the `--headers` flag or by configuring it in your configuration file.

In this example we configure an `Authorization` header with a JWT token as a value. This header will be passed with all `corectl` commands towards your engine.

```yaml
headers:
  authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0" #generated at jwt.io with the password passw0rd
```
