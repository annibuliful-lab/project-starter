package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/sysulq/dataloader-go"
)

var PermissionsByIdsDataloader = dataloader.New(
	func(ctx context.Context, ids []int) []dataloader.Result[string] {
		results := make([]dataloader.Result[string], len(ids))
		// db, err := clients.NewPostgreSQLClient()

		// if err != nil {
		//     return results
		// }
		for i, key := range ids {
			results[i] = dataloader.Wrap(fmt.Sprintf("Result for %d", key), nil)
		}
		return results
	},
	dataloader.WithCache(100, time.Minute),
	dataloader.WithBatchSize(50),
	dataloader.WithWait(5*time.Millisecond),
)
