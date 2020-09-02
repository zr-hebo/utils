package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	// _ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	dbTypeMysql = "mysql"
)

// Host  主机
type Host struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
	Port   int    `json:"port"`
}

// UnanimityHost  id标示的主机
type UnanimityHost struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (uh *UnanimityHost) String() string {
	return fmt.Sprintf("%s:%d", uh.Host, uh.Port)
}

// UnanimityHostWithDomains   带域名的id标示的主机
type UnanimityHostWithDomains struct {
	UnanimityHost
	IP      string   `json:"ip"`
	Domains []string `json:"domains"`
}

// Field 字段
type Field struct {
	Name string
	Type string
}

// FieldType Common type include "STRING", "FLOAT", "INT", "BOOL"
func (f *Field) FieldType() string {
	return f.Type
}

// QueryRow 查询单行数据
type QueryRow struct {
	Fields []Field
	Record map[string]interface{}
}

// QueryRows 查询多行数据
type QueryRows struct {
	Fields  []Field
	Records []map[string]interface{}
}

func newQueryRow() *QueryRow {
	queryRow := new(QueryRow)
	queryRow.Fields = make([]Field, 0)
	queryRow.Record = make(map[string]interface{})
	return queryRow
}

func newQueryRows() *QueryRows {
	queryRows := new(QueryRows)
	queryRows.Fields = make([]Field, 0)
	queryRows.Records = make([]map[string]interface{}, 0)
	return queryRows
}

// MySQL Mysql主机实例
type MySQL struct {
	Host
	UserName        string
	Passwd          string
	DatabaseType    string
	DBName          string
	MultiStatements bool
	ConnectTimeout  int
	QueryTimeout    time.Duration
	connMaxLifetime time.Duration
	maxIdleConns    int
	maxOpenConns    int

	connectionLock sync.Mutex
	stmtDB         *sql.DB
}

// NewMySQL 创建MySQL数据库
func NewMySQL(
	ip string, port int, userName, passwd, dbName string) (mysql *MySQL, err error) {
	mysql = new(MySQL)
	mysql.DatabaseType = dbTypeMysql
	mysql.QueryTimeout = 30 * time.Second
	mysql.IP = ip
	mysql.Port = port
	mysql.UserName = userName
	mysql.Passwd = passwd
	mysql.DBName = dbName

	return
}

// NewMySQL 创建MySQL数据库
func NewMySQLWithTimeout(
	ip string, port int, userName, passwd, dbName string, timeout time.Duration) (mysql *MySQL, err error) {
	mysql = new(MySQL)
	mysql.DatabaseType = dbTypeMysql
	mysql.QueryTimeout = timeout
	mysql.IP = ip
	mysql.Port = port
	mysql.UserName = userName
	mysql.Passwd = passwd
	mysql.DBName = dbName

	return
}

// SetConnMaxLifetime 设置连接超时时间
func (m *MySQL) SetConnMaxLifetime(d time.Duration) {
	m.connMaxLifetime = d
	return
}

// SetMaxIdleConns 设置最大空闲连接
func (m *MySQL) SetMaxIdleConns(n int) {
	m.maxIdleConns = n
}

// SetMaxOpenConns 设置最大连接数
func (m *MySQL) SetMaxOpenConns(n int) {
	m.maxOpenConns = n
}

// Close 关闭数据库连接
func (m *MySQL) Close() (err error) {
	if m.stmtDB != nil {
		return m.stmtDB.Close()
	}
	return
}

// GetConnection 获取数据库连接
func (m *MySQL) OpenSession(ctx context.Context) (session *sql.Conn, err error) {
	m.connectionLock.Lock()
	defer m.connectionLock.Unlock()

	if m.stmtDB == nil {
		db, err := sql.Open(m.DatabaseType, m.fillConnStr())
		if err != nil {
			return nil, err
		}

		db.SetConnMaxLifetime(m.QueryTimeout)
		db.SetMaxOpenConns(m.maxOpenConns)
		db.SetMaxIdleConns(m.maxIdleConns)
		m.stmtDB = db
	}

	session, err = m.stmtDB.Conn(ctx)
	return
}

// QueryRows 执行MySQL Query语句，返回多条数据
func (m *MySQL) QueryRows(querySQL string, args ...interface{}) (queryRows *QueryRows, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("query rows on %s:%d failed <-- %s", m.IP, m.Port, err.Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	session, err := m.OpenSession(ctx)
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return nil, err
	}

	rawRows, err := session.QueryContext(ctx, querySQL, args...)
	// rawRows, err := db.Query(stmt)
	if rawRows != nil {
		defer rawRows.Close()
	}
	if err != nil {
		return
	}

	colTypes, err := rawRows.ColumnTypes()
	if err != nil {
		return
	}

	fields := make([]Field, 0, len(colTypes))
	for _, colType := range colTypes {
		fields = append(fields, Field{Name: colType.Name(), Type: getDataType(colType.DatabaseTypeName())})
	}

	queryRows = newQueryRows()
	queryRows.Fields = fields
	for rawRows.Next() {
		receiver := createReceivers(fields)
		err = rawRows.Scan(receiver...)
		if err != nil {
			err = fmt.Errorf("scan rows failed <-- %s", err.Error())
			return
		}

		if err = rawRows.Err(); err != nil {
			return
		}

		queryRows.Records = append(queryRows.Records, getRecordFromReceiver(receiver, fields))
	}
	return
}

// QueryRows 执行MySQL Query语句，返回多条数据
func QueryRowsInTx(ctx context.Context, tx *sql.Tx, querySQL string, args ...interface{}) (queryRows *QueryRows, err error) {
	rawRows, err := tx.QueryContext(ctx, querySQL, args...)
	// rawRows, err := db.Query(stmt)
	if rawRows != nil {
		defer rawRows.Close()
	}
	if err != nil {
		return
	}

	colTypes, err := rawRows.ColumnTypes()
	if err != nil {
		return
	}

	fields := make([]Field, 0, len(colTypes))
	for _, colType := range colTypes {
		fields = append(fields, Field{Name: colType.Name(), Type: getDataType(colType.DatabaseTypeName())})
	}

	queryRows = newQueryRows()
	queryRows.Fields = fields
	for rawRows.Next() {
		receiver := createReceivers(fields)
		err = rawRows.Scan(receiver...)
		if err != nil {
			err = fmt.Errorf("scan rows failed <-- %s", err.Error())
			return
		}

		if err = rawRows.Err(); err != nil {
			return
		}

		queryRows.Records = append(queryRows.Records, getRecordFromReceiver(receiver, fields))
	}
	return
}

// QueryRow 执行MySQL Query语句，返回１条或０条数据
func QueryRowInTx(ctx context.Context, tx *sql.Tx, stmt string, args ...interface{}) (row *QueryRow, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("query row failed <-- %s", err.Error())
		}
	}()

	queryRows, err := QueryRowsInTx(ctx, tx, stmt, args...)
	if err != nil || queryRows == nil {
		return
	}

	if len(queryRows.Records) < 1 {
		return
	}

	row = newQueryRow()
	row.Fields = queryRows.Fields
	row.Record = queryRows.Records[0]

	return
}

// QueryRow 执行MySQL Query语句，返回１条或０条数据
func (m *MySQL) QueryRow(stmt string, args ...interface{}) (row *QueryRow, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("query row failed <-- %s", err.Error())
		}
	}()

	queryRows, err := m.QueryRows(stmt, args...)
	if err != nil || queryRows == nil {
		return
	}

	if len(queryRows.Records) < 1 {
		return
	}

	row = newQueryRow()
	row.Fields = queryRows.Fields
	row.Record = queryRows.Records[0]

	return
}

// BatchQuery 适合返回大量数据的情况
func (m *MySQL) BatchQuery(querySQL string, args ...interface{}) (
	fields []Field, recordChan chan map[string]interface{}, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("batch query rows on %s:%d failed <-- %s", m.IP, m.Port, err.Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	session, err := m.OpenSession(ctx)
	defer func() {
		if session != nil {
			session.Close()
		}
	}()
	if err != nil {
		return
	}

	queryFieldSQL := fmt.Sprintf("%s LIMIT 1", querySQL)
	rawRows, err := session.QueryContext(ctx, queryFieldSQL, args...)
	defer func() {
		if rawRows != nil {
			rawRows.Close()
		}
	}()
	if err != nil {
		return
	}

	colTypes, err := rawRows.ColumnTypes()
	if err != nil {
		return
	}

	fields = make([]Field, 0, len(colTypes))
	for _, colType := range colTypes {
		fields = append(fields, Field{Name: colType.Name(), Type: getDataType(colType.DatabaseTypeName())})
	}

	recordChan = make(chan map[string]interface{}, 10)
	go m.fetchRowsAsync(recordChan, fields, querySQL, args...)
	return
}

func (m *MySQL) fetchRowsAsync(
	recordChan chan map[string]interface{}, fields []Field, querySQL string, args ...interface{}) {
	var err error
	defer func() {
		if err != nil {
			panic(err.Error())
		}
		close(recordChan)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	session, err := m.OpenSession(ctx)
	defer func() {
		if session != nil {
			session.Close()
		}
	}()
	if err != nil {
		return
	}

	rawRows, err := session.QueryContext(ctx, querySQL, args...)
	defer func() {
		if rawRows != nil {
			rawRows.Close()
		}
	}()
	if err != nil {
		session.Close()
		return
	}

	for {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("async query context canceled <-- %s", ctx.Err().Error())

		default:
			if rawRows.Next() {
				receiver := createReceivers(fields)
				err = rawRows.Scan(receiver...)
				if err != nil {
					panic(fmt.Sprintf("scan rows failed <-- %s", err.Error()))
				}
				if err = rawRows.Err(); err != nil {
					panic(err.Error())
				}

				recordChan <- getRecordFromReceiver(receiver, fields)

			} else {
				return
			}
		}
	}
}

func createReceivers(fields []Field) (receivers []interface{}) {
	receivers = make([]interface{}, 0, len(fields))
	for _, field := range fields {
		switch field.Type {
		case "string":
			{
				var val sql.NullString
				receivers = append(receivers, &val)
			}
		case "int32":
			{
				var val sql.NullInt32
				receivers = append(receivers, &val)
			}
		case "int64":
			{
				var val sql.NullInt64
				receivers = append(receivers, &val)
			}
		case "float32":
			{
				var val sql.NullFloat64
				receivers = append(receivers, &val)
			}
		case "float64":
			{
				var val sql.NullFloat64
				receivers = append(receivers, &val)
			}
		case "bool":
			{
				var val sql.NullBool
				receivers = append(receivers, &val)
			}
		case "blob":
			{
				var val sql.RawBytes
				receivers = append(receivers, &val)
			}
		default:
			var val sql.NullString
			receivers = append(receivers, &val)
		}
	}

	return
}

func getRecordFromReceiver(receiver []interface{}, fields []Field) (record map[string]interface{}) {
	record = make(map[string]interface{}, len(fields))
	for idx := 0; idx < len(fields); idx++ {
		field := fields[idx]
		value := receiver[idx]
		switch field.Type {
		case "string":
			{
				nullVal := value.(*sql.NullString)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.String
				}
			}
		case "int32":
			{
				nullVal := value.(*sql.NullInt32)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.Int32
				}
			}
		case "int64":
			{
				nullVal := value.(*sql.NullInt64)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.Int64
				}
			}
		case "float64":
			{
				nullVal := value.(*sql.NullFloat64)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.Float64
				}
			}
		case "float32":
			{
				nullVal := value.(*sql.NullFloat64)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = float32(nullVal.Float64)
				}
			}
		case "bool":
			{
				nullVal := value.(*sql.NullBool)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.Bool
				}
			}
		case "blob":
			{
				rawVal := value.(*sql.RawBytes)
				if rawVal != nil && *rawVal != nil {
					val := make([]byte, len(*rawVal))
					copy(val, *rawVal)
					record[field.Name] = val
				} else {
					record[field.Name] = nil
				}
			}
		default:
			{
				nullVal := value.(*sql.NullString)
				record[field.Name] = nil
				if nullVal.Valid {
					record[field.Name] = nullVal.String
				}
			}
		}
	}
	return
}

var columnTypeDict = map[string]string{
	"VARCHAR":  "string",
	"TEXT":     "string",
	"NVARCHAR": "string",
	"DATETIME": "string",
	"DECIMAL":  "string",
	"FLOAT":    "float32",
	"DOUBLE":   "float64",
	"BOOL":     "bool",
	"INT":      "int32",
	"BIGINT":   "int64",
	"BLOB":     "blob",
}

func getDataType(dbColType string) (colType string) {
	colType, ok := columnTypeDict[dbColType]
	if ok {
		return
	}

	colType = "string"
	return
}

func (m *MySQL) fillConnStr() string {
	dbServerInfoStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=%v",
		m.UserName, m.Passwd, m.IP, m.Port, m.DBName, m.MultiStatements)
	if m.ConnectTimeout > 0 {
		dbServerInfoStr = fmt.Sprintf("%s&timeout=%ds&readTimeout=%ds&writeTimeout=%ds",
			dbServerInfoStr, m.ConnectTimeout, m.QueryTimeout, m.QueryTimeout)
	}

	return dbServerInfoStr
}

// Execute 执行MySQL dml语句，返回执行结果
func (m *MySQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.stmtDB.ExecContext(context.Background(), query, args...)
}

// ExecContext 执行MySQL dml语句，返回执行结果
func (m *MySQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.stmtDB.ExecContext(ctx, query, args...)
}
