SOURCE_MAKE=. make.sh

b: test
test:
	@${SOURCE_MAKE} && test

d: run
run:
	@${SOURCE_MAKE} && run