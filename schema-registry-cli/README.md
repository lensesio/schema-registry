Schema Registry CLI
===================

compatible <subject.version>
setcompatability [subject]
getcompatability [subject]

```
A command line interface for the Confluent schema registry

Usage:
  schemaregistrycli [command]

Available Commands:
  add         registers the schema provided through stdin
  exists      checks if the schema provided through stdin exists for the subject
  get         retrieves a schema specified by id or subject
  subjects    lists all registered subjects
  versions    lists all available versions

Flags:
  -h, --help         help for schemaregistrycli
  -e, --url string   schema registry url (default "http://localhost:8081")
  -v, --verbose      be verbose

Use "schemaregistrycli [command] --help" for more information about a command.
```
