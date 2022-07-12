package sources

import (
	"errors"
)

var MoreThenOneAuthenticationForSourceErr = errors.New("more then one authentication")
var AuthenticationForSourcesNotFoundErr = errors.New("authentications for source weren't found in sources app")
var SourcesClientErr = errors.New("sources client error")
var ApplicationNotFoundErr = errors.New("application not found is sources app")
var ApplicationTypesFetchUnsuccessful = errors.New("failed to fetch ApplicationTypes")
var ApplicationTypeNotFound = errors.New("application type 'provisioning' has not been found in types supported by sources")
var ApplicationTypeCacheFailed = errors.New("application type id failed to write to cache")
