package simulation

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/jackalLabs/canine-chain/x/rns/keeper"
	"github.com/jackalLabs/canine-chain/x/rns/types"
)

func SimulateMsgBid(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgBid{
			Creator: simAccount.Address.String(),
		}

		// finding a random bid
		forSale := k.GetAllForsale(ctx)
		numForSale := len(forSale)
		if numForSale < 1 {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "No domains for sale"), nil, nil
		}

		saleI := simtypes.RandIntBetween(r, 0, numForSale)
		bidDomain := forSale[saleI]

		// making the bid
		bidPrice, err := strconv.ParseFloat(bidDomain.Price, 64)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Couldn't convert bidPrice"), nil, nil
		}
		lowerBid := int(bidPrice * 0.75)
		upperBid := int(bidPrice * 1.25)
		rPrice := simtypes.RandIntBetween(r, lowerBid, upperBid)

		// building the message
		msg.Bid = fmt.Sprint(rPrice) + "ujkl"
		msg.Name = bidDomain.Name

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
