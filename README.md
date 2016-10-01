# Toggl⥃CSV: CSV-based export and import for Toggl

Toggl⥃CSV is a command line utility for Windows, Mac OS and Linux that imports or exports time reports to and from a Toggl account via CSV files

[Toggl](https://www.toggl.com) is a very convenient web-based time tracking utility that I use to track efforts for different projects during my daily work as a software developer.

Unfortunately Toggls' (as of August 2016) bulk editing, export and import are a bit frustrating to use if you have time records in multiple workspaces, across multiple years.
That's why I have created **Toggl⥃CSV**.

**Toggl⥃CSV** can be used to

- **export** toggl time records as CSV from a given account
- **import** CSV based time records into a given account

This enables you to create regular **backups** of **all** your time tracking reports, to do **bulk editing** of your time records **via Excel** and to **transfer data** from one account or workspace to another.

[![Build Status](https://travis-ci.org/andreaskoch/togglcsv.svg?branch=master)](https://travis-ci.org/andreaskoch/togglcsv)

## Usage

togglcsv `<action>` `Your-Toggl-API-Token`

### Export

Export all time records from your Toggl account starting from a given **start date** until the given **end date**:

togglcsv **export** `Your-Toggl-API-Token` `Start-Date (required, e.g. 2015-01-23)` `End-Date (optional, e.g. 2015-12-31)`

```bash
togglcsv export 1971800d4d82861d8f2c1651fea4d212 2015-01-01 2016-08-12
```

The **end date** parameter is optional. If you don't specify an end date the current date will be used.

### Import

Pipe the a given CSV file into **togglcsv** and import them into your Toggl account:

togglcsv **import** `Your-Toggl-API-Token` `<` `report.csv`

```bash
togglcsv import 1971800d4d82861d8f2c1651fea4d212 < files/toggl-report-sample.csv
```

Projects and Tags that don't exist are created automatically. But please make sure that the workspace you are assigning in your [CSV](files/toggl-report-sample.csv) does exist because workspaces cannot be created via the [Toggl API](https://github.com/toggl/toggl_api_docs).

## The CSV Format

The CSV files created by the **export** action have the following format:

| Start                | Stop                 | Workspace Name | Project Name | Client Name | Tags(s)          | Description     |
|:---------------------|:---------------------|:---------------|:-------------|:------------|:-----------------|:----------------|
| 2016-08-12T07:54:47Z | 2016-08-12T08:19:02Z | My Workspace   | Project A    | A Client    | Meetings, Sprint | Retrospective   |
| 2016-08-12T08:19:03Z | 2016-08-12T08:26:25Z | My Workspace   | Project A    | A Client    | Meetings, Sprint | Sprint Review   |
| 2016-08-12T08:26:32Z | 2016-08-12T09:34:15Z | My Workspace   | Project B    | A Client    | Meetings, Sprint | Sprint Planning |
| 2016-08-12T10:28:00Z | 2016-08-12T11:01:09Z | My Workspace   | Project C    | A Client    | Meetings, Sprint | Sprint Planning |
| 2016-08-12T11:01:09Z | 2016-08-12T13:20:32Z | My Workspace   | Project A    | A Client    | Bugs             | Fixing Bug XY   |
| ...                  | ...                  | ...            | ...          | ...         | ...              | ...             |

**CSV-file parameters**

- Header: `yes`
  - Columns: `7`
    1. Start (Date format: [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601))
    2. Stop (Date format: [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601))
    3. Workspace name (Note: The workspace must exist before the import)
    4. Project name
    5. Client name
    6. Tags (comma separated)
    7. Description of the time record
- Column Delimiter: `,`
- Row Delimiter: `\n`
- Encoding: `UTF-8`

Example: [toggl-report-sample.csv](files/toggl-report-sample.csv)

## Licensing

Toggl⥃CSV is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.

## Related Resources

### Toggl API

Toggl⥃CSV uses [github.com/andreaskoch/togglapi](https://github.com/andreaskoch/togglapi) for the communication with the [Toggl API](https://github.com/toggl/toggl_api_docs/blob/master/chapters/time_entries.md).
