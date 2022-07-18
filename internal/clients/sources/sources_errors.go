package sources

import (
	"encoding/json"
	"errors"
)

var MoreThenOneAuthenticationForSourceErr = errors.New("more then one authentication")
var AuthenticationForSourcesNotFoundErr = errors.New("authentications for source weren't found in sources app")
var SourcesClientErr = errors.New("sources client error")
var ApplicationNotFoundErr = errors.New("application not found is sources app")
var ApplicationTypesFetchUnsuccessfulErr = errors.New("failed to fetch ApplicationTypes")
var ApplicationTypeNotFoundErr = errors.New("application type 'provisioning' has not been found in types supported by sources")
var ApplicationTypeCacheFailedErr = errors.New("application type id failed to write to cache")
var CantMarshalErr = errors.New("cant marshal array to json")

func ParseErrorNotFoundToJSON(sourcesErrors ErrorNotFound) ([]byte, error) {
	j, err := json.Marshal(sourcesErrors)
	if err != nil {
		return nil, CantMarshalErr
	}
	return j, nil
}

func ParseErrorBadRequestToJSON(sourcesErrors ErrorBadRequest) ([]byte, error) {
	j, err := json.Marshal(sourcesErrors)
	if err != nil {
		return nil, CantMarshalErr
	}
	return j, nil
}
