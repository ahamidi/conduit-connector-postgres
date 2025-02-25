# Conduit Connector PostgreSQL

# Source

The Postgres Source Connector connects to a database with the provided `url` and starts creating records for each change
detected in a table.

Upon starting, the source takes a snapshot of a given table in the database, then switches into CDC mode. In CDC mode,
the plugin reads from a buffer of CDC events.

## Snapshot Capture

When the connector first starts, snapshot mode is enabled. The connector acquires a read-only lock on the table, and
then reads all rows of the table into Conduit. Once all rows in that initial snapshot are read the connector releases
its lock and switches into CDC mode.

This behavior is enabled by default, but can be turned off by adding `"snapshotMode":"never"` to the Source
configuration.

## Change Data Capture

This connector implements CDC features for PostgreSQL by creating a logical replication slot and a publication that
listens to changes in the configured table. Every detected change is converted into a record and returned in the call to
`Read`. If there is no record available at the moment `Read` is called, it blocks until a record is available or the
connector receives a stop signal.

If logical replication isn't available on the PostgreSQL instance, the connector can be configured to use long polling
by adding `"cdcMode":"longPolling"` to the Source configuration.

### Logical Replication Configuration

When the connector switches to CDC mode, it attempts to run the initial setup commands to create its logical replication
slot and publication. It will connect to an existing slot if one with the configured name exists.

The Postgres user specified in the connection URL must have sufficient privileges to run all of these setup commands, or
it will fail.

Example configuration for CDC features:

```json
{
  "url": "url",
  "key": "key",
  "table": "records",
  "columns": "key,column1,column2,column3",
  "cdcMode": "logrepl",
  "logrepl.publicationName": "meroxademo",
  "logrepl.slotName": "meroxademo"
}
```

## Key Handling

If no `key` field is provided, then the connector will attempt to look up the primary key column of the table. If that
can't be determined it will fail.

## Columns

If no column names are provided in the config, then the connector will assume that all columns in the table should be
returned.

## Configuration Options

| name                      | description                                                                                                                         | required | default                |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------------------- |
| `url`                     | Connection string for the Postgres database.                                                                                        | true     |                        |
| `table`                   | The name of the table in Postgres that the connector should read.                                                                   | true     |                        |
| `columns`                 | Comma separated list of column names that should be included in each Record's payload.                                              | false    | (all columns)          |
| `key`                     | Column name that records should use for their `Key` fields.                                                                         | false    | (primary key of table) |
| `snapshotMode`            | Whether or not the plugin will take a snapshot of the entire table before starting cdc mode (allowed values: `initial` or `never`). | false    | `initial`              |
| `cdcMode`                 | Determines the CDC mode (allowed values: `auto`, `logrepl` or `long_polling`).                                                      | false    | `auto`                 |
| `logrepl.publicationName` | Name of the publication to listen for WAL events.                                                                                   | false    | `conduitpub`           |
| `logrepl.slotName`        | Name of the slot opened for replication events.                                                                                     | false    | `conduitslot`          |

# Destination

The Postgres Destination takes a `record.Record` and parses it into a valid SQL query. The Destination is designed to
handle different payloads and keys. Because of this, each record is individually parsed and upserted.

## Table Name

If a record contains a `table` property in its metadata it will be inserted in that table, otherwise it will fall back
to use the table configured in the connector. This way the Destination can support multiple tables in the same
connector, provided the user has proper access to those tables.

## Upsert Behavior

If the target table already contains a record with the same key, the Destination will upsert with its current received
values. Because Keys must be unique, this can overwrite and thus potentially lose data, so keys should be assigned
correctly from the Source.

If there is no key, the record will be simply appended.

## Configuration Options

| name    | description                                                                 | required | default |
| ------- | --------------------------------------------------------------------------- | -------- | ------- |
| `url`   | Connection string for the Postgres database.                                | true     |         |
| `table` | The name of the table in Postgres that the connector should write to.       | false    |         |
| `key`   | Column name used to detect if the target table already contains the record. | false    |         |

# Testing

Run `make test` to run all the unit and integration tests, which require Docker to be installed and running. The command
will handle starting and stopping docker containers for you.

# References

- https://github.com/bitnami/bitnami-docker-postgresql-repmgr
- https://github.com/Masterminds/squirrel
