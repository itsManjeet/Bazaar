package main

func handleFlatpak(data string) {
	if err := doProcess([]string{"flatpak", "install", data}, "/tmp"); err != nil {
		showError(err.Error())
	}
}

func handleRecipie(data string) {
	data = data[0 : len(data)-7]
	if err := doProcess([]string{"sys-appbuild"}, data); err != nil {
		showError(err.Error())
	}
}

func handlePkg(data string) {
	if err := doProcess([]string{"sys-appadd", data}, "/tmp"); err != nil {
		showError(err.Error())
	}
}
