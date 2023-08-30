package args

// TODO "r" to set the refresh rate
type Arguments struct {
	TuiPort          uint64
	AlgodURL         string
	AlgodToken       string
	AlgodAdminToken  string
	AlgodDataDir     string
	AddressWatchList []string
	VersionFlag      bool
}
