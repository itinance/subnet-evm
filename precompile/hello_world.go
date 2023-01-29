// Code generated
// This file is a generated precompile contract with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

// There are some must-be-done changes waiting in the file. Each area requiring you to add your code is marked with CUSTOM CODE to make them easy to find and modify.
// Additionally there are other files you need to edit to activate your precompile.
// These areas are highlighted with comments "ADD YOUR PRECOMPILE HERE".
// For testing take a look at other precompile tests in core/stateful_precompile_test.go

/* General guidelines for precompile development:
1- Read the comment and set a suitable contract address in precompile/params.go. E.g:
	HelloWorldAddress = common.HexToAddress("ASUITABLEHEXADDRESS")
2- Set gas costs here
3- It is recommended to only modify code in the highlighted areas marked with "CUSTOM CODE STARTS HERE". Modifying code outside of these areas should be done with caution and with a deep understanding of how these changes may impact the EVM.
Typically, custom codes are required in only those areas.
4- Add your upgradable config in params/precompile_config.go
5- Add your precompile upgrade in params/config.go
6- Add your solidity interface and test contract to contract-examples/contracts
7- Write solidity tests for your precompile in contract-examples/test
8- Create your genesis with your precompile enabled in tests/e2e/genesis/
9- Create e2e test for your solidity test in tests/e2e/solidity/suites.go
10- Run your e2e precompile Solidity tests with './scripts/run_ginkgo.sh'

*/

package precompile

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/vmerrs"

	"github.com/ethereum/go-ethereum/common"
)

const (
	SayHelloGasCost    uint64 = readGasCostPerSlot  // SET A GAS COST HERE
	SetGreetingGasCost uint64 = writeGasCostPerSlot // SET A GAS COST HERE

	// HelloWorldRawABI contains the raw ABI of HelloWorld contract.
	HelloWorldRawABI = "[{\"inputs\":[],\"name\":\"sayHello\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"name\":\"setGreeting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	_ StatefulPrecompileConfig = &HelloWorldConfig{}

	HelloWorldABI abi.ABI // will be initialized by init function

	HelloWorldPrecompile StatefulPrecompiledContract // will be initialized by init function

)

// HelloWorldConfig implements the StatefulPrecompileConfig
// interface while adding in the HelloWorld specific precompile address.
type HelloWorldConfig struct {
	UpgradeableConfig
}

func init() {
	parsed, err := abi.JSON(strings.NewReader(HelloWorldRawABI))
	if err != nil {
		panic(err)
	}
	HelloWorldABI = parsed

	HelloWorldPrecompile = createHelloWorldPrecompile(HelloWorldAddress)
}

// NewHelloWorldConfig returns a config for a network upgrade at [blockTimestamp] that enables
// HelloWorld .
func NewHelloWorldConfig(blockTimestamp *big.Int) *HelloWorldConfig {
	return &HelloWorldConfig{

		UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableHelloWorldConfig returns config for a network upgrade at [blockTimestamp]
// that disables HelloWorld.
func NewDisableHelloWorldConfig(blockTimestamp *big.Int) *HelloWorldConfig {
	return &HelloWorldConfig{
		UpgradeableConfig: UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Equal returns true if [s] is a [*HelloWorldConfig] and it has been configured identical to [c].
func (c *HelloWorldConfig) Equal(s StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*HelloWorldConfig)
	if !ok {
		return false
	}
	// CUSTOM CODE STARTS HERE
	// modify this boolean accordingly with your custom HelloWorldConfig, to check if [other] and the current [c] are equal
	// if HelloWorldConfig contains only UpgradeableConfig  you can skip modifying it.
	equals := c.UpgradeableConfig.Equal(&other.UpgradeableConfig)
	return equals
}

// String returns a string representation of the HelloWorldConfig.
func (c *HelloWorldConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

// Address returns the address of the HelloWorld. Addresses reside under the precompile/params.go
// Select a non-conflicting address and set it in the params.go.
func (c *HelloWorldConfig) Address() common.Address {
	return HelloWorldAddress
}

// Configure configures [state] with the initial configuration.
func (c *HelloWorldConfig) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	// This will be called in the first block where HelloWorld stateful precompile is enabled.
	// 1) If BlockTimestamp is nil, this will not be called
	// 2) If BlockTimestamp is 0, this will be called while setting up the genesis block
	// 3) If BlockTimestamp is 1000, this will be called while processing the first block
	// whose timestamp is >= 1000
	//
	// Set the initial value under [common.BytesToHash([]byte("storageKey")] to "Hello World!"
	res := common.LeftPadBytes([]byte("Hello World!"), common.HashLength)
	state.SetState(HelloWorldAddress, common.BytesToHash([]byte("storageKey")), common.BytesToHash(res))
}

// Contract returns the singleton stateful precompiled contract to be used for HelloWorld.
func (c *HelloWorldConfig) Contract() StatefulPrecompiledContract {
	return HelloWorldPrecompile
}

// Verify tries to verify HelloWorldConfig and returns an error accordingly.
func (c *HelloWorldConfig) Verify() error {

	// CUSTOM CODE STARTS HERE
	// Add your own custom verify code for HelloWorldConfig here
	// and return an error accordingly
	return nil
}

// PackSayHello packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackSayHello() ([]byte, error) {
	return HelloWorldABI.Pack("sayHello")
}

// PackSayHelloOutput attempts to pack given result of type string
// to conform the ABI outputs.
func PackSayHelloOutput(result string) ([]byte, error) {
	return HelloWorldABI.PackOutput("sayHello", result)
}

func sayHello(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, SayHelloGasCost); err != nil {
		return nil, 0, err
	}

	// Get the current state
	currentState := accessibleState.GetStateDB()
	// Get the value set at recipient
	value := currentState.GetState(HelloWorldAddress, common.BytesToHash([]byte("storageKey")))
	packedOutput, err := PackSayHelloOutput(string(common.TrimLeftZeroes(value.Bytes())))
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackSetGreetingInput attempts to unpack [input] into the string type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSetGreetingInput(input []byte) (string, error) {
	res, err := HelloWorldABI.UnpackInput("setGreeting", input)
	if err != nil {
		return "", err
	}
	unpacked := *abi.ConvertType(res[0], new(string)).(*string)
	return unpacked, nil
}

// PackSetGreeting packs [response] of type string into the appropriate arguments for setGreeting.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackSetGreeting(response string) ([]byte, error) {
	return HelloWorldABI.Pack("setGreeting", response)
}

func setGreeting(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, SetGreetingGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// Attempts to unpack [input] into the arguments to the SetGreetingInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStr, err := UnpackSetGreetingInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	// Check if the input string is longer than 32 bytes
	if len(inputStr) > 32 {
		return nil, 0, errors.New("input string is longer than 32 bytes")
	}

	// setGreeting is the execution function
	// "SetGreeting(name string)" and sets the storageKey
	// in the string returned by hello world
	res := common.LeftPadBytes([]byte(inputStr), common.HashLength)
	accessibleState.GetStateDB().SetState(HelloWorldAddress, common.BytesToHash([]byte("storageKey")), common.BytesToHash(res))

	// This function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createHelloWorldPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.

func createHelloWorldPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	var functions []*statefulPrecompileFunction

	methodSayHello, ok := HelloWorldABI.Methods["sayHello"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodSayHello.ID, sayHello))

	methodSetGreeting, ok := HelloWorldABI.Methods["setGreeting"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodSetGreeting.ID, setGreeting))

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, functions)
	return contract
}
