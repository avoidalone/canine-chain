package app

import (
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/bouncybulldog"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/recovery"
	v121 "github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/121"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/alpha11"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/alpha13"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/async"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/beta6"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/beta7"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/fixstrays"
	"github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/killdeals"
	paramUpgrade "github.com/jackalLabs/canine-chain/v3/app/upgrades/testnet/params"
	v3 "github.com/jackalLabs/canine-chain/v3/app/upgrades/v3"
	v4 "github.com/jackalLabs/canine-chain/v3/app/upgrades/v4"
)

func (app *JackalApp) registerTestnetUpgradeHandlers() {
	app.registerUpgrade(alpha11.NewUpgrade(app.mm, app.configurator, app.OracleKeeper))
	app.registerUpgrade(alpha13.NewUpgrade(app.mm, app.configurator))
	app.registerUpgrade(killdeals.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(fixstrays.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(async.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(paramUpgrade.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(beta6.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(beta7.NewUpgrade(app.mm, app.configurator, app.NotificationsKeeper))
	app.registerUpgrade(v121.NewUpgrade(app.mm, app.configurator))
	app.registerUpgrade(recovery.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
}

func (app *JackalApp) registerMainnetUpgradeHandlers() {
	app.registerUpgrade(bouncybulldog.NewUpgrade(app.mm, app.configurator, app.OracleKeeper))
	app.registerUpgrade(recovery.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(v3.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
	app.registerUpgrade(v4.NewUpgrade(app.mm, app.configurator, app.StorageKeeper))
}

// registerUpgrade registers the given upgrade to be supported by the app
func (app *JackalApp) registerUpgrade(upgrade upgrades.Upgrade) {
	app.upgradeKeeper.SetUpgradeHandler(upgrade.Name(), upgrade.Handler())

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgrade.Name() && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrade.StoreUpgrades()))
	}
}
