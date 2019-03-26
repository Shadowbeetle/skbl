package upower_test

import (
	"testing"

	"github.com/Shadowbeetle/set-kbd-blight/upower"
	"github.com/Shadowbeetle/set-kbd-blight/upower/mock"
)

func TestGetObject(t *testing.T) {
	mockConn := mock.NewDbusConnection()
	o := upower.GetObject(mockConn)

	if !mockConn.IsObjectStubCalled {
		t.Error("expected DbusConnection.Object to be called")
	}

	if mockConn.Args.Dest != "" {
		t.Errorf("expeceted s got %s", mockConn.Args.Dest)
	}
}
