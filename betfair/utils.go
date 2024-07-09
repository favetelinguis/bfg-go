package betfair

import "fmt"

func createUrl(endpoint string, method string) string {
	return endpoint + method
}

func createJsonRpcMethodName(name string) string {
	return fmt.Sprintf("SportsAPING/v1.0/%s", name)
}
