package integration

import (
	"testing"

	fileservice "github.com/aim/aim/internal/file/service"
	msgmodel "github.com/aim/aim/internal/message/model"
	msgservice "github.com/aim/aim/internal/message/service"
	"github.com/aim/aim/internal/pkg/database"
	relservice "github.com/aim/aim/internal/relation/service"

	airepo "github.com/aim/aim/internal/ai/repo"
	aiservice "github.com/aim/aim/internal/ai/service"
	filerepo "github.com/aim/aim/internal/file/repo"
	msgrepo "github.com/aim/aim/internal/message/repo"
	relrepo "github.com/aim/aim/internal/relation/repo"
	userrepo "github.com/aim/aim/internal/user/repo"
	userservice "github.com/aim/aim/internal/user/service"
)

// TestE2E_UserMessageFlow tests the full user registration -> messaging flow.
func TestE2E_UserMessageFlow(t *testing.T) {
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}

	userrepo.NewUserRepo(db).AutoMigrate()
	msgrepo.NewMessageRepo(db).AutoMigrate()

	userSvc := userservice.NewUserService(db)
	msgSvc := msgservice.NewMessageService(db)

	alice, err := userSvc.Register("alice_e2e", "pass", "Alice")
	if err != nil {
		t.Fatalf("Register alice error = %v", err)
	}

	bob, err := userSvc.Register("bob_e2e", "pass", "Bob")
	if err != nil {
		t.Fatalf("Register bob error = %v", err)
	}

	msg, err := msgSvc.Send(alice.ID, bob.ID, msgmodel.TargetUser, "text", "hello bob")
	if err != nil {
		t.Fatalf("Send error = %v", err)
	}
	if msg.FromID != alice.ID {
		t.Errorf("FromID = %d, want %d", msg.FromID, alice.ID)
	}

	msgs, err := msgSvc.GetMessages(bob.ID, alice.ID, msgmodel.TargetUser, 0, 20)
	if err != nil {
		t.Fatalf("GetMessages error = %v", err)
	}
	if len(msgs) < 1 {
		t.Fatal("expected at least 1 message")
	}

	loggedIn, err := userSvc.Login("alice_e2e", "pass")
	if err != nil {
		t.Fatalf("Login error = %v", err)
	}
	if loggedIn.ID != alice.ID {
		t.Errorf("login ID mismatch")
	}
}

// TestE2E_GroupChatFlow tests group creation -> invite -> messaging -> transfer.
func TestE2E_GroupChatFlow(t *testing.T) {
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}

	relrepo.NewRelationRepo(db).AutoMigrate()
	msgrepo.NewMessageRepo(db).AutoMigrate()

	relSvc := relservice.NewRelationService(db)
	msgSvc := msgservice.NewMessageService(db)

	group, err := relSvc.CreateGroup("dev-team-e2e", 1)
	if err != nil {
		t.Fatalf("CreateGroup error = %v", err)
	}
	if group.OwnerID != 1 {
		t.Errorf("OwnerID = %d, want 1", group.OwnerID)
	}

	relSvc.InviteMember(group.ID, 2, 1)
	relSvc.InviteMember(group.ID, 3, 1)

	members, _ := relSvc.GetGroupMembers(group.ID)
	if len(members) != 3 {
		t.Errorf("member count = %d, want 3", len(members))
	}

	_, err = msgSvc.Send(1, group.ID, msgmodel.TargetGroup, "text", "hello team")
	if err != nil {
		t.Fatalf("Send group msg error = %v", err)
	}

	relSvc.TransferOwner(group.ID, 2, 1)
	updatedGroup, _ := relSvc.GetGroup(group.ID)
	if updatedGroup.OwnerID != 2 {
		t.Errorf("OwnerID after transfer = %d, want 2", updatedGroup.OwnerID)
	}
}

// TestE2E_AIChatFlow tests bot creation -> AI chat -> billing.
func TestE2E_AIChatFlow(t *testing.T) {
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}

	airepo.NewAIRepo(db).AutoMigrate()

	aiSvc := aiservice.NewAIService(db, nil)

	bot, err := aiSvc.CreateBot("AIMBot-E2E", "openai", "gpt-4", "helpful", true, 0.03)
	if err != nil {
		t.Fatalf("CreateBot error = %v", err)
	}

	reply, err := aiSvc.Chat(1, bot.ID, 2, "user", "hello")
	if err != nil {
		t.Fatalf("Chat error = %v", err)
	}
	if reply == "" {
		t.Error("expected non-empty reply")
	}

	bots, _ := aiSvc.ListBots()
	if len(bots) != 1 {
		t.Errorf("bot count = %d, want 1", len(bots))
	}

	billing, _ := aiSvc.GetBilling(1, 20)
	if len(billing) != 1 {
		t.Errorf("billing rec count = %d, want 1", len(billing))
	}
}

// TestE2E_FileUploadFlow tests file save -> list -> delete.
func TestE2E_FileUploadFlow(t *testing.T) {
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}

	filerepo.NewFileRepo(db).AutoMigrate()
	fileSvc := fileservice.NewFileService(db, "./upload")

	f, err := fileSvc.SaveFile(1, "report.pdf", "/upload/report.pdf", 51200, "application/pdf")
	if err != nil {
		t.Fatalf("SaveFile error = %v", err)
	}

	files, _ := fileSvc.ListFiles(1, 0, 10)
	if len(files) != 1 {
		t.Errorf("file count = %d, want 1", len(files))
	}

	fileSvc.DeleteFile(f.ID, 1)
	files, _ = fileSvc.ListFiles(1, 0, 10)
	if len(files) != 0 {
		t.Errorf("file count after delete = %d, want 0", len(files))
	}
}

// TestE2E_FriendFlow tests add friend -> remark -> delete.
func TestE2E_FriendFlow(t *testing.T) {
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}

	relrepo.NewRelationRepo(db).AutoMigrate()
	relSvc := relservice.NewRelationService(db)

	relSvc.AddFriend(1, 2, "bestie")
	friends, _ := relSvc.ListFriends(1)
	if len(friends) != 1 {
		t.Errorf("friend count = %d, want 1", len(friends))
	}

	relSvc.UpdateRemark(1, 2, "BFF")
	relSvc.DeleteFriend(1, 2)

	friends, _ = relSvc.ListFriends(1)
	if len(friends) != 0 {
		t.Errorf("friend count after delete = %d, want 0", len(friends))
	}
}
