package node

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)


func setWithdrawalAddress(c *cli.Context, withdrawalAddress common.Address) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Print the "pending" disclaimer
    colorReset := "\033[0m"
    colorRed := "\033[31m"
    colorYellow := "\033[33m"
    var confirm bool
    fmt.Println("You are about to change your withdrawal address. All future ETH & RPL rewards/refunds will be sent here.")
    if !c.Bool("force") {
        confirm = false
        fmt.Println("By default, this will put your new withdrawal address into a \"pending\" state.")
        fmt.Println("Rocket Pool will continue to use your old withdrawal address until you confirm that you own the new address via the Rocket Pool website.")
        fmt.Println("You will need to use a web3-compatible wallet (such as MetaMask) with your new address to confirm it.")
        fmt.Printf("%sIf you cannot use such a wallet, or if you want to bypass this step and force Rocket Pool to use the new address immediately, please re-run this command with the \"--force\" flag.\n\n%s", colorYellow, colorReset)
    } else {
        confirm = true
        fmt.Printf("%sYou have specified the \"--force\" option, so your new address will take effect immediately.\n", colorRed)
        fmt.Printf("Please ensure that you have the correct address - you will not be able to change this once set!%s\n\n", colorReset)
    }

    // Set node's withdrawal address
    canResponse, err := rp.CanSetNodeWithdrawalAddress(withdrawalAddress, confirm)
    if err != nil {
        return err
    }

    // Prompt for a test transaction
    if cliutils.Confirm("Would you like to send a test transaction to make sure you have the correct address?") {
        inputAmount := cliutils.Prompt(fmt.Sprintf("Please enter an amount of ETH to send to %s:", withdrawalAddress), "^\\d+(\\.\\d+)?$", "Invalid amount")
        testAmount, err := strconv.ParseFloat(inputAmount, 64)
        if err != nil {
            return fmt.Errorf("Invalid test amount '%s': %w\n", inputAmount, err)
        }
        amountWei := eth.EthToWei(testAmount)
        response, err := rp.NodeSend(amountWei, "eth", withdrawalAddress)
        if err != nil {
            return err
        }

        if !cliutils.Confirm(fmt.Sprintf("Please confirm you want to send %f ETH to %s.", testAmount, withdrawalAddress)) {
            fmt.Println("Cancelled.")
            return nil
        }

        fmt.Printf("Sending ETH to %s...\n", withdrawalAddress.Hex())
        cliutils.PrintTransactionHash(rp, response.TxHash)
        if _, err = rp.WaitForTransaction(response.TxHash); err != nil {
            return err
        }

        fmt.Printf("Successfully sent the test transaction.\nPlease verify that your withdrawal address received it before confirming it below.\n\n")
    }

    // Display gas estimate
    rp.PrintGasInfo(canResponse.GasInfo)

    // Prompt for confirmation
    if !cliutils.Confirm(fmt.Sprintf("Are you sure you want to set your node's withdrawal address to %s?", withdrawalAddress.Hex())) {
        fmt.Println("Cancelled.")
        return nil
    }

    // Set node's withdrawal address
    response, err := rp.SetNodeWithdrawalAddress(withdrawalAddress, confirm)
    if err != nil {
        return err
    }

    fmt.Printf("Setting withdrawal address...\n")
    cliutils.PrintTransactionHash(rp, response.TxHash)
    if _, err = rp.WaitForTransaction(response.TxHash); err != nil {
        return err
    }

    // Log & return
    fmt.Printf("The node's withdrawal address was successfully set to %s.\n", withdrawalAddress.Hex())
    return nil

}

