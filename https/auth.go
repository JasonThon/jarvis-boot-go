package https

import (
	"thingworks.net/thingworks/jarvis-boot/datastructure/nonlinear"
)

var permissionSet = nonlinear.NewStringSet()

func AddPermission(path string) {
	permissionSet.Add(path)
}
