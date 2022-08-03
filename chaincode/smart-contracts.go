package main

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type DeviceInfo struct {
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Name string `json:"Name"`
	Region string `json:"Region"`
	IPFSHash string `json:"IPFSHash"`
	AuthorizedDevices []string `json:"AuthorizedDevices"`
	AuthorizedUsers []string `json:"AuthorizedUsers"`
	UpdatedAt string `json:"UpdatedAt"`
}

var count = 1

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateNewDevice(ctx contractapi.TransactionContextInterface, owner string, name string, region string) (string, error) {
	id := "dev" + strconv.Itoa(count)
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return "", err
	}
	if exists {
		return "The device " + id + " already exists", nil
	}
	
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	timeString := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).String()
	asset := DeviceInfo{
		ID: id,
		Owner: owner,
		Name: name,
		Region: region,
		IPFSHash: "",
		AuthorizedDevices: []string{},
		AuthorizedUsers: []string{},
		UpdatedAt: timeString}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}
	count++
	ctx.GetStub().PutState(id, assetJSON)
	return id, nil
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
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, name string, region string,
iPFSHash string, authorizedDevices []string, authorizedUsers []string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	timestamp, err := ctx.GetStub().GetTxTimestamp()
	timeString := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).String()
   	// overwriting original asset with new asset
	asset := DeviceInfo{
		ID: id,
		Owner: owner,
		Name: name,
		Region: region,
		IPFSHash: iPFSHash,
		AuthorizedDevices: authorizedDevices,
		AuthorizedUsers: authorizedUsers,
		UpdatedAt: timeString}
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
/*func fetchIPFSHashForDeviceFromDevice(ctx contractapi.TransactionContextInterface, requestingdeviceID, targetDeviceID ) {
	// TODO
}*/

// check if requesting user has access to given device's data
/*func (s *SmartContract) fetchIPFSHashForDeviceFromUser(ctx contractapi.TransactionContextInterface, requestinguserEmail string, targetDeviceID string) (string, error) {
	exists, err := s.AssetExists(ctx, targetDeviceID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "The device " + targetDeviceID + " does not exist", nil
	}
	
	asset, err := s.ReadAsset(ctx, targetDeviceID)
	for _, v2 := range asset.AuthorizedUsers {
		if requestinguserEmail == v2 {
			return asset.IPFSHash, nil
		}
	}
	return "The user " + requestinguserEmail + " does not have access to device " + targetDeviceID, nil
}*/
