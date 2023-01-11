package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
)

func main() {
	start := time.Now().UTC()
	err := getSSZFilePrysm()
	if err != nil {
		log.Fatal(err)
		return
	}

	dat, err := os.ReadFile("beacon_state_prysm.ssz")
	if err != nil {
		log.Fatal(err)
		return
	}

	state := eth.BeaconStateBellatrix{}

	err = state.UnmarshalSSZ(dat)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(state.Slot)
	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)
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
