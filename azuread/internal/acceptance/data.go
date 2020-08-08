package acceptance

import "os"

func init() {
	// unit testing
	if os.Getenv("TF_ACC") == "" {
		return
	}

	EnsureProvidersAreInitialised()
}
