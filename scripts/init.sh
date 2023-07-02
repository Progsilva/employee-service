#!/bin/bash

/opt/mssql-tools/bin/sqlcmd -S mssql -U SA -P Password.1 -d master -i /opt/mssql_scripts/create.sql