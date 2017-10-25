package main

import (
	"ims_api_connector"
	"fmt"
)

func main() {
	connector := ims_api_connector.New("some_username", "some_password", "some.server.com:8000", 5)

	connector.Authenticate()

	assets, _ := connector.GetAssets()

	fmt.Println(assets)
}
