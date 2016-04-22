package modules

import (
	"github.com/NebulousLabs/Sia/types"
)

const (
	// HostDir names the directory that contains the host persistence.
	HostDir = "host"
)

type (
	// HostFinancialMetrics provides financial statistics for the host,
	// including money that is locked in contracts. Though verbose, these
	// statistics should provide a clear picture of where the host's money is
	// currently being used. The front end can consolidate stats where desired.
	// Potential revenue refers to revenue that is available in a file
	// contract, but the file contract window has not yet closed.
	HostFinancialMetrics struct {
		// Every time a renter forms a contract with a host, a contract fee is
		// paid by the renter. These stats track the total contract fees.
		ContractCompensation          types.Currency `json:"contractcompensation"`
		PotentialContractCompensation types.Currency `json:"potentialcontractcompensation"`

		// Metrics related to storage proofs, collateral, and submitting
		// transactions to the blockchain.
		LockedStorageCollateral types.Currency `json:"lockedstoragecollateral"`
		LostRevenue             types.Currency `json:"lostrevenue"`
		LostStorageCollateral   types.Currency `json:"loststoragecollateral"`
		PotentialStorageRevenue types.Currency `json:"potentialerevenue"`
		RiskedStorageCollateral types.Currency `json:"riskedstoragecollateral"`
		StorageRevenue          types.Currency `json:"storagerevenue"`
		TransactionFeeExpenses  types.Currency `json:"transactionfeeexpenses"`

		// Bandwidth financial metrics.
		DownloadBandwidthRevenue          types.Currency `json:"downloadbandwidthrevenue"`
		PotentialDownloadBandwidthRevenue types.Currency `json:"potentialdownloadbandwidthrevenue"`
		PotentialUploadBandwidthRevenue   types.Currency `json:"potentialuploadbandwidthrevenue"`
		UploadBandwidthRevenue            types.Currency `json:"uploadbandwidthrevenue"`
	}

	// HostInternalSettings contains a list of settings that can be changed.
	HostInternalSettings struct {
		AcceptingContracts   bool              `json:"acceptingcontracts"`
		MaxDuration          types.BlockHeight `json:"maxduration"`
		MaxDownloadBatchSize uint64            `json:"maxdownloadbatchsize"`
		MaxReviseBatchSize   uint64            `json:"maxrevisebatchsize"`
		NetAddress           NetAddress        `json:"netaddress"`
		WindowSize           types.BlockHeight `json:"windowsize"`

		Collateral            types.Currency `json:"collateral"`
		CollateralBudget      types.Currency `json:"collateralbudget"`
		MaxCollateralFraction types.Currency `json:"maxcollateralfraction"`
		MaxCollateral         types.Currency `json:"maxcollateral"`

		DownloadLimitGrowth uint64 `json:"downloadlimitgrowth"` // Bytes per second that get added to the limit for how much download bandwidth the host is allowed to use.
		DownloadLimitCap    uint64 `json:"downloadlimitcap"`    // The maximum size of the limit for how much download bandwidth the host is allowed to use.
		DownloadSpeedLimit  uint64 `json:"downloadspeedlimit"`  // The maximum download speed for all combined host connections.
		UploadLimitGrowth   uint64 `json:"uploadlimitgrowth"`   // Bytes per second that get added to the limit for how much upload bandwidth the host is allowed to use.
		UploadLimitCap      uint64 `json:"uploadlimitcap"`      // The maximum size of the limit for how much upload bandwidth the host is allowed to use.
		UploadSpeedLimit    uint64 `json:"uploadspeedlimit"`    // The maximum upload speed for all combined host connections.

		MinimumContractPrice          types.Currency `json:"contractprice"`
		MinimumDownloadBandwidthPrice types.Currency `json:"minimumdownloadbandwidthprice"`
		MinimumStoragePrice           types.Currency `json:"storageprice"`
		MinimumUploadBandwidthPrice   types.Currency `json:"minimumuploadbandwidthprice"`
	}

	// HostNetworkMetrics reports the quantity of each type of RPC call that
	// has been made to the host.
	HostNetworkMetrics struct {
		NetAddress NetAddress

		DownloadBandwidthConsumed uint64 `json:"downloadbandwidthconsumed"`
		UploadBandwidthConsumed   uint64 `json:"uploadbandwidthconsumed"`

		DownloadCalls     uint64 `json:"downloadcalls"`
		ErrorCalls        uint64 `json:"errorcalls"`
		FormContractCalls uint64 `json:"formcontractcalls"`
		RenewCalls        uint64 `json:"renewcalls"`
		ReviseCalls       uint64 `json:"revisecalls"`
		SettingsCalls     uint64 `json:"settingscalls"`
		UnrecognizedCalls uint64 `json:"unrecognizedcalls"`
	}

	// A Host can take storage from disk and offer it to the network, managing
	// things such as announcements, settings, and implementing all of the RPCs
	// of the host protocol.
	Host interface {
		// Announce submits a host announcement to the blockchain.
		Announce() error

		// AnnounceAddress submits an announcement using the given address.
		AnnounceAddress(NetAddress) error

		// FinancialMetrics returns the financial statistics of the host.
		FinancialMetrics() HostFinancialMetrics

		// InternalSettings returns the host's internal settings.
		InternalSettings() HostInternalSettings

		// NetworkMetrics returns information on the types of RPC calls that
		// have been made to the host.
		NetworkMetrics() HostNetworkMetrics

		// SetInternalSettings sets the hosting parameters of the host.
		SetInternalSettings(HostInternalSettings) error

		// The storage manager provides an interface for adding and removing
		// storage folders and data sectors to the host.
		StorageManager
	}
)

// BandwidthPriceToConsensus converts a human bandwidth price, having the unit
// 'Siacoins per Terabyte', to a consensus storage price, having the unit
// 'Hastings per Byte'.
func BandwidthPriceToConsensus(siacoinsTB uint64) (hastingsByte types.Currency) {
	hastingsTB := types.NewCurrency64(siacoinsTB).Mul(types.SiacoinPrecision)
	return hastingsTB.Div(types.NewCurrency64(1e12))
}

// BandwidthPriceToHuman converts a consensus bandwidth price, having the unit
// 'Hastings per Byte' to a human bandwidth price, having the unit 'Siacoins
// per Terabyte'.
func BandwidthPriceToHuman(hastingsByte types.Currency) (siacoinsTB uint64, err error) {
	hastingsTB := hastingsByte.Mul(types.NewCurrency64(1e12))
	if hastingsTB.Cmp(types.SiacoinPrecision.Div(types.NewCurrency64(2))) < 0 {
		// The result of the final division is going to be less than 0.5,
		// therefore 0 should be returned.
		return 0, nil
	}
	if hastingsTB.Cmp(types.SiacoinPrecision) < 0 {
		// The result of the final division is going to be greater than or
		// equal to 0.5, but less than 1, therefore 1 should be returned.
		return 1, nil
	}
	return hastingsTB.Div(types.SiacoinPrecision).Uint64()
}

// StoragePriceToConsensus converts a human storage price, having the unit
// 'Siacoins per Month per Terabyte', to a consensus storage price, having the
// unit 'Hastings per Block per Byte'.
func StoragePriceToConsensus(siacoinsMonthTB uint64) (hastingsBlockByte types.Currency) {
	// Perform multiplication first to preserve precision.
	hastingsMonthTB := types.NewCurrency64(siacoinsMonthTB).Mul(types.SiacoinPrecision)
	hastingsBlockTB := hastingsMonthTB.Div(types.NewCurrency64(4320))
	return hastingsBlockTB.Div(types.NewCurrency64(1e12))
}

// StoragePriceToHuman converts a consensus storage price, having the unit
// 'Hastings per Block per Byte', to a human storage price, having the unit
// 'Siacoins per Month per Terabyte'. An error is returned if the result would
// overflow a uint64. If the result is between 0 and 1, the value is rounded to
// the nearest value.
func StoragePriceToHuman(hastingsBlockByte types.Currency) (siacoinsMonthTB uint64, err error) {
	// Perform multiplication first to preserve precision.
	hastingsMonthByte := hastingsBlockByte.Mul(types.NewCurrency64(4320))
	hastingsMonthTB := hastingsMonthByte.Mul(types.NewCurrency64(1e12))
	if hastingsMonthTB.Cmp(types.SiacoinPrecision.Div(types.NewCurrency64(2))) < 0 {
		// The result of the final division is going to be less than 0.5,
		// therefore 0 should be returned.
		return 0, nil
	}
	if hastingsMonthTB.Cmp(types.SiacoinPrecision) < 0 {
		// The result of the final division is going to be greater than or
		// equal to 0.5, but less than 1, therefore 1 should be returned.
		return 1, nil
	}
	return hastingsMonthTB.Div(types.SiacoinPrecision).Uint64()
}
