package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/encoding/ssz/detect"
	"github.com/prysmaticlabs/prysm/v3/runtime/version"

	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	err := loadBeaconStateLocalnet()
	if err != nil {
		log.Fatal(err)
	}
}

func loadBeaconStateLocalnet() error {
	data, err := os.ReadFile("beacon_state_lodestar_localnet.ssz")
	if err != nil {
		return errors.Wrap(err, "could not open beacon state file")
	}

	localNetConfig := params.MinimalSpecConfig()

	localNetConfig.SlotsPerEpoch = 4
	localNetConfig.EpochsPerSyncCommitteePeriod = 8
	localNetConfig.SyncCommitteeSize = 32

	params.OverrideBeaconConfig(localNetConfig)

	cf, err := detect.FromState(data)
	if err != nil {
		return errors.Wrap(err, "could not sniff config+fork for origin state bytes")
	}

	config := params.BeaconConfig()

	fmt.Printf(config.ConfigName)

	_, ok := params.BeaconConfig().ForkVersionSchedule[cf.Version]
	if !ok {
		return fmt.Errorf("config mismatch, beacon node configured to connect to %s, detected state is for %s", params.BeaconConfig().ConfigName, cf.Config.ConfigName)
	}

	log.Printf("detected supported config for state & block version, config name=%s, fork name=%s", cf.Config.ConfigName, version.String(cf.Fork))
	state, err := cf.UnmarshalBeaconState(data)
	if err != nil {
		return errors.Wrap(err, "failed to initialize origin state w/ bytes + config+fork")
	}

	stateHash, err := state.HashTreeRoot(context.Background())
	if err != nil {
		return errors.Wrap(err, "unable to hash tree root")
	}

	fmt.Printf("leaf: %s\n", common.BytesToHash(state.FinalizedCheckpoint().Root))
	fmt.Printf("root: %s\n", common.BytesToHash(stateHash[:]))

	proof, err := state.FinalizedRootProof(context.Background())
	if err != nil {
		return err
	}

	for _, proofItem := range proof {
		fmt.Printf("proofitem: %s\n", common.BytesToHash(proofItem))
	}

	return nil
}

func loadBeaconStateGoerli() error {
	data, err := os.ReadFile("beacon_state_prysm.ssz")
	if err != nil {
		return errors.Wrap(err, "could not open beacon state file")
	}

	praterConfig := params.PraterConfig()

	params.OverrideBeaconConfig(praterConfig)

	cf, err := detect.FromState(data)
	if err != nil {
		return errors.Wrap(err, "could not sniff config+fork for origin state bytes")
	}

	config := params.BeaconConfig()

	fmt.Printf(config.ConfigName)

	_, ok := params.BeaconConfig().ForkVersionSchedule[cf.Version]
	if !ok {
		return fmt.Errorf("config mismatch, beacon node configured to connect to %s, detected state is for %s", params.BeaconConfig().ConfigName, cf.Config.ConfigName)
	}

	log.Printf("detected supported config for state & block version, config name=%s, fork name=%s", cf.Config.ConfigName, version.String(cf.Fork))
	state, err := cf.UnmarshalBeaconState(data)
	if err != nil {
		return errors.Wrap(err, "failed to initialize origin state w/ bytes + config+fork")
	}

	stateHash, err := state.HashTreeRoot(context.Background())
	if err != nil {
		return errors.Wrap(err, "unable to hash tree root")
	}

	fmt.Printf("leaf: %s\n", common.BytesToHash(state.FinalizedCheckpoint().Root))
	fmt.Printf("root: %s\n", common.BytesToHash(stateHash[:]))

	proof, err := state.FinalizedRootProof(context.Background())
	if err != nil {
		return err
	}

	for _, proofItem := range proof {
		fmt.Printf("proofitem: %s\n", common.BytesToHash(proofItem))
	}

	return nil
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
