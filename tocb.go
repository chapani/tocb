package main

import(
  "fmt"
  "github.com/spf13/viper"
  "os"
  "time"
)


func main () {

  // execution time started
  start := time.Now()

  // read from config.toml
  vp := viper.New()
  vp.SetConfigName("config")
  vp.AddConfigPath("./")
  err := vp.ReadInConfig()
  if err != nil {
    fmt.Println("Config file not found...", err)
    os.Exit(1)
  }

  // get bucket
  cbBucket, cbErr := getBucket(vp)
  if cbErr != nil {
    fmt.Println(cbErr)
    os.Exit(1)
  }

  // get handler type from config
  handler := vp.GetString("handler")
  if handler == "" {
    fmt.Println("Handler type is missing")
    os.Exit(1)
  }

  // run mysql handler
  if handler == "mysql" {
    count, err := mysqlHandler(vp, cbBucket)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    fmt.Printf("Number of documents added: %d\n", count)
  }

  // run sqlite handler
  if handler == "sqlite" {
    count, err := sqliteHandler(vp, cbBucket)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    fmt.Printf("Number of documents added: %d\n", count)
  }

  // print out elapsed time
  elapsed := time.Since(start)
  fmt.Printf("Time Elapsed: %f secs\n", elapsed.Seconds())
}
