package main

import(
  "encoding/json"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/spf13/viper"
  "github.com/chapani/sue"
  "gopkg.in/couchbase/gocb.v1"
)


// fetches data from mysql and adds to couchbase
func mysqlHandler (vp *viper.Viper, cbBucket *gocb.Bucket) (int, interface{}) {

  // get user name from config.toml
  user := vp.GetString("mysql.user")
  if user == "" {
    return 0, "MySQL connection failed: user is missing"
  }

  // get password
  password := vp.GetString("mysql.password")
  if password == "" {
    return 0, "MySQL connection failed: password is missing"
  }

  // get database name
  dbname := vp.GetString("mysql.dbname")
  if dbname == "" {
    return 0, "MySQL connection failed: database name is missing"
  }

  // connect to database
  db, err := sql.Open("mysql", user + ":" + password + "@/" + dbname)
  if err != nil {
    return 0, err.Error()
  }
  defer db.Close()

  // number of added documents
  count := 0

  // query.sql array from config.toml
  var queries []map[string]interface{}
  vp.UnmarshalKey("query", &queries)

  // loop over sql queries
  for i := range queries {

    query := queries[i]

    // fetch data from database
    rows, err := db.Query(query["sql"].(string))
    if err != nil {
      return 0, err.Error()
    }

    // get column names
    columns, err := rows.Columns()
    if err != nil {
      return 0, err.Error()
    }

    // make a slice for the values
    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
      scanArgs[i] = &values[i]
    }

    // fetch rows
    for rows.Next() {
      err = rows.Scan(scanArgs...)
      if err != nil {
        return 0, err.Error()
      }

      // add field values to document
      var value interface{}
      doc := map[string]interface{}{}
      for i, col := range values {
        if col == nil {
          value = "NULL"
        } else {
          value = string(col)
        }
        doc[columns[i]] = value
      }

      // add additional field values
      if props, ok := query["props"].(map[string]interface{}); ok {
        for propKey, propValue := range props {
          doc[propKey] = propValue
        }
      }

      // convert doc to json and upsert into couchbase
      jsonDoc, _ := json.Marshal(doc)
      _, err := upsert(cbBucket, sue.New2(), jsonDoc)
      if err != nil {
        return 0, err
      }

      count++
    }
    if err = rows.Err(); err != nil {
      return 0, err.Error()
    }
  }
  return count, nil
}
