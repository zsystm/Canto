package keeper_test

import (
	"testing"
	"time"

	"github.com/Canto-Network/Canto/v2/app"
	"github.com/Canto-Network/Canto/v2/x/csr/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/tests"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

type KeeperTestSuite struct {
	suite.Suite
	// use keeper for tests
	ctx sdk.Context

	app         *app.Canto
	consAddress sdk.ConsAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	// instantiate app
	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState())
	// initialize ctx for tests
	suite.SetupApp()
}

func (suite *KeeperTestSuite) SetupApp() {
	t := suite.T()
	// account key
	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.consAddress = sdk.ConsAddress(privCons.PubKey().Address())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "canto_9001-1",
		Time:            time.Now().UTC(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})
}

// helper function to generate test addresses for CheckPools tests
func generateSdkAddr(numAccts int) []string {
	// generate pks
	accts := make([]string, numAccts)
	for i := 0; i < numAccts; i++ {
		// fill PrivKeyField
		pk := ed25519.GenPrivKey().PubKey()
		// generate account
		accts[i] = sdk.AccAddress(pk.Bytes()).String()
	}

	return accts
}

func generateAddresses(len int) ([]string, []*types.UIntArray) {
	deployer := tests.GenerateAddress()
	contracts := make([]string, len)
	nonces := make([]*types.UIntArray, len)
	for i := 0; i < len; i++ {
		// generate nonces
		nonces[i] = &types.UIntArray{Value: []uint64{uint64(i + 1)}}
		// generate contract addresses
		contracts[i] = crypto.CreateAddress(deployer, uint64(i+1)).String()
	}
	return contracts, nonces
}