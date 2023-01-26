package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/jackalLabs/canine-chain/x/rns/keeper"
	"github.com/jackalLabs/canine-chain/x/rns/types"
)

func SimulateMsgList(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgList{
			Creator: simAccount.Address.String(),
		}

		// choosing names that the simAccount already owns
		wctx := sdk.WrapSDKContext(ctx)
		nReq := &types.QueryListOwnedNamesRequest{
			Address: simAccount.Address.String(),
		}

		regNames, err := k.ListOwnedNames(wctx, nReq)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Couldn't fetch names"), nil, err
		}

		names := regNames.GetNames()
		if names == nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "No names to list for bidding"), nil, nil
		}

		// choosing a random name
		nameI := simtypes.RandIntBetween(r, 0, len(names))
		tName := names[nameI]

		// checking if the name is listable
		if ctx.BlockHeight() > tName.Expires {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Expired domain"), nil, nil
		}
		if tName.Locked > ctx.BlockHeight() {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Can't list a free name"), nil, nil
		}

		// generating the fees
		price := sdk.NewInt(0) // listing is free?
		spendable := bk.SpendableCoins(ctx, simAccount.Address)
		coins, hasNeg := spendable.SafeSub(sdk.NewCoins(sdk.NewCoin("ujkl", price)))

		var fees sdk.Coins

		if !hasNeg {
			var err error
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
			}
		}

		// transferring to the randomly generated simulation account
		msg.Name = tName.Name + "." + tName.Tld
		msg.Price = fmt.Sprint(simtypes.RandIntBetween(r, 0, 10000000)) + "ujkl"

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    types.ModuleName,
		}

		return simulation.GenAndDeliverTx(txCtx, fees)
	}
}
