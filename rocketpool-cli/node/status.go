package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)


func getStatus(c *cli.Context) error {

    // Get RP client
    rp, err := rocketpool.NewClientFromCtx(c)
    if err != nil { return err }
    defer rp.Close()

    // Get node status
    status, err := rp.NodeStatus()
    if err != nil {
        return err
    }

    // Account address & balances
    fmt.Printf(
        "The node %s has a balance of %.6f ETH and %.6f RPL.\n",
        status.AccountAddress.Hex(),
        math.RoundDown(eth.WeiToEth(status.AccountBalances.ETH), 6),
        math.RoundDown(eth.WeiToEth(status.AccountBalances.RPL), 6))
    if status.AccountBalances.FixedSupplyRPL.Cmp(big.NewInt(0)) > 0 {
        fmt.Printf("The node has a balance of %.6f old RPL which can be swapped for new RPL.\n", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyRPL), 6))
    }

    // Registered node details
    if status.Registered {

        // Withdrawal address & balances
        if !bytes.Equal(status.AccountAddress.Bytes(), status.WithdrawalAddress.Bytes()) {
            fmt.Printf(
                "The node's withdrawal address %s has a balance of %.6f ETH and %.6f RPL.\n",
                status.WithdrawalAddress.Hex(),
                math.RoundDown(eth.WeiToEth(status.WithdrawalBalances.ETH), 6),
                math.RoundDown(eth.WeiToEth(status.WithdrawalBalances.RPL), 6))
        }
        fmt.Println("")

        // Node status
        fmt.Printf("The node is registered with Rocket Pool with a timezone location of %s.\n", status.TimezoneLocation)
        if status.Trusted {
            fmt.Println("The node is a member of the oracle DAO - it can create unbonded minipools, vote on DAO proposals and perform watchtower duties.")
        }
        fmt.Println("")

        // RPL stake details
        fmt.Printf(
            "The node has a total stake of %.6f RPL and an effective stake of %.6f RPL, allowing it to run %d minipool(s) in total.\n",
            math.RoundDown(eth.WeiToEth(status.RplStake), 6),
            math.RoundDown(eth.WeiToEth(status.EffectiveRplStake), 6),
            status.MinipoolLimit)
        fmt.Printf(
            "This is currently a %.2f%% collateral ratio.\n",
            status.CollateralRatio * 100,
        )

        // Minipool details
        if status.MinipoolCounts.Total > 0 {

            // RPL stake
            fmt.Printf("The node must keep at least %.6f RPL staked to collateralize its minipools and claim RPL rewards.\n", math.RoundDown(eth.WeiToEth(status.MinimumRplStake), 6))
            fmt.Println("")

            // Minipools
            fmt.Printf("The node has a total of %d minipool(s):\n", status.MinipoolCounts.Total)
            if status.MinipoolCounts.Initialized > 0 {
                fmt.Printf("- %d initialized\n", status.MinipoolCounts.Initialized)
            }
            if status.MinipoolCounts.Prelaunch > 0 {
                fmt.Printf("- %d at prelaunch\n", status.MinipoolCounts.Prelaunch)
            }
            if status.MinipoolCounts.Staking > 0 {
                fmt.Printf("- %d staking\n", status.MinipoolCounts.Staking)
            }
            if status.MinipoolCounts.Withdrawable > 0 {
                fmt.Printf("- %d withdrawable (after withdrawal delay)\n", status.MinipoolCounts.Withdrawable)
            }
            if status.MinipoolCounts.Dissolved > 0 {
                fmt.Printf("- %d dissolved\n", status.MinipoolCounts.Dissolved)
            }
            if status.MinipoolCounts.RefundAvailable > 0 {
                fmt.Printf("* %d minipool(s) have refunds available!\n", status.MinipoolCounts.RefundAvailable)
            }
            if status.MinipoolCounts.WithdrawalAvailable > 0 {
                fmt.Printf("* %d minipool(s) are ready for withdrawal!\n", status.MinipoolCounts.WithdrawalAvailable)
            }
            if status.MinipoolCounts.CloseAvailable > 0 {
                fmt.Printf("* %d dissolved minipool(s) can be closed!\n", status.MinipoolCounts.CloseAvailable)
            }

        } else {
            fmt.Println("The node does not have any minipools yet.")
        }
        
    } else {
        fmt.Println("The node is not registered with Rocket Pool.")
    }

    // Return
    return nil

}

