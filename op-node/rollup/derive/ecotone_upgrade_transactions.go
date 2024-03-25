package derive

import (
	"bytes"
	"fmt"
	"github.com/kroma-network/kroma/kroma-bindings/bindings"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum-optimism/optimism/op-service/solabi"
	"github.com/kroma-network/kroma/kroma-bindings/predeploys"
)

const UpgradeToFuncSignature = "upgradeTo(address)"

var (
	// known address w/ zero txns
	L1BlockDeployerAddress        = common.HexToAddress("0x4210000000000000000000000000000000000000")
	GasPriceOracleDeployerAddress = common.HexToAddress("0x4210000000000000000000000000000000000001")

	newL1BlockAddress        = crypto.CreateAddress(L1BlockDeployerAddress, 0)
	newGasPriceOracleAddress = crypto.CreateAddress(GasPriceOracleDeployerAddress, 0)

	deployL1BlockSource        = UpgradeDepositSource{Intent: "Ecotone: L1 Block Deployment"}
	deployGasPriceOracleSource = UpgradeDepositSource{Intent: "Ecotone: Gas Price Oracle Deployment"}
	updateL1BlockProxySource   = UpgradeDepositSource{Intent: "Ecotone: L1 Block Proxy Update"}
	updateGasPriceOracleSource = UpgradeDepositSource{Intent: "Ecotone: Gas Price Oracle Proxy Update"}
	enableEcotoneSource        = UpgradeDepositSource{Intent: "Ecotone: Gas Price Oracle Set Ecotone"}
	beaconRootsSource          = UpgradeDepositSource{Intent: "Ecotone: beacon block roots contract deployment"}

	enableEcotoneInput = crypto.Keccak256([]byte("setEcotone()"))[:4]

	EIP4788From         = common.HexToAddress("0x0B799C86a49DEeb90402691F1041aa3AF2d3C875")
	eip4788CreationData = common.FromHex("0x60618060095f395ff33373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500")
	UpgradeToFuncBytes4 = crypto.Keccak256([]byte(UpgradeToFuncSignature))[:4]

	l1BlockDeploymentBytecode        = common.FromHex(bindings.L1BlockMetaData.Bin)
	gasPriceOracleDeploymentBytecode = common.FromHex(bindings.GasPriceOracleMetaData.Bin)
)

func EcotoneNetworkUpgradeTransactions() ([]hexutil.Bytes, error) {
	upgradeTxns := make([]hexutil.Bytes, 0, 5)

	deployL1BlockTransaction, err := types.NewTx(&types.DepositTx{
		SourceHash: deployL1BlockSource.SourceHash(),
		From:       L1BlockDeployerAddress,
		To:         nil,
		Mint:       big.NewInt(0),
		Value:      big.NewInt(0),
		Gas:        500_000,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		Data: l1BlockDeploymentBytecode,
	}).MarshalBinary()

	if err != nil {
		return nil, err
	}

	upgradeTxns = append(upgradeTxns, deployL1BlockTransaction)

	deployGasPriceOracle, err := types.NewTx(&types.DepositTx{
		SourceHash: deployGasPriceOracleSource.SourceHash(),
		From:       GasPriceOracleDeployerAddress,
		To:         nil,
		Mint:       big.NewInt(0),
		Value:      big.NewInt(0),
		Gas:        1_000_000,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		Data: gasPriceOracleDeploymentBytecode,
	}).MarshalBinary()

	if err != nil {
		return nil, err
	}

	upgradeTxns = append(upgradeTxns, deployGasPriceOracle)

	updateL1BlockProxy, err := types.NewTx(&types.DepositTx{
		SourceHash: updateL1BlockProxySource.SourceHash(),
		From:       common.Address{},
		To:         &predeploys.L1BlockAddr,
		Mint:       big.NewInt(0),
		Value:      big.NewInt(0),
		Gas:        50_000,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		Data: upgradeToCalldata(newL1BlockAddress),
	}).MarshalBinary()

	if err != nil {
		return nil, err
	}

	upgradeTxns = append(upgradeTxns, updateL1BlockProxy)

	updateGasPriceOracleProxy, err := types.NewTx(&types.DepositTx{
		SourceHash: updateGasPriceOracleSource.SourceHash(),
		From:       common.Address{},
		To:         &predeploys.GasPriceOracleAddr,
		Mint:       big.NewInt(0),
		Value:      big.NewInt(0),
		Gas:        50_000,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		Data: upgradeToCalldata(newGasPriceOracleAddress),
	}).MarshalBinary()

	if err != nil {
		return nil, err
	}

	upgradeTxns = append(upgradeTxns, updateGasPriceOracleProxy)

	enableEcotone, err := types.NewTx(&types.DepositTx{
		SourceHash: enableEcotoneSource.SourceHash(),
		From:       L1InfoDepositerAddress,
		To:         &predeploys.GasPriceOracleAddr,
		Mint:       big.NewInt(0),
		Value:      big.NewInt(0),
		Gas:        80_000,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		Data: enableEcotoneInput,
	}).MarshalBinary()
	if err != nil {
		return nil, err
	}
	upgradeTxns = append(upgradeTxns, enableEcotone)

	deployEIP4788, err := types.NewTx(&types.DepositTx{
		From:  EIP4788From,
		To:    nil, // contract-deployment tx
		Mint:  big.NewInt(0),
		Value: big.NewInt(0),
		Gas:   0x3d090, // hex constant, as defined in EIP-4788
		Data:  eip4788CreationData,
		// [Kroma: START]
		// IsSystemTransaction: false,
		// [Kroma: END]
		SourceHash: beaconRootsSource.SourceHash(),
	}).MarshalBinary()

	if err != nil {
		return nil, err
	}

	upgradeTxns = append(upgradeTxns, deployEIP4788)

	return upgradeTxns, nil
}

func upgradeToCalldata(addr common.Address) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4+20))
	if err := solabi.WriteSignature(buf, UpgradeToFuncBytes4); err != nil {
		panic(fmt.Errorf("failed to write upgradeTo signature data: %w", err))
	}
	if err := solabi.WriteAddress(buf, addr); err != nil {
		panic(fmt.Errorf("failed to write upgradeTo address data: %w", err))
	}
	return buf.Bytes()
}
