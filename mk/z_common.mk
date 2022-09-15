#
# Common targets (included as the last, undocumented, used on CI)
#

MAKE_DOC=docs/make.md

.PHONY: validate
validate: validate-spec validate-clients

