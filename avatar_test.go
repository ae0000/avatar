package avatar

import (
	"os"
	"testing"
)

func TestCleanString(t *testing.T) {
	// Input should equal output
	tests := map[string]string{}
	tests["AE"] = "AE"
	tests["ae"] = "AE"
	tests["a e"] = "AE"
	tests["andrew edwards"] = "AE"
	tests["andrew   edwards"] = "AN"
	tests["a"] = "A"
	tests["123"] = "12"
	tests["A 3"] = "A3"
	tests["B 3"] = "B3"

	for k, v := range tests {
		if cleanString(k) != v {
			t.Errorf("Received '%s', was expecting '%s' from '%s'", cleanString(k), v, k)
		}
	}
}

func TestGetFont(t *testing.T) {
	f, err := getFont("rubbish")
	if err == nil {
		t.Error("should not have been able to return a font")
	}
	if f != nil {
		t.Error("should not have been able to return a font")
	}
	f, err = getFont(defaultfontFace)
	if err != nil {
		t.Error("should of been able to return a font", err)
	}
	if f == nil {
		t.Error("should of been able to return a font")
	}
}

func TestCreateAvatar(t *testing.T) {
	tests := []struct {
		initials  string
		bgColor   string
		fontColor string
	}{
		{"AE", "", ""},
		{"JD", "", ""},
		{"A", "", ""},
		{"XY", "#FF0000", "#FFFFFF"},
		{"TZ", "#00FF00", ""},
	}

	for _, tt := range tests {
		img, err := createAvatar(tt.initials, tt.bgColor, tt.fontColor)
		if err != nil {
			t.Errorf("createAvatar(%s, %s, %s) returned error: %v", tt.initials, tt.bgColor, tt.fontColor, err)
			continue
		}
		if img == nil {
			t.Errorf("createAvatar(%s, %s, %s) returned nil image", tt.initials, tt.bgColor, tt.fontColor)
		}
		bounds := img.Bounds()
		if bounds.Dx() != int(imageWidth) || bounds.Dy() != int(imageHeight) {
			t.Errorf("createAvatar(%s, %s, %s) returned image with wrong dimensions: got %dx%d, want %dx%d",
				tt.initials, tt.bgColor, tt.fontColor, bounds.Dx(), bounds.Dy(), int(imageWidth), int(imageHeight))
		}
	}
}

func TestToDisk(t *testing.T) {
	tmpFile := "test_avatar.png"
	defer os.Remove(tmpFile)

	err := ToDisk("TD", tmpFile)
	if err != nil {
		t.Errorf("ToDisk returned error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("ToDisk did not create file")
	}
}

func TestToDiskCustom(t *testing.T) {
	tmpFile := "test_avatar_custom.png"
	defer os.Remove(tmpFile)

	err := ToDiskCustom("TC", tmpFile, "#FF0000", "#FFFFFF")
	if err != nil {
		t.Errorf("ToDiskCustom returned error: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("ToDiskCustom did not create file")
	}
}

func TestCacheFunctionality(t *testing.T) {
	// Use unique initials to avoid conflicts with other tests
	initials := "CF"

	// Delete from cache if it exists
	imageCache.Delete(initials)

	// First call should create and cache the image
	img1, err := createAvatar(initials, "", "")
	if err != nil {
		t.Fatalf("createAvatar failed: %v", err)
	}

	// Check that image is in cache
	cachedImg := getImage(initials)
	if cachedImg == nil {
		t.Error("Image was not cached after creation")
	}

	// Second call should return cached image
	img2, err := createAvatar(initials, "", "")
	if err != nil {
		t.Fatalf("createAvatar failed on second call: %v", err)
	}

	// Both should point to the same underlying image (from cache)
	if img1 != img2 {
		t.Error("Second call did not return cached image")
	}

	// Cleanup
	imageCache.Delete(initials)
}

func TestSetGetImage(t *testing.T) {
	initials := "SG"

	// Delete from cache if it exists
	imageCache.Delete(initials)

	// Should return nil when not in cache
	img := getImage(initials)
	if img != nil {
		t.Error("getImage should return nil for non-existent entry")
	}

	// Create an image and cache it
	avatar, err := createAvatar(initials, "", "")
	if err != nil {
		t.Fatalf("createAvatar failed: %v", err)
	}

	setImage(initials, avatar)

	// Should now return the cached image
	cachedImg := getImage(initials)
	if cachedImg == nil {
		t.Error("getImage returned nil after setImage")
	}
	if cachedImg != avatar {
		t.Error("getImage returned different image than what was set")
	}

	// Cleanup
	imageCache.Delete(initials)
}
