package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Creates a new instance of the CSR object. This consumes a fully created CSRPool object.
func NewCSR(deployer string, contracts []string, csrPool *CSRPool) CSR {
	return CSR{
		Deployer:  deployer,
		Contracts: contracts,
		CsrPool:   csrPool,
	}
}

// Creates a new instance of a CSRPool. This function will look through the entire supply of NFTs and
// create a new internal representation of the NFT based on id.
func NewCSRPool(nftSupply uint64, poolAddress string) CSRPool {
	csrNFTs := []*CSRNFT{}
	id := uint64(0)
	for ; id < nftSupply; id++ {
		csrNFT := NewCSRNFT(id, poolAddress)
		csrNFTs = append(csrNFTs, &csrNFT)
	}

	return CSRPool{
		CsrNfts:     csrNFTs,
		NftSupply:   nftSupply,
		PoolAddress: poolAddress,
	}
}

// Creates a new instance of a CSRNFT. This will only be called when the CSRNFTs are initially created upon
// registration. As such, the period will default to 0 for every minted NFT.
func NewCSRNFT(id uint64, address string) CSRNFT {
	return CSRNFT{
		Period:  0,
		Id:      id,
		Address: address,
	}
}

// Validate performs stateless validation of a CSR object
func (csr CSR) Validate() error {
	// Check if the address of the deployer is valid
	deployer := csr.Deployer
	if _, err := sdk.AccAddressFromBech32(deployer); err != nil {
		return err
	}

	seenSmartContracts := make(map[string]bool)
	for _, smartContract := range csr.Contracts {
		if seenSmartContracts[smartContract] {
			return sdkerrors.Wrapf(ErrDuplicateSmartContracts, "CSR::Validate there are duplicate NFTs in this CSR.")
		}
	}

	// Ensure that there is at least one smart contract in the CSR Pool
	numSmartContracts := len(csr.Contracts)
	if numSmartContracts < 1 {
		return sdkerrors.Wrapf(ErrSmartContractSupply, "CSRPool::Validate # of smart contracts must be greater than 0 got: %d", numSmartContracts)
	}

	// Validate the CSR Pool that belongs to the
	if err := csr.CsrPool.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate performs stateless validation of a CSRPool object
func (csrPool CSRPool) Validate() error {
	// Ensure the NFT smart contract address is not empty
	if _, err := sdk.AccAddressFromBech32(csrPool.PoolAddress); err != nil {
		return err
	}

	// The total supply of NFTs must be greater than 0
	nftSupply := csrPool.NftSupply
	if nftSupply < 1 {
		return sdkerrors.Wrapf(ErrNFTSupply, "The total supply of NFTs must be greater than 0 got: %d", nftSupply)
	}

	return nil
}
