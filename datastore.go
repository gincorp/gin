package main

import (
    "gopkg.in/redis.v5"
)

type Datastore struct {
    db *redis.Client
}


func NewDatastore(uri string) (d Datastore, err error) {
    var opts *redis.Options

    if opts, err = redis.ParseURL(uri); err != nil {
        return
    }

    d.db = redis.NewClient(opts)

    _,err = d.db.Ping().Result()
    return
}
