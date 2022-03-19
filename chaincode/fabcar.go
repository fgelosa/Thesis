package main

import (
  "encoding/json"
  "fmt"
  "log"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// This SmartContract provides functions for managing an Hash
   type SmartContract struct {
      contractapi.Contract
    }

// Hash describes basic details of what makes up a simple hash
   type Hash struct {
      Id             string `json:"Id"`
      HashedInfo     string `json:"hashedinfo"`
      Owner          string `json:"owner"`
    }

// InitLedger adds a base set of hashes to the ledger
   func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    hashes := []Hash{
      {Id: "-1", HashedInfo: "00000000", Owner: "Genesis"},
    }

    for _, hash := range hashes {
      hashJSON, err := json.Marshal(hash)
      if err != nil {
        return err
      }

      err = ctx.GetStub().PutState(hash.Id, hashJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

// CreateHash issues a new hash to the world state with given details.
   func (s *SmartContract) CreateHash(ctx contractapi.TransactionContextInterface, id string, hashedinfo string, owner string) error {
    exists, err := s.HashExists(ctx, id)
    if err != nil {
      return err
    }
    if exists {
      return fmt.Errorf("the hash %s already exists", id)
    }

    hash := Hash{
      Id:             id,
      HashedInfo:     hashedinfo,
      Owner:          owner,
    }
    hashJSON, err := json.Marshal(hash)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, hashJSON)
  }

// ReadHash returns the hash stored in the world state with given id.
   func (s *SmartContract) ReadHash(ctx contractapi.TransactionContextInterface, id string) (*Hash, error) {
    hashJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if hashJSON == nil {
      return nil, fmt.Errorf("the hash %s does not exist", id)
    }

    var hash Hash
    err = json.Unmarshal(hashJSON, &hash)
    if err != nil {
      return nil, err
    }

    return &hash, nil
  }

// ReadHash returns the hash stored in the world state with given id.
    func (s *SmartContract) CheckHash(ctx contractapi.TransactionContextInterface, id string) (string, error) {
     hashJSON, err := ctx.GetStub().GetState(id)
     if err != nil {
       return "", fmt.Errorf("failed to read from world state: %v", err)
     }
     if hashJSON == nil {
       return "", fmt.Errorf("the hash %s does not exist", id)
     }

     var hash Hash
     err = json.Unmarshal(hashJSON, &hash)
     if err != nil {
       return "", err
     }

     return hash.HashedInfo, nil
   }

// UpdateHash updates an existing hash in the world state with provided parameters.
   func (s *SmartContract) UpdateHash(ctx contractapi.TransactionContextInterface, id string, hashedinfo string, owner string) error {
    exists, err := s.HashExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the hash %s does not exist", id)
    }

    // overwriting original hash with new hash
    hash := Hash{
          Id:             id,
          HashedInfo:     hashedinfo,
          Owner:          owner,
    }
    hashJSON, err := json.Marshal(hash)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, hashJSON)
  }

  // DeleteHash deletes an given hash from the world state.
  func (s *SmartContract) DeleteHash(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.HashExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the hash %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// HashExists returns true when hash with given ID exists in world state
   func (s *SmartContract) HashExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    hashJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return hashJSON != nil, nil
  }

// TransferHash updates the owner field of hash with given id in world state.
   func (s *SmartContract) TransferHash(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
    hash, err := s.ReadHash(ctx, id)
    if err != nil {
      return err
    }

    hash.Owner = newOwner
    hashJSON, err := json.Marshal(hash)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, hashJSON)
  }

// GetAllHashes returns all hashes found in world state
   func (s *SmartContract) GetAllHashes(ctx contractapi.TransactionContextInterface) ([]*Hash, error) {
// range query with empty string for startKey and endKey does an
// open-ended query of all hashes in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
      return nil, err
    }
    defer resultsIterator.Close()

    var hashes []*Hash
    for resultsIterator.HasNext() {
      queryResponse, err := resultsIterator.Next()
      if err != nil {
        return nil, err
      }

      var hash Hash
      err = json.Unmarshal(queryResponse.Value, &hash)
      if err != nil {
        return nil, err
      }
      hashes = append(hashes, &hash)
    }

    return hashes, nil
  }

  func main() {
    hashChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating hash-transfer-basic chaincode: %v", err)
    }

    if err := hashChaincode.Start(); err != nil {
      log.Panicf("Error starting hash-transfer-basic chaincode: %v", err)
    }
  }
