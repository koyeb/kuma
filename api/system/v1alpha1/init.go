package v1alpha1

import (
	"encoding/gob"
)

func init() {
	gob.Register(&DataSource_Secret{})
}
