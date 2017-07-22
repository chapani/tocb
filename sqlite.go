package main

import (
  "github.com/mxk/go-sqlite/sqlite3"
  "encoding/json"
  "github.com/spf13/viper"
  "github.com/chapani/sue"
  "gopkg.in/couchbase/gocb.v1"
)


// fetches data from mysql and adds to couchbase
func sqliteHandler (vp *viper.Viper, cbBucket *gocb.Bucket) (int, interface{}) {

  // get database name
  dbname := vp.GetString("sqlite.dbpath")
  if dbname == "" {
    return 0, "No database path has been provided"
  }

  // connect to database
  db, err := sqlite3.Open(dbname)
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
    doc := make(sqlite3.RowMap)

    for s, dbErr := db.Query(query["sql"].(string)); dbErr == nil; dbErr = s.Next() {
      s.Scan(doc)
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
  return count, nil
}
