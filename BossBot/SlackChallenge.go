package BossBot

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type myData struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

func ChallengeServer() {
	http.HandleFunc("/bossbot", func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data myData
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", data)

		rw.WriteHeader(200)
		rw.Write([]byte(data.Challenge))
	})

	http.ListenAndServe(":8332", nil)
}
