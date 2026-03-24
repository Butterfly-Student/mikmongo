//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/internal/repository/postgres"
)

// createTestMikrotikRouter inserts a minimal mikrotik_routers row and returns its ID.
func createTestMikrotikRouter(t *testing.T, suite *TestSuite) string {
	t.Helper()
	id := uuid.New().String()
	err := suite.DB.WithContext(suite.Ctx).Exec(
		`INSERT INTO mikrotik_routers (id, name, address, username, password_encrypted, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		id, "Test Router", "192.168.1.1", "admin", "enc_placeholder",
	).Error
	require.NoError(t, err)
	return id
}

// createTestHotspotSale inserts a single hotspot_sale row and returns the model.
func createTestHotspotSale(t *testing.T, suite *TestSuite, routerID string, agentID *string, profile, batchCode string) *model.HotspotSale {
	t.Helper()
	sale := &model.HotspotSale{
		RouterID:     routerID,
		Username:     "user-" + uuid.New().String()[:8],
		Profile:      profile,
		Price:        5000,
		SellingPrice: 7000,
		Prefix:       "ts",
		BatchCode:    batchCode,
		SalesAgentID: agentID,
	}
	repos := postgres.NewRepository(suite.DB)
	err := repos.HotspotSaleRepo.Create(suite.Ctx, sale)
	require.NoError(t, err)
	return sale
}

func TestHotspotSale_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	sale := &model.HotspotSale{
		RouterID:     routerID,
		Username:     "voucher001",
		Profile:      "10mb",
		Price:        5000,
		SellingPrice: 7000,
		Prefix:       "ts",
		BatchCode:    "ABCD",
	}

	err := repos.HotspotSaleRepo.Create(suite.Ctx, sale)
	require.NoError(t, err)
	assert.NotEmpty(t, sale.ID)

	got, err := repos.HotspotSaleRepo.GetByID(suite.Ctx, uuid.MustParse(sale.ID))
	require.NoError(t, err)
	assert.Equal(t, sale.Username, got.Username)
	assert.Equal(t, sale.Profile, got.Profile)
	assert.Equal(t, sale.Price, got.Price)
	assert.Equal(t, sale.SellingPrice, got.SellingPrice)
	assert.Equal(t, sale.BatchCode, got.BatchCode)
	assert.Nil(t, got.SalesAgentID)
}

func TestHotspotSale_CreateBatch_Empty(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	err := repos.HotspotSaleRepo.CreateBatch(suite.Ctx, []model.HotspotSale{})
	require.NoError(t, err)
}

func TestHotspotSale_CreateBatch(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	sales := make([]model.HotspotSale, 5)
	for i := range sales {
		sales[i] = model.HotspotSale{
			RouterID:  routerID,
			Username:  fmt.Sprintf("batch-user-%d", i),
			Profile:   "10mb",
			BatchCode: "BATCH1",
		}
	}

	err := repos.HotspotSaleRepo.CreateBatch(suite.Ctx, sales)
	require.NoError(t, err)

	count, err := repos.HotspotSaleRepo.Count(suite.Ctx, repository.HotspotSaleFilter{BatchCode: "BATCH1"})
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestHotspotSale_List_NoFilter(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "C1")
	createTestHotspotSale(t, suite, routerID, nil, "20mb", "C2")
	createTestHotspotSale(t, suite, routerID, nil, "5mb", "C3")

	list, err := repos.HotspotSaleRepo.List(suite.Ctx, repository.HotspotSaleFilter{}, 100, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
}

func TestHotspotSale_List_FilterRouterID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerA := createTestMikrotikRouter(t, suite)
	routerB := createTestMikrotikRouter(t, suite)

	createTestHotspotSale(t, suite, routerA, nil, "10mb", "RA1")
	createTestHotspotSale(t, suite, routerA, nil, "10mb", "RA2")
	createTestHotspotSale(t, suite, routerB, nil, "10mb", "RB1")

	routerAID := uuid.MustParse(routerA)
	filter := repository.HotspotSaleFilter{RouterID: &routerAID}

	list, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 100, 0)
	require.NoError(t, err)
	assert.Len(t, list, 2)
	for _, s := range list {
		assert.Equal(t, routerA, s.RouterID)
	}
}

func TestHotspotSale_List_FilterProfile(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "P1")
	createTestHotspotSale(t, suite, routerID, nil, "20mb", "P2")

	filter := repository.HotspotSaleFilter{Profile: "10mb"}
	list, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 100, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
	for _, s := range list {
		assert.Equal(t, "10mb", s.Profile)
	}
}

func TestHotspotSale_List_FilterBatchCode(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "BCH1")
	createTestHotspotSale(t, suite, routerID, nil, "10mb", "BCH1")
	createTestHotspotSale(t, suite, routerID, nil, "10mb", "BCH2")

	filter := repository.HotspotSaleFilter{BatchCode: "BCH1"}
	list, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 100, 0)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestHotspotSale_List_FilterDateRange(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "DR1")

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	filter := repository.HotspotSaleFilter{DateFrom: &yesterday, DateTo: &tomorrow}
	list, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 100, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// Future range — should be empty
	future := now.AddDate(0, 0, 2)
	filterFuture := repository.HotspotSaleFilter{DateFrom: &tomorrow, DateTo: &future}
	listFuture, err := repos.HotspotSaleRepo.List(suite.Ctx, filterFuture, 100, 0)
	require.NoError(t, err)
	assert.Len(t, listFuture, 0)
}

func TestHotspotSale_Count(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	rID := uuid.MustParse(routerID)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "CT1")
	createTestHotspotSale(t, suite, routerID, nil, "10mb", "CT2")

	filter := repository.HotspotSaleFilter{RouterID: &rID}
	count, err := repos.HotspotSaleRepo.Count(suite.Ctx, filter)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestHotspotSale_List_Pagination(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	rID := uuid.MustParse(routerID)

	for i := 0; i < 5; i++ {
		createTestHotspotSale(t, suite, routerID, nil, "10mb", fmt.Sprintf("PG%d", i))
	}

	filter := repository.HotspotSaleFilter{RouterID: &rID}

	page1, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 2, 0)
	require.NoError(t, err)
	assert.Len(t, page1, 2)

	page2, err := repos.HotspotSaleRepo.List(suite.Ctx, filter, 2, 2)
	require.NoError(t, err)
	assert.Len(t, page2, 2)

	// IDs should not overlap
	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

func TestHotspotSale_ListByBatchCode(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	rID := uuid.MustParse(routerID)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "LBC1")
	createTestHotspotSale(t, suite, routerID, nil, "10mb", "LBC1")
	createTestHotspotSale(t, suite, routerID, nil, "20mb", "LBC2")

	list, err := repos.HotspotSaleRepo.ListByBatchCode(suite.Ctx, rID, "LBC1")
	require.NoError(t, err)
	assert.Len(t, list, 2)
	for _, s := range list {
		assert.Equal(t, "LBC1", s.BatchCode)
	}
}

func TestHotspotSale_DeleteByBatchCode(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	rID := uuid.MustParse(routerID)

	createTestHotspotSale(t, suite, routerID, nil, "10mb", "DEL1")
	createTestHotspotSale(t, suite, routerID, nil, "10mb", "DEL1")

	err := repos.HotspotSaleRepo.DeleteByBatchCode(suite.Ctx, rID, "DEL1")
	require.NoError(t, err)

	remaining, err := repos.HotspotSaleRepo.ListByBatchCode(suite.Ctx, rID, "DEL1")
	require.NoError(t, err)
	assert.Empty(t, remaining)
}
