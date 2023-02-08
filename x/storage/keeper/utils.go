package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/jackalLabs/canine-chain/x/storage/types"
)

const (
	TwoGigs = 2000000000
)

// func MakeFid(data []byte) (string, error) {
//	return bech32.ConvertAndEncode(types.FidPrefix, data)
// }

func MakeCid(data []byte) (string, error) {
	return bech32.ConvertAndEncode(types.CidPrefix, data)
}

func (k Keeper) GetPaidAmount(ctx sdk.Context, address string) int64 {
	payInfo, found := k.GetStoragePaymentInfo(
		ctx,
		address,
	)
	if !found {
		return 0
	}

	return payInfo.SpaceAvailable
}

func (k Keeper) GetProviderUsing(ctx sdk.Context, provider string) int64 {
	allDeals := k.GetAllActiveDeals(ctx)

	var space int64
	for i := 0; i < len(allDeals); i++ {
		deal := allDeals[i]
		if deal.Provider != provider {
			continue
		}
		size, ok := sdk.NewIntFromString(deal.Filesize)
		if !ok {
			continue
		}

		space += size.Int64()

	}

	return space
}

// Calculate storage cost in ujkl
// Uses gigabytes and months to calculate how much user has to pay
// cost(in jkl) = ((gbs*(0.008/3))*months)/jkl_price
func (k Keeper) GetStorageCost(ctx sdk.Context, gbs int64, hours int64) sdk.Int {
	pricePerTBPerMonth := sdk.NewDec(8)
	quantifiedPricePerTBPerMonth := pricePerTBPerMonth.QuoInt64(3)
	pricePerGbPerMonth := quantifiedPricePerTBPerMonth.QuoInt64(1000)
	pricePerGbPerHour := pricePerGbPerMonth.QuoInt64(720)

	pricePerHour := pricePerGbPerHour.MulInt64(gbs)

	totalCost := pricePerHour.MulInt64(hours)

	jklPrice := k.GetJklPrice(ctx)

	// TODO: fetch denom unit from bank module
	var ujklUnit int64 = 1000000
	jklCost := totalCost.Quo(jklPrice)

	ujklCost := jklCost.MulInt64(ujklUnit)

	return ujklCost.TruncateInt()
}

// Use oracle module to get jkl price
// Returns 0.20 if feed doesn't exist or failed to unmarshal data
// Unmarshal failure is logged
func (k Keeper) GetJklPrice(ctx sdk.Context) (price sdk.Dec) {
	price = sdk.MustNewDecFromStr("0.20")

	feed, found := k.oraclekeeper.GetFeed(ctx, "jklprice")
	if found {
		type data struct {
			Price  string `json:"price"`
			Change string `json:"24h_change"`
		}

		var d data
		err := json.Unmarshal([]byte(feed.Data), &d)
		if err != nil {
			ctx.Logger().Debug("Failed to unmarshal Feed.Data (%s)", feed.Data)
		}

		p, err := sdk.NewDecFromStr(d.Price)
		if err != nil {
			ctx.Logger().Debug("Failed to convert Feed.Data.Price to sdk.Dec (%s)", d.Price)
		} else {
			price = p
		}
	}

	return price
}
