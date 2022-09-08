#
# This file only contains the rule that generate the
# help content from the comments in the different files.
#
# Use '##@ My group text' at the beginning of a line to
# print out a group text.
#
# Use '## My help text' at the end of a rule to print out
# content related with a rule. Try to short the description.
#

.PHONY: help
help: ## Print out the help content
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
