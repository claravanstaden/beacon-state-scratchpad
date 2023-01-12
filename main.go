package main

import (
	"fmt"
	ssz "github.com/ferranbt/fastssz/spectests"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	dat, err := os.ReadFile("beacon_state_prysm.ssz")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("read file")

	sszState := ssz.BeaconStateBellatrix{}
	err = sszState.UnmarshalSSZ(dat)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("unmarshal done")

	tree, err := sszState.GetTree()
	if err != nil {
		return
	}

	log.Printf("get tree done")

	result, err := tree.Get(303104)
	if err != nil {
		return
	}

	log.Printf("get index at tree done")

	fmt.Println(result.Hash())
}

func getSSZFileLodestar() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://lodestar-goerli.chainsafe.io/eth/v2/debug/beacon/states/head", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create("beacon_state_lodestar.ssz")
	if err != nil {
		return err
	}

	defer out.Close()
	io.Copy(out, resp.Body)

	return nil
}

func getSSZFilePrysm() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://127.0.0.1:3500/eth/v2/debug/beacon/states/head", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create("beacon_state_prysm.ssz")
	if err != nil {
		return err
	}

	defer out.Close()
	io.Copy(out, resp.Body)

	return nil
}
