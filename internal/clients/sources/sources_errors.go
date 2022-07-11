package sources

import (
	"errors"
)

var MoreThenOneAuthenticationForSourceErr = errors.New("more then one authentication")
var AuthenticationForSourcesNotFoundErr = errors.New("authentications for source weren't found in sources app")
var SourcesClientErr = errors.New("sources client error")
var ApplicationNotFoundErr = errors.New("application not found is sources app")
