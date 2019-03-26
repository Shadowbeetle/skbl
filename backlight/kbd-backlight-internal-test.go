package backlight

// import (
// 	"testing"
// 	"time"
// )

// func TestNewKbdBacklight(t *testing.T) {
// 	// expectedConnDest := "org.freedesktop.UPower"

// 	mockDConn := &MockDbusConnection{}

// 	conf := Config{
// 		IdleWaitTime:   time.Duration(5) * time.Second,
// 		InputPaths:     []string{"/test/input/kbd"},
// 		dbusConnection: mockDConn,
// 	}

// 	_ /*kbl*/, err := NewKbdBacklight(conf)
// 	if err != nil {
// 		t.Fatalf("expected nil error got %s instead", err.Error())
// 	}

// 	if !mockDConn.isObjectStubCalled {
// 		t.Errorf("expected MockDbusConnection.Object stub to be called, got false")
// 	}

// 	// if mockDConn.args.dest !=
// 	// TODO: should call conn.Object("org.freedesktop.UPower", "/org/freedesktop/UPower/KbdBacklight")
// 	// TODO: should set brPtr and should set kbl.desiredBrightness
// 	// TODO: should call busObject.AddMatchSignal("org.freedesktop.UPower.KbdBacklight", "BrightnessChangedWithSource")
// 	// TODO: should retun &KbdBacklight with proper setup
// }
