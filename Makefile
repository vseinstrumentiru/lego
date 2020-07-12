include env.mk
include main.mk
# Update main targets
init:
	curl https://raw.githubusercontent.com/sagikazarmark/makefiles/master/go-binary/main.mk > main.mk
