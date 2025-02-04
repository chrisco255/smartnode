package node

import (
	"context"

	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
)


func getSyncProgress(c *cli.Context) (*api.NodeSyncProgressResponse, error) {

    // Get services
    if err := services.RequireNodeWallet(c); err != nil { return nil, err }
    if err := services.RequireRocketStorage(c); err != nil { return nil, err }

    // Response
    response := api.NodeSyncProgressResponse{}

    // Get eth1 client
    ec, err := services.GetEthClient(c)
    if err != nil {
        return nil, err
    }

    // Get eth1 sync progress
    progress, err := ec.SyncProgress(context.Background())
    if err != nil {
        return nil, err
    }
    if progress != nil {
        p := float64(progress.CurrentBlock - progress.StartingBlock) / float64(progress.HighestBlock - progress.StartingBlock)
        if p > 1 {
            p = 1
        } 
        response.Eth1Progress = p
        response.Eth1Synced = false
    } else {
        response.Eth1Progress = 1
        response.Eth1Synced = true
    }

    // Get eth2 client
    bc, err := services.GetBeaconClient(c)
    if err != nil {
        return nil, err
    }

    // Get eth2 sync progress
    syncStatus, err := bc.GetSyncStatus()
    if err != nil {
        return nil, err
    }
    if syncStatus.Syncing {
        response.Eth2Progress = syncStatus.Progress
        response.Eth2Synced = false
    } else {
        response.Eth2Progress = 1
        response.Eth2Synced = true
    }

    // Return response
    return &response, nil

}

