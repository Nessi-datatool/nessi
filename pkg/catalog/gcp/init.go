package gcp

import (
	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
)

func init() {
	nessitypes.GetCatalogFactory().RegisterProvider(nessitypes.GCPDataCatalog, NewDataCatalog)
}
