/*
Package usrerr provides error type and predefined values for user-facing HTTP errors.

Since our application makes heavy use of Go error wrapping, many errors end up very
long in the form of err: err: err: err: something happened which is not very user-friendly.
A custom user message can be provided using usrerr.New function and such error can
be wrapped into other errors.

In addition, HTTP code can be provided too which will be used by the rendering package
as the result REST HTTP code.

While HTTP code is required (must be not zero), user message can be empty which means
the error message itself will be used as user error.

Do not wrap multiple usrerr.Error types into a single error value as only the first one
(the top one) will be used for the rendering.

Use all-lowercase style for both errors and user messages.

Some predefined error values are available like ErrUnauthorized401 for use.
*/
package usrerr
