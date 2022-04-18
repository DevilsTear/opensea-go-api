include META
.PHONY: push
push:
	@echo "[BUILD] Committing and pushing to remote repository"
	@echo " - Committing"
	@git add META
	@git commit -am "v$(VERSION)"
	@echo " - Tagging"
	@git tag v${VERSION}
	@echo " - Pushing"
	@git push --tags origin ${BRANCH}