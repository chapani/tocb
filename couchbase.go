package main

import (
  "gopkg.in/couchbase/gocb.v1"
  "github.com/spf13/viper"
)


// returns couchbase bucket struct
func getBucket (vp *viper.Viper) (*gocb.Bucket, interface{}) {

  // bucket name
  bucketName := vp.GetString("couchbase.bucket")
  if bucketName == "" {
    return &gocb.Bucket{}, "Bucket name not found in config"
  }

  // host data
  host := vp.GetString("couchbase.host")
  if host == "" {
    return &gocb.Bucket{}, "Host data not found in config"
  }

  // password (optional, if bucket is not password-protected)
  password := vp.GetString("couchbase.password")

  // connect to couchbase
  cluster, connErr := gocb.Connect("couchbase://" + host)
  if connErr != nil {
    return &gocb.Bucket{}, "Connection to Couchbase failed"
  }

  return cluster.OpenBucket(bucketName, password)
} // getBucket


// inserts or updates a document into a bucket
func upsert(b *gocb.Bucket, id string, doc []uint8) (interface{}, interface{}) {
  return b.Upsert(id, doc, 0)
}
