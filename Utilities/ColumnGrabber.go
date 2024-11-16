package Utilities

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func ReadCSVFile(filename string, tableNames ...string) ([]Table, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(absPath, ".csv") {
		return nil, fmt.Errorf("invalid file type, only CSV files are allowed")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var tables []Table
	var currentTable Table
	var currentTableName string

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) == 0 {
			if len(currentTable.Headers) > 0 {
				if len(tableNames) == 0 || contains(tableNames, currentTableName) {
					tables = append(tables, currentTable)
				}
				currentTable = Table{}
			}
		} else {
			if len(currentTable.Headers) == 0 {
				currentTable.Headers = row
				currentTableName = row[0]
			} else {
				currentTable.Rows = append(currentTable.Rows, row)
			}
		}
	}

	if len(currentTable.Headers) > 0 {
		if len(tableNames) == 0 || contains(tableNames, currentTableName) {
			tables = append(tables, currentTable)
		}
	}

	return tables, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractColumnValues(table Table, columnName string) []string {
	columnValues := []string{}
	columnIndex := -1

	for i, header := range table.Headers {
		if header == columnName {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return columnValues
	}

	for _, row := range table.Rows {
		if len(row) > columnIndex {
			columnValues = append(columnValues, row[columnIndex])
		}
	}

	return columnValues
}

func ColumnGrabber() {
	var tableNames []string
	filename := CLI_Handlers.GetFilePath()

	Interface.Write("Enter the column name to grab values from")
	Interface.Input()
	var columnName string
	_, err := fmt.Scanln(&columnName)
	CLI_Handlers.LogError(err)

	if len(columnName) == 0 {
		fmt.Println("No column name provided.")
		return
	}
	tables, err := ReadCSVFile(filename, tableNames...)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, table := range tables {
		columnValues := extractColumnValues(table, columnName)
		for _, value := range columnValues {
			fmt.Println(value)
			err := CLI_Handlers.AppendToFile(Helpers.OutputPath+"/"+columnName+".txt", []string{value})
			CLI_Handlers.LogError(err)
		}
	}

	Interface.Write("Press enter to go Main Menu")
	_, err = fmt.Scanln()
	CLI_Handlers.LogError(err)
}
