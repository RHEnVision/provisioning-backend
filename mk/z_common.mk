#
# Common targets (included as the last, undocumented, used on CI)
#

MAKE_DOC=docs/make.md

.PHONY: validate-make
validate-make:
	echo '# Make documentation' > $(MAKE_DOC)
	echo '```' >> $(MAKE_DOC)
	make help | sed -r "s/\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]//g" >> $(MAKE_DOC)
	echo '```' >> $(MAKE_DOC)

.PHONY: validate
validate: validate-make validate-spec validate-clients

