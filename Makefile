.PHONY: init

init:
	chmod +x ./hooks/pre-commit.sh && ln -f ./hooks/pre-commit.sh ./.git/hooks/pre-commit
