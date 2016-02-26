package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func vaultInit(port int) string {
	type vaultReq struct {
		SecretShares    int `json:"secret_shares"`
		SecretThreshold int `json:"secret_threshold"`
	}

	var vReq vaultReq

	vault_path := "http://127.0.0.1:" + strconv.Itoa(port) + "/v1/sys/init"
	vReq.SecretShares = 5
	vReq.SecretThreshold = 3

	client := &http.Client{}

	buf, _ := json.Marshal(vReq)

	req, err := http.NewRequest("PUT", vault_path, bytes.NewBufferString(string(buf)))
	if err != nil {
		fmt.Println("Could not connect to vault service")
		os.Exit(1)
	}

	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		type initResp struct {
			Keys      []string `json:"keys"`
			RootToken string   `json:"root_token"`
		}

		var tok initResp

		body, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &tok)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		output := ""
		for index := range tok.Keys {
			output += "vault_key_" + strconv.Itoa(index) + ": " + tok.Keys[index] + "\n"
		}
		output += "vault_root_token: " + tok.RootToken

		return output
	}

	os.Exit(1)
	return ""
}

func vaultUnseal(port int, keys []string) bool {
	type unsealReq struct {
		Key string `json:"key"`
	}

	var uReq unsealReq

	vault_path := "http://127.0.0.1:" + strconv.Itoa(port) + "/v1/sys/unseal"

	client := &http.Client{}

	unsealed := false
	for index := range keys {

		uReq.Key = keys[index]
		buf, _ := json.Marshal(uReq)

		req, err := http.NewRequest("PUT", vault_path, bytes.NewBufferString(string(buf)))
		if err != nil {
			fmt.Println("Could not connect to vault service")
			os.Exit(1)
		}

		response, err := client.Do(req)
		defer response.Body.Close()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if response.StatusCode == 200 {
			type unsealResp struct {
				Sealed   bool `json:"sealed"`
				T        int  `json:"t"`
				N        int  `json:"n"`
				Progress int  `json:"progress"`
			}

			var tok unsealResp

			body, _ := ioutil.ReadAll(response.Body)
			err = json.Unmarshal(body, &tok)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !tok.Sealed {
				unsealed = true
				break
			}
		} else {
			fmt.Println("bad status", response.StatusCode)
			break
		}
	}

	return unsealed
}
