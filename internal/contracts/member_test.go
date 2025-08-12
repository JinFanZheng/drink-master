package contracts

import (
	"testing"
	"time"
)

func TestUpdateMemberRequest(t *testing.T) {
	req := UpdateMemberRequest{
		ID:       "member-123",
		Nickname: "张三",
		Avatar:   "https://example.com/avatar.jpg",
	}

	if req.ID != "member-123" {
		t.Errorf("expected ID to be 'member-123', got '%s'", req.ID)
	}

	if req.Nickname != "张三" {
		t.Errorf("expected Nickname to be '张三', got '%s'", req.Nickname)
	}

	if req.Avatar != "https://example.com/avatar.jpg" {
		t.Errorf("expected Avatar to be 'https://example.com/avatar.jpg', got '%s'", req.Avatar)
	}
}

func TestCreateFranchiseIntentionRequest(t *testing.T) {
	req := CreateFranchiseIntentionRequest{
		MemberID:         "member-123",
		ContactName:      "张三",
		ContactPhone:     "13800138000",
		IntendedLocation: "北京市朝阳区",
		Remarks:          "希望在繁华地段开店",
	}

	if req.MemberID != "member-123" {
		t.Errorf("expected MemberID to be 'member-123', got '%s'", req.MemberID)
	}

	if req.ContactName != "张三" {
		t.Errorf("expected ContactName to be '张三', got '%s'", req.ContactName)
	}
}

func TestFranchiseIntentionStatus(t *testing.T) {
	testCases := []struct {
		status   FranchiseIntentionStatus
		expected string
	}{
		{FranchiseIntentionStatusPending, "Pending"},
		{FranchiseIntentionStatusApproved, "Approved"},
		{FranchiseIntentionStatusRejected, "Rejected"},
	}

	for _, tc := range testCases {
		if string(tc.status) != tc.expected {
			t.Errorf("expected status '%s', got '%s'", tc.expected, string(tc.status))
		}
	}
}

func TestGetMemberInfoResponse(t *testing.T) {
	now := time.Now()
	intentions := []FranchiseIntentionSummary{
		{
			ID:               "fi-123",
			ContactName:      "张三",
			IntendedLocation: "北京市",
			Status:           FranchiseIntentionStatusPending,
			CreatedAt:        now,
		},
	}

	response := GetMemberInfoResponse{
		ID:                  "member-123",
		Nickname:            "张三",
		Avatar:              "https://example.com/avatar.jpg",
		WeChatOpenID:        "wx_openid_123",
		Role:                "Member",
		IsAdmin:             false,
		CreatedAt:           now,
		UpdatedAt:           now,
		FranchiseIntentions: intentions,
	}

	if response.ID != "member-123" {
		t.Errorf("expected ID to be 'member-123', got '%s'", response.ID)
	}

	if len(response.FranchiseIntentions) != 1 {
		t.Errorf("expected 1 franchise intention, got %d", len(response.FranchiseIntentions))
	}

	if response.FranchiseIntentions[0].Status != FranchiseIntentionStatusPending {
		t.Errorf("expected status to be 'Pending', got '%s'", string(response.FranchiseIntentions[0].Status))
	}
}

func TestMemberErrorCodes(t *testing.T) {
	testCases := []struct {
		code     string
		expected string
	}{
		{ErrorCodeMemberNotFound, "MEMBER_NOT_FOUND"},
		{ErrorCodeMemberUpdateFailed, "MEMBER_UPDATE_FAILED"},
		{ErrorCodeFranchiseExists, "FRANCHISE_INTENTION_EXISTS"},
		{ErrorCodeFranchiseCreateFail, "FRANCHISE_CREATE_FAILED"},
		{ErrorCodeInvalidMemberData, "INVALID_MEMBER_DATA"},
	}

	for _, tc := range testCases {
		if tc.code != tc.expected {
			t.Errorf("expected error code '%s', got '%s'", tc.expected, tc.code)
		}
	}
}
