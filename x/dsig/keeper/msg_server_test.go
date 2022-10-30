package keeper_test

import (
	"context"
	"testing"

	keepertest "dsig/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jackalLabs/canine-chain/x/dsig/keeper"
	"github.com/jackalLabs/canine-chain/x/dsig/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.DsigKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
