install:
		touch ~/.ghuser
		cp .env ~/.ghuser
homebrew:
		touch ${HOME}/.ghuser
		cp .env-sample ${HOME}/.ghuser
uninstall:
		rm ~/.ghuser
