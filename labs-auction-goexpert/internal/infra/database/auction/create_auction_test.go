package auction

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionShouldCloseAutomatically(t *testing.T) {
	repo, ctx := newTestAuctionRepository(t)

	const auctionDuration = 50 * time.Millisecond
	t.Setenv("AUCTION_DURATION", auctionDuration.String())

	auction := &auction_entity.Auction{
		Id:          fmt.Sprintf("auction-test-%d", time.Now().UnixNano()),
		ProductName: "Test Product",
		Category:    "Test",
		Description: "Test Description long enough",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	err := repo.CreateAuction(ctx, auction)
	require.NoError(t, err)

	found, err := repo.FindAuctionById(ctx, auction.Id)
	require.NoError(t, err)
	assert.Equal(t, found.Status, auction_entity.Active, "expected auction status Active, got %v", found.Status)

	time.Sleep(auctionDuration + 10*time.Millisecond)

	found, err = repo.FindAuctionById(ctx, auction.Id)
	require.NoError(t, err)
	assert.Equal(t, found.Status, auction_entity.Completed, "expected auction status Completed, got %v", found.Status)
}

func newTestAuctionRepository(t *testing.T) (*AuctionRepository, context.Context) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in -short")
	}

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		t.Skip("set MONGODB_URL to run this integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = client.Disconnect(cleanupCtx)
	})

	require.NoError(t, client.Ping(ctx, nil))

	db := client.Database("auctions_test_" + time.Now().Format("20060102150405.000000000"))
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})

	return NewAuctionRepository(db), ctx
}
