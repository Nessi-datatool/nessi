package azure

import (
	nessitypes "github.com/nessi-dev/nessi-dev/pkg/api/types"
)

func init() {
	nessitypes.GetCatalogFactory().RegisterProvider(nessitypes.AzurePurview, NewPurviewCatalog)
}
