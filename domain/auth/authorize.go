package auth

import (
	"github.com/guilhermealegre/go-clean-arch-core-lib/helpers"
)

func Authorize(toAuthorize string, userAuthorizations []string, acl map[string][]string) bool {
	requiredAuthorizations, exists := acl[toAuthorize]
	if !exists {
		// toAuthorize not subjected to acl
		return true
	}

	for _, authorization := range requiredAuthorizations {
		if helpers.StringInSlice(authorization, userAuthorizations) {
			// user has the authorization required
			return true
		}
	}

	return false
}
