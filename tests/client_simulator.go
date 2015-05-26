package tests

import (
  "fmt"
  "bytes"
  "io/ioutil"
  "net/http"
)

func TestServer(server string) {
  fmt.Println("URL:>", server)

  var jsonStr = []byte(`
    "method": "get_all_builds",
    "params": {
      "device": "hammerhead",
      "channels": [
          "stable",
          "snapshot",
          "RC",
          "nightly",
      ]
    }
  `)
  req, err := http.NewRequest("POST", server, bytes.NewBuffer(jsonStr))
  req.Header.Set("User-Agent", "com.copperhead.updater/0.1")
  req.Header.Set("Cache-control", "no-cache")
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  fmt.Println("response Headers:", resp.Header)
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
  // 'method' : 'get_all_builds',
  // 'params' : {
  //     'device' : 'hammerhead',
  //     'channels': [
  //         'stable',
  //         'snapshot',
  //         'RC',
  //         'nightly'
  //     ],
  //     // Optional: use this to get always the newest zips based on the current one.
  //     // If not used: will get all the zips.
  //     //'source_incremental' : ''
  // }
  // 'Cache-control' : 'no-cache',
  // 'Content-type' : 'application/json',
  // 'User-Agent' : 'com.cyanogenmod.updater/2.2'
}
