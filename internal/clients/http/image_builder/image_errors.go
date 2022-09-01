package image_builder

import (
	"errors"
)

var ClientErr = errors.New("image builder client error")
var ComposeNotFoundErr = errors.New("image was not not found")
var ImageStatusErr = errors.New("build of requested image has not finished yet")
var BadImageTypeErr = errors.New("image type is wrong and can not be used for target platform")
