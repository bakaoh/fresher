package counter

import (
	"context"
	fmt "fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func BenchmarkCounter(b *testing.B) {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("PORT")), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := NewCounterServiceClient(conn)
	b.ResetTimer()

	b.Run("Set", func(b *testing.B) {
		rand.Seed(56789)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err := client.SetBalance(context.Background(), &UserReq{
					UserId:  strconv.FormatInt(rand.Int63(), 16),
					Balance: 100,
				})
				assert.Nil(b, err)
			}
		})
	})

	b.Run("Get", func(b *testing.B) {
		rand.Seed(56789)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				balanceRes, err := client.GetBalance(context.Background(), &UserReq{
					UserId: strconv.FormatInt(rand.Int63(), 16),
				})
				assert.Nil(b, err)
				assert.Contains(b, []int64{0, 100}, balanceRes.Balance)
			}
		})
	})

	b.Run("Inc", func(b *testing.B) {
		rand.Seed(56789)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err := client.IncreaseBalance(context.Background(), &UserReq{
					UserId:  strconv.FormatInt(rand.Int63(), 16),
					Balance: 1,
				})
				assert.Nil(b, err)
			}
		})
	})

	b.Run("Dec", func(b *testing.B) {
		rand.Seed(56789)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				_, err := client.DecreaseBalance(context.Background(), &UserReq{
					UserId:  strconv.FormatInt(rand.Int63(), 16),
					Balance: 1,
				})
				assert.Nil(b, err)
			}
		})
	})

	rand.Seed(56789)
	for i := 0; i < 1000; i++ {
		balanceRes, err := client.GetBalance(context.Background(), &UserReq{
			UserId: strconv.FormatInt(rand.Int63(), 16),
		})
		assert.Nil(b, err)
		assert.Equal(b, int64(100), balanceRes.Balance)
	}
}
