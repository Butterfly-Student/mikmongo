//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
)

// createTestSalesAgent inserts a minimal SalesAgent for the given routerID and returns the model.
func createTestSalesAgent(t *testing.T, suite *TestSuite, routerID string) *model.SalesAgent {
	t.Helper()
	repos := postgres.NewRepository(suite.DB)
	agent := &model.SalesAgent{
		RouterID:     routerID,
		Name:         "Test Agent " + uuid.New().String()[:8],
		Username:     "agent-" + uuid.New().String()[:8],
		PasswordHash: "$2a$04$placeholder",
		Status:       "active",
		VoucherMode:  "mix",
		VoucherLength: 6,
		VoucherType:  "upp",
	}
	err := repos.SalesAgentRepo.Create(suite.Ctx, agent)
	require.NoError(t, err)
	return agent
}

func TestSalesAgent_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	phone := "08123456789"
	agent := &model.SalesAgent{
		RouterID:      routerID,
		Name:          "Budi Santoso",
		Phone:         &phone,
		Username:      "budisantoso",
		PasswordHash:  "$2a$04$placeholder",
		Status:        "active",
		VoucherMode:   "mix",
		VoucherLength: 8,
		VoucherType:   "upp",
		BillDiscount:  10.0,
	}

	err := repos.SalesAgentRepo.Create(suite.Ctx, agent)
	require.NoError(t, err)
	assert.NotEmpty(t, agent.ID)

	got, err := repos.SalesAgentRepo.GetByID(suite.Ctx, uuid.MustParse(agent.ID))
	require.NoError(t, err)
	assert.Equal(t, "Budi Santoso", got.Name)
	assert.Equal(t, "budisantoso", got.Username)
	assert.Equal(t, "active", got.Status)
	assert.Equal(t, 8, got.VoucherLength)
	assert.Equal(t, 10.0, got.BillDiscount)
	require.NotNil(t, got.Phone)
	assert.Equal(t, phone, *got.Phone)
}

func TestSalesAgent_Create_DuplicateUsername(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	agent1 := &model.SalesAgent{
		RouterID: routerID, Name: "Agent One", Username: "dupuser",
		PasswordHash: "$2a$04$ph", Status: "active", VoucherMode: "mix", VoucherType: "upp",
	}
	agent2 := &model.SalesAgent{
		RouterID: routerID, Name: "Agent Two", Username: "dupuser",
		PasswordHash: "$2a$04$ph", Status: "active", VoucherMode: "mix", VoucherType: "upp",
	}

	require.NoError(t, repos.SalesAgentRepo.Create(suite.Ctx, agent1))
	err := repos.SalesAgentRepo.Create(suite.Ctx, agent2)
	assert.Error(t, err, "duplicate username should fail")
}

func TestSalesAgent_GetByUsername(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	got, err := repos.SalesAgentRepo.GetByUsername(suite.Ctx, agent.Username)
	require.NoError(t, err)
	assert.Equal(t, agent.ID, got.ID)
}

func TestSalesAgent_GetByUsername_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	_, err := repos.SalesAgentRepo.GetByUsername(suite.Ctx, "nonexistent-user-xyz")
	assert.Error(t, err)
}

func TestSalesAgent_Update(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	phone := "08111222333"
	agent.Name = "Updated Name"
	agent.Phone = &phone
	agent.Status = "inactive"

	err := repos.SalesAgentRepo.Update(suite.Ctx, agent)
	require.NoError(t, err)

	got, err := repos.SalesAgentRepo.GetByID(suite.Ctx, uuid.MustParse(agent.ID))
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "inactive", got.Status)
	require.NotNil(t, got.Phone)
	assert.Equal(t, phone, *got.Phone)
}

func TestSalesAgent_Delete_SoftDelete(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	agentID := uuid.MustParse(agent.ID)
	err := repos.SalesAgentRepo.Delete(suite.Ctx, agentID)
	require.NoError(t, err)

	// GetByID should return error for soft-deleted record
	_, err = repos.SalesAgentRepo.GetByID(suite.Ctx, agentID)
	assert.Error(t, err, "soft-deleted agent should not be found")
}

func TestSalesAgent_Delete_SoftDelete_RecordStillExists(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	agentID := uuid.MustParse(agent.ID)
	require.NoError(t, repos.SalesAgentRepo.Delete(suite.Ctx, agentID))

	// Row should still exist in DB with deleted_at set (not hard deleted)
	var count int64
	err := suite.DB.WithContext(suite.Ctx).
		Unscoped().
		Model(&model.SalesAgent{}).
		Where("id = ? AND deleted_at IS NOT NULL", agentID).
		Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSalesAgent_List_NoFilter(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)

	createTestSalesAgent(t, suite, routerID)
	createTestSalesAgent(t, suite, routerID)

	list, err := repos.SalesAgentRepo.List(suite.Ctx, nil, 100, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

func TestSalesAgent_List_FilterRouterID(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerA := createTestMikrotikRouter(t, suite)
	routerB := createTestMikrotikRouter(t, suite)

	createTestSalesAgent(t, suite, routerA)
	createTestSalesAgent(t, suite, routerA)
	createTestSalesAgent(t, suite, routerB)

	rA := uuid.MustParse(routerA)
	list, err := repos.SalesAgentRepo.List(suite.Ctx, &rA, 100, 0)
	require.NoError(t, err)
	assert.Len(t, list, 2)
	for _, a := range list {
		assert.Equal(t, routerA, a.RouterID)
	}
}

func TestSalesAgent_Count(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	rID := uuid.MustParse(routerID)

	createTestSalesAgent(t, suite, routerID)
	createTestSalesAgent(t, suite, routerID)

	count, err := repos.SalesAgentRepo.Count(suite.Ctx, &rID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestSalesAgent_UpsertProfilePrice_Create(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	price := &model.SalesProfilePrice{
		SalesAgentID: agent.ID,
		ProfileName:  "10mb",
		BasePrice:    5000,
		SellingPrice: 7000,
		IsActive:     true,
	}

	err := repos.SalesAgentRepo.UpsertProfilePrice(suite.Ctx, price)
	require.NoError(t, err)
	assert.NotEmpty(t, price.ID)

	got, err := repos.SalesAgentRepo.GetProfilePrice(suite.Ctx, uuid.MustParse(agent.ID), "10mb")
	require.NoError(t, err)
	assert.Equal(t, 5000.0, got.BasePrice)
	assert.Equal(t, 7000.0, got.SellingPrice)
}

func TestSalesAgent_UpsertProfilePrice_Update(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)
	agentID := uuid.MustParse(agent.ID)

	// Create
	price := &model.SalesProfilePrice{
		SalesAgentID: agent.ID, ProfileName: "10mb",
		BasePrice: 5000, SellingPrice: 7000, IsActive: true,
	}
	require.NoError(t, repos.SalesAgentRepo.UpsertProfilePrice(suite.Ctx, price))

	// Update with new prices
	priceUpdate := &model.SalesProfilePrice{
		SalesAgentID: agent.ID, ProfileName: "10mb",
		BasePrice: 8000, SellingPrice: 10000, IsActive: true,
	}
	err := repos.SalesAgentRepo.UpsertProfilePrice(suite.Ctx, priceUpdate)
	require.NoError(t, err)

	got, err := repos.SalesAgentRepo.GetProfilePrice(suite.Ctx, agentID, "10mb")
	require.NoError(t, err)
	assert.Equal(t, 8000.0, got.BasePrice)
	assert.Equal(t, 10000.0, got.SellingPrice)
}

func TestSalesAgent_GetProfilePrice_NotFound(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	_, err := repos.SalesAgentRepo.GetProfilePrice(suite.Ctx, uuid.MustParse(agent.ID), "nonexistent-profile")
	assert.Error(t, err)
}

func TestSalesAgent_ListProfilePrices(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)

	for _, profile := range []string{"10mb", "20mb", "5mb"} {
		p := &model.SalesProfilePrice{
			SalesAgentID: agent.ID, ProfileName: profile,
			BasePrice: 5000, SellingPrice: 7000, IsActive: true,
		}
		require.NoError(t, repos.SalesAgentRepo.UpsertProfilePrice(suite.Ctx, p))
	}

	list, err := repos.SalesAgentRepo.ListProfilePrices(suite.Ctx, uuid.MustParse(agent.ID))
	require.NoError(t, err)
	assert.Len(t, list, 3)
	// Should be ordered alphabetically by profile_name
	assert.Equal(t, "10mb", list[0].ProfileName)
	assert.Equal(t, "20mb", list[1].ProfileName)
	assert.Equal(t, "5mb", list[2].ProfileName)
}

func TestSalesAgent_ProfilePrice_Cascade(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	routerID := createTestMikrotikRouter(t, suite)
	agent := createTestSalesAgent(t, suite, routerID)
	agentID := uuid.MustParse(agent.ID)

	p := &model.SalesProfilePrice{
		SalesAgentID: agent.ID, ProfileName: "10mb",
		BasePrice: 5000, SellingPrice: 7000, IsActive: true,
	}
	require.NoError(t, repos.SalesAgentRepo.UpsertProfilePrice(suite.Ctx, p))

	// Hard delete (unscoped) to trigger CASCADE
	err := suite.DB.WithContext(suite.Ctx).Unscoped().Delete(&model.SalesAgent{}, "id = ?", agentID).Error
	require.NoError(t, err)

	// Profile prices should be cascaded
	var count int64
	suite.DB.WithContext(suite.Ctx).Model(&model.SalesProfilePrice{}).
		Where("sales_agent_id = ?", agentID).Count(&count)
	assert.Equal(t, int64(0), count)
}
