package jklmining

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jackal-dao/canine/x/jklmining/keeper"
	"github.com/jackal-dao/canine/x/jklmining/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the saveRequests
	for _, elem := range genState.SaveRequestsList {
		k.SetSaveRequests(ctx, elem)
	}
	// Set all the miners
	for _, elem := range genState.MinersList {
		k.SetMiners(ctx, elem)
	}
	// Set all the mined
	for _, elem := range genState.MinedList {
		k.SetMined(ctx, elem)
	}

	// Set mined count
	k.SetMinedCount(ctx, genState.MinedCount)
	// Set all the minerClaims
	for _, elem := range genState.MinerClaimsList {
		k.SetMinerClaims(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.SaveRequestsList = k.GetAllSaveRequests(ctx)
	genesis.MinersList = k.GetAllMiners(ctx)
	genesis.MinedList = k.GetAllMined(ctx)
	genesis.MinedCount = k.GetMinedCount(ctx)
	genesis.MinerClaimsList = k.GetAllMinerClaims(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
