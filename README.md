# laplace-api-golang

Client library for Laplace API for the US stock market and BIST (Istanbul stock market) fundamental financial data.

## Instantiating a Client

```go
import (
	"context"
	"fmt"
	laplace "github.com/Laplace-Analytics/laplace-api-golang"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a new client
	client := laplace.NewClient(laplace.Config{
		APIKey: "your_api_key_here",
	}, logrus.New())

	// Create a context
	ctx := context.Background()

	// Example: Get all stocks
	stocks, err := client.GetAllStocks(ctx, laplace.RegionUs)
	if err != nil {
		fmt.Println("Error fetching stocks:", err)
		return
	}

	// Print the stocks
	for _, stock := range stocks {
		fmt.Printf("Stock: %s, Name: %s\n", stock.Symbol, stock.Name)
	}
}
```
