package cmd

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"math/rand"
	"time"
)

func RandId() string {

	return xid.New().String()

}

func RandUid() string {

	limit := 8
	gid := xid.New().String()

	var xx []string
	for _, v := range gid {
		xx = append(xx, string(v))
	}

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(gid))

	var id string
	for i := 0; i < limit; i++ {
		id += xx[perm[i]]
	}

	return id

}

func RandUuid() string {

	return uuid.NewString()

}

func JsonMarshal(v interface{}) string {

	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)

}
