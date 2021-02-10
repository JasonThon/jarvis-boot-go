package https

import "sync"

var permitAllPaths []string
var once sync.Once

func PermitAll(paths ...string) {
	once.Do(func() {
		for _, path := range paths {
			permitAllPaths = append(permitAllPaths, path)
		}
	})
}
