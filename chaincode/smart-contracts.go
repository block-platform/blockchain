package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type DeviceInfo struct {
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	IPFSHash string `json:"IPFSHash"`
	AuthorizedDevices []string `json:"AuthorizedDevices"`
	AuthorizedUsers []string `json:"AuthorizedUsers"`
}

var count = 1

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []DeviceInfo{
		{ID: "device1", Owner: "Tomoko", IPFSHash: "jgf783y4uf", 
		AuthorizedDevices: []string{"device2", "device3"}, 
		AuthorizedUsers: []string{"sandhya.shekar@sjsu.edu", "dylan.zhang@sjsu.edu"}},
		{ID: "device2", Owner: "Tomoko2", IPFSHash: "csd214h", 
		AuthorizedDevices: []string{"device2", "device3"}, 
		AuthorizedUsers: []string{"sandhya.shekar@sjsu.edu", "dylan.zhang@sjsu.edu"}},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}


// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateNewDevice(ctx contractapi.TransactionContextInterface, owner string) error {
	id := "device" + strconv.Itoa(count)
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}
	asset := DeviceInfo{
		ID: id,
		Owner: owner,
		IPFSHash: "",
		AuthorizedDevices: []string{"any"},
		AuthorizedUsers: []string{"any"}}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	count++
	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*DeviceInfo, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*DeviceInfo
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset DeviceInfo
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*DeviceInfo, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset DeviceInfo
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
// TODO there is an issue invoking this method via peer CLI due to array parameters. Looking for workaround.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, 
iPFSHash string, authorizedDevices []string, authorizedUsers []string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := DeviceInfo{
		ID: id,
		Owner: owner,
		IPFSHash: iPFSHash,
		AuthorizedDevices: authorizedDevices,
		AuthorizedUsers: authorizedUsers}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// check if requesting device has access to given device's data
// assests dataset argument that can be changed to other dataset
func fetchIPFSHashForDeviceFromDevice(requestDeviceID string, assets []DeviceInfo) string {
	// TODO Dylan
	var ipfsCode string
	if len(assets) <= 0 {
		fmt.Println("Not found any data...")
	}
	for _, v := range assets {
		for _, v2 := range v.AuthorizedDevices {
			if requestDeviceID == v2 {
				fmt.Printf("Yes, the requestDeviceID: %v has access to given devices's data (IPFSHash): %v\n", requestDeviceID, v.IPFSHash)
				ipfsCode = v.IPFSHash
			} else {
				fmt.Printf("No, the requestDeviceID: %v hasn't access to given devices's data\n", requestDeviceID)
			}
		}
	}
	return ipfsCode
}

// check if requesting user has access to given device's data
// assests dataset argument that can be changed to other dataset
func fetchIPFSHashForDeviceFromUser(requestUserEmail string, assets []DeviceInfo) []string {
	// TODO Dylan
	var authUsers []string
	if len(assets) <= 0 {
		fmt.Println("Not found any data...")
	}
	for _, v := range assets {
		for _, v2 := range v.AuthorizedUsers {
			if requestUserEmail == v2 {
				fmt.Printf("Yes, the requestUserEmail: %v has access to given devices's data (IPFSHash): %v\n", requestUserEmail, v.IPFSHash)
				authUsers = append(authUsers, v.IPFSHash)
			} else {
				fmt.Printf("No, the requestUserEmail: %v hasn't access to given devices's data\n", requestUserEmail)
			}
		}
	}
	return authUsers
}