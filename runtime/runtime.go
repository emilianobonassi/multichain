package runtime

import (
	"context"
	"fmt"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/multichain/compat/ethereumcompat"
	"github.com/renproject/multichain/compat/substratecompat"
	"github.com/renproject/pack"
)

type (
	// BitcoinCompatClients is a mapping from chains to their
	// Bitcoin-compatible clients. Clients are responsible for interacting with
	// the chain (using through an RPC interface). Chains that are not
	// Bitcoin-compatible chains will no appear in this mapping.
	BitcoinCompatClients map[multichain.Chain]bitcoincompat.Client

	// BitcoinCompatAddressDecoders is a mapping from chains to their
	// Bitcoin-compatible address decoders. Address decoders are responsible for
	// converting from strings to well-typed addresses. Chains that are not
	// Bitcoin-compatible chains will no appear in this mapping.
	BitcoinCompatAddressDecoders map[multichain.Chain]bitcoincompat.AddressDecoder

	// BitcoinCompatTxBuilders is a mapping from chains to their
	// Bitcoin-compatible transaction builders. Transaction builders are
	// responsible for building simple pay-to-address transactions, and are used
	// to build release transactions on the underlying chain. Chains that are
	// not Bitcoin-compatible chains will no appear in this mapping.
	BitcoinCompatTxBuilders map[multichain.Chain]bitcoincompat.TxBuilder

	// BitcoinCompatGasEstimators is a mapping from chains to their
	// Bitcoin-compatible gas estimators (we refer to the underlying chain fees
	// as "gas", but they are also known as miner fees). Gas estimators are
	// responsible for estimating the number of SATs/byte that a transaction
	// should pay to the underlying chain. This estimate is used in conjunction
	// with transaction building to build transactions that are highly probably
	// to be accepted by the underlying chain in a reasonable time. Chains that
	// are not Bitcoin-compatible chains will no appear in this mapping.
	BitcoinCompatGasEstimators map[multichain.Chain]bitcoincompat.GasEstimator
)

type (
	// EthereumCompatClients is a mapping from chains to their
	// Ethereum-compatible clients. Clients are responsible for interacting with
	// the chain (using through an RPC interface). Chains that are not
	// Ethereum-compatible chains will no appear in this mapping.
	EthereumCompatClients map[multichain.Chain]ethereumcompat.Client

	// EthereumCompatAddressDecoders is a mapping from chains to their
	// Ethereum-compatible address decoders. Address decoders are responsible
	// for converting from strings to well-typed addresses. Chains that are not
	// Ethereum-compatible chains will no appear in this mapping.
	EthereumCompatAddressDecoders map[multichain.Chain]ethereumcompat.AddressDecoder
)

type (
	// SubstrateCompatClients is a mapping from chains to their
	// Substrate-compatible clients. Clients are responsible for interacting
	// with the chain (using through an RPC interface). Chains that are not
	// Substrate-compatible chains will no appear in this mapping.
	SubstrateCompatClients map[multichain.Chain]substratecompat.Client

	// SubstrateCompatAddressDecoders is a mapping from chains to their
	// Substrate-compatible address decoders. Address decoders are responsible
	// for converting from strings to well-typed addresses. Chains that are not
	// Substrate-compatible chains will no appear in this mapping.
	SubstrateCompatAddressDecoders map[multichain.Chain]substratecompat.AddressDecoder
)

// The Runtime exposes all of the functionality of the underlying chains that
// are supported by the multichain. Execution engines can use this functionality
// necessary to implement cross-chain interoperability (both centralised, and
// decentralised). Often, cross-chain interoperability requires the generation
// of private keys, and the use of private keys to sign transactions. This
// functionality is explicitly excluded from the multichain runtime, allowing
// execution engines to customise these flows.
//
// The APIs exposed by the runtime are grouped by compatibility. For example,
// support for Bitcoin, Bitcoin Cash, and Zcash are all supported through the
// Bitcoin-compatibility API. The specific chain can be selected by specifying
// the "chain" argument when calling any of the BitcoinXXX methods.
// Ethereum-compatible and Substrate-compatible chains are supported through
// similar APIs. If a chain is selected that is not actually compatible with the
// API call, then an "unsupported chain" error will be returned.
//
// When new chains are added to the multichain, implementors must decide whether
// to implement one of the existing compatibility APIs, or create a new
// compatibility API. Over time, the multichain will evolve to support a wide
// enough range of compatibility APIs that most new chains will not need to
// define their own, and can simply implement and existing one.
//
// Bitcoin-compatibility API:
//
//  BitcoinDecodeAddress
//  BitcoinOutput
//  BitcoinGasPerByte
//  BitcoinBuildTx
//  BitcoinSubmitTx
//
// Ethereum-compatibility API:
//
//  EthereumDecodeAddress
//  EthereumBurnEvent
//
// Substrate-compatibility API:
//
//  SubstrateDecodeAddress
//  SubstrateBurnEvent
//
type Runtime struct {
	// Bitcoin-compatibility
	bitcoinCompatClients         BitcoinCompatClients
	bitcoinCompatAddressDecoders BitcoinCompatAddressDecoders
	bitcoinCompatTxBuilders      BitcoinCompatTxBuilders
	bitcoinCompatGasEstimators   BitcoinCompatGasEstimators
	// Ethereum-compatibility
	ethereumCompatClients         EthereumCompatClients
	ethereumCompatAddressDecoders EthereumCompatAddressDecoders
	// Substrate-compatiblity
	substrateCompatClients         SubstrateCompatClients
	substrateCompatAddressDecoders SubstrateCompatAddressDecoders
}

// NewRuntime returns a new instance of the multichain runtime. The mappings
// passed to this function define the underlying chains that are supported by
// the runtime. If a chain is not in all the mappings for its relevant
// compatibility API, then it will not be supported by the runtime.
//
// By allowing chains to be enabled/disabled through these mappings, the
// multichain can acquire support for new underlying chains as quickly as
// possible, and developers can have the flexibility to pick and choose which
// ones will be enabled for their specific use-case.
func NewRuntime(
	// Bitcoin-compatibility
	bitcoinCompatClients BitcoinCompatClients,
	bitcoinCompatAddressDecoders BitcoinCompatAddressDecoders,
	bitcoinCompatTxBuilders BitcoinCompatTxBuilders,
	bitcoinCompatGasEstimators BitcoinCompatGasEstimators,
	// Ethereum-compatibility
	ethereumCompatClients EthereumCompatClients,
	ethereumCompatAddressDecoders EthereumCompatAddressDecoders,
	// Substrate-compatiblity
	substrateCompatClients SubstrateCompatClients,
	substrateCompatAddressDecoders SubstrateCompatAddressDecoders,
) *Runtime {
	return &Runtime{
		// Bitcoin-compatibility
		bitcoinCompatClients:         bitcoinCompatClients,
		bitcoinCompatAddressDecoders: bitcoinCompatAddressDecoders,
		bitcoinCompatTxBuilders:      bitcoinCompatTxBuilders,
		bitcoinCompatGasEstimators:   bitcoinCompatGasEstimators,
		// Ethereum-compatibility
		ethereumCompatClients:         ethereumCompatClients,
		ethereumCompatAddressDecoders: ethereumCompatAddressDecoders,
		// Substrate-compatiblity
		substrateCompatClients:         substrateCompatClients,
		substrateCompatAddressDecoders: substrateCompatAddressDecoders,
	}
}

// BitcoinDecodeAddress decodes a string into a Bitcoin-compatible address.
// Address encodings are often specific to the chain. If the chain is not
// Bitcoin-compatible, then an "unsupported chain" error is returned.
func (rt *Runtime) BitcoinDecodeAddress(chain multichain.Chain, encoded pack.String) (bitcoincompat.Address, error) {
	addressDecoder, ok := rt.bitcoinCompatAddressDecoders[chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain %v", chain)
	}
	return addressDecoder.DecodeAddress(encoded)
}

// BitcoinOutput returns the Bitcoin-compatible transaction output associated
// with the given Bitcoin-compatible transaction outpoint. If the outpoint
// cannot be found, or it does not have sufficient confirmations, this method
// will return an error. If the chain is not Bitcoin-compatible, then an
// "unsupported chain" error is returned.
func (rt *Runtime) BitcoinOutput(ctx context.Context, chain multichain.Chain, asset multichain.Asset, outpoint bitcoincompat.Outpoint) (bitcoincompat.Output, error) {
	client, ok := rt.bitcoinCompatClients[chain]
	if !ok {
		return bitcoincompat.Output{}, fmt.Errorf("unsupported chain %v", chain)
	}
	// Get the tx output.
	output, confirmations, err := client.Output(ctx, outpoint)
	if err != nil {
		return bitcoincompat.Output{}, fmt.Errorf("bad output: %v", err)
	}
	// Check the tx confirmations.
	if confirmations < 1 { // TODO: This must be configurable.
		return bitcoincompat.Output{}, fmt.Errorf("insufficient confirmations: %v", confirmations)
	}
	return output, nil
}

// BitcoinGasPerByte returns the gas-per-byte that must be paid to the chain as
// a fee. This is required so that transactions do not stay pending in the
// mempool for long periods of time. If the chain is not Bitcoin-compatible,
// then an "unsupported chain" error is returned.
func (rt *Runtime) BitcoinGasPerByte(ctx context.Context, chain multichain.Chain) (pack.U64, error) {
	gasEstimator, ok := rt.bitcoinCompatGasEstimators[chain]
	if !ok {
		return pack.NewU64(0), fmt.Errorf("unsupported chain %v", chain)
	}
	return gasEstimator.GasPerByte(ctx)
}

// BitcoinBuildTx builds and returns a Bitcoin-compatible transaction that
// consumes the given transaction outputs as inputs, and produces a new set of
// transaction outputs that send funds to the given recipients. If the chain is
// not Bitcoin-compatible, then an "unsupported chain" error is returned.
func (rt *Runtime) BitcoinBuildTx(ctx context.Context, chain multichain.Chain, asset multichain.Asset, inputs []bitcoincompat.Output, recipients []bitcoincompat.Recipient) (bitcoincompat.Tx, error) {
	txBuilder, ok := rt.bitcoinCompatTxBuilders[chain]
	if !ok {
		return nil, fmt.Errorf("unsupported chain %v", chain)
	}
	return txBuilder.BuildTx(inputs, recipients)
}

// BitcoinSubmitTx will submit a signed Bitcoin-compatible transaction to the
// underlying chain, and return the transaction hash. If submission fails, an
// error is returned. If the chain is not Bitcoin-compatible, then an
// "unsupported chain" error is returned.
//
// Signing the Bitcoin-compatible transaction is not the responsibility of the
// multichain, and must be done by the execution engine. Below is an example of
// a centralised execution engine that signs the transaction using a random
// private key:
//
//  privKey := id.NewPrivKey()
//  sighashes, _ := tx.Sighashes()
//  signatures := make([]pack.Bytes65, len(sighashes))
//  for i := range sighashes {
//      hash := id.Hash(sighashes[i].Bytes32())
//      signature, _ := privKey.Sign(&hash)
//      signatures[i] = pack.NewBytes65(signature)
//  }
//  _ = tx.Sign(signatures, pack.NewBytes(wif.SerializePubKey()))
//
func (rt *Runtime) BitcoinSubmitTx(ctx context.Context, chain multichain.Chain, tx bitcoincompat.Tx) (pack.Bytes32, error) {
	client, ok := rt.bitcoinCompatClients[chain]
	if !ok {
		return pack.Bytes32{}, fmt.Errorf("unsupported chain %v", chain)
	}
	return client.SubmitTx(ctx, tx)
}

// EthereumDecodeAddress decodes a string into a Ethereum-compatible address.
// Address encodings are often specific to the chain. If the chain is not
// Ethereum-compatible, then an "unsupported chain" error is returned.
func (rt *Runtime) EthereumDecodeAddress(chain multichain.Chain, encoded pack.String) (ethereumcompat.Address, error) {
	addressDecoder, ok := rt.ethereumCompatAddressDecoders[chain]
	if !ok {
		return ethereumcompat.Address{}, fmt.Errorf("unsupported chain %v", chain)
	}
	return addressDecoder.DecodeAddress(encoded)
}

// EthereumBurnEvent returns the amount and recipient of a burn event, given the
// nonce of the burn event. If the nonce cannot be found, or the event does not
// have sufficient confirmations, this method will return an error. If the chain
// is not Ethereum-compatible, then an "unsupported chain" error is returned.
func (rt *Runtime) EthereumBurnEvent(ctx context.Context, chain multichain.Chain, asset multichain.Asset, nonce pack.Bytes32) (pack.U256, pack.String, error) {
	client, ok := rt.ethereumCompatClients[chain]
	if !ok {
		return pack.U256{}, pack.String(""), fmt.Errorf("unsupported chain %v", chain)
	}
	amount, to, confirmations, err := client.BurnEvent(ctx, asset, nonce)
	if err != nil {
		return pack.U256{}, pack.String(""), fmt.Errorf("bad burn event: %v", err)
	}
	if confirmations < 1 { // TODO: This must be configurable.
		return pack.U256{}, pack.String(""), fmt.Errorf("insufficient confirmations: %v", confirmations)
	}
	return amount, to, nil
}

// SubstrateDecodeAddress decodes a string into a Substrate-compatible address.
// Address encodings are often specific to the chain. If the chain is not
// Substrate-compatible, then an "unsupported chain" error is returned.
func (rt *Runtime) SubstrateDecodeAddress(chain multichain.Chain, encoded pack.String) (substratecompat.Address, error) {
	addressDecoder, ok := rt.substrateCompatAddressDecoders[chain]
	if !ok {
		return substratecompat.Address{}, fmt.Errorf("unsupported chain %v", chain)
	}
	return addressDecoder.DecodeAddress(encoded)
}

// SubstrateBurnEvent returns the amount and recipient of a burn event, given
// the nonce of the burn event. If the nonce cannot be found, or the event does
// not have sufficient confirmations, this method will return an error. If the
// chain is not Substrate-compatible, then an "unsupported chain" error is
// returned.
func (rt *Runtime) SubstrateBurnEvent(ctx context.Context, chain multichain.Chain, asset multichain.Asset, nonce pack.Bytes32) (pack.U256, pack.String, error) {
	client, ok := rt.substrateCompatClients[chain]
	if !ok {
		return pack.U256{}, pack.String(""), fmt.Errorf("unsupported chain %v", chain)
	}
	amount, to, confirmations, err := client.BurnEvent(ctx, asset, nonce)
	if err != nil {
		return pack.U256{}, pack.String(""), fmt.Errorf("bad burn event: %v", err)
	}
	if confirmations < 1 { // TODO: This must be configurable.
		return pack.U256{}, pack.String(""), fmt.Errorf("insufficient confirmations: %v", confirmations)
	}
	return amount, to, nil
}
