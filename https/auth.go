package https

import (
	"thingworks.net/thingworks/common/datastructure/nonlinear"
)

var permissionSet = nonlinear.NewStringSet()

func AddPermission(path string) {
	permissionSet.Add(path)
}
