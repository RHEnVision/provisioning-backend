#
# Common targets (included as the last, undocumented, used on CI)
#

.PHONY: validate
validate: validate-spec validate-clients validate-example-config

