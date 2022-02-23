package util

import (
	"testing"
)

func TestCompareTwoLinks(t *testing.T) {
	similarity := CompareTwoLinks("https://discord.com", "https://discord.com")

	if similarity != 1 {
		t.Errorf("Comparing a link with itself did not return a similarity of 1")
	}

	similarity = CompareTwoLinks("", "https://discord.com")
	if similarity != 0 {
		t.Errorf("Comparing one link with an empty string should return a similarity of 0")
	}

	similarity = CompareTwoLinks("https://dlscord.com", "https://discord.com")
	if similarity < 0.5 {
		t.Errorf("Comparint one links with only one character difference should return a high similarity")
	}

}
func TestExtractLinks(t *testing.T) {
	links := ExtractLinks("https://discord.com")

	if len(links) != 1 || links[0] != "discord.com" {
		t.Errorf("Could not detect https link")
	}

	links = ExtractLinks("http://discord.com")

	if len(links) != 1 || links[0] != "discord.com" {
		t.Errorf("Could not detect http link")
	}

	links = ExtractLinks("[[https://discord.com]]()jdfkajdkfj")

	if len(links) != 1 || links[0] != "discord.com" {
		t.Errorf("Could not detect link inside a message")
	}

	links = ExtractLinks("discord.com")

	if len(links) != 0 {
		t.Errorf("Detected a link without a scheme")
	}

	links = ExtractLinks("https://discord.co.uk")

	if len(links) != 1 || links[0] != "discord.co.uk" {
		t.Errorf("Did not detect link with double TLD (.co.uk)")
	}

}
