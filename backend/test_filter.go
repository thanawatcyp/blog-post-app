package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"blog-app-backend/services"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}

func main() {
	// Load environment variables from .env file
	loadEnv()

	// Test cases for content filtering
	testCases := []struct {
		title    string
		content  string
		expected string
	}{
		{
			title:    "Clean Post",
			content:  "This is a wonderful day and I love programming!",
			expected: "CLEAN",
		},
		{
			title:    "Hate Speech Test",
			content:  "I hate all people from different backgrounds and they should leave",
			expected: "INAPPROPRIATE",
		},
		{
			title:    "Profanity Test",
			content:  "This is some damn bullshit content with bad words",
			expected: "INAPPROPRIATE",
		},
		{
			title:    "Technical Content",
			content:  "Here's how to implement a REST API with Go and Fiber framework",
			expected: "CLEAN",
		},
	}

	// Check if API key is set
	if os.Getenv("DEEPSEEK_API_KEY") == "" {
		log.Fatal("Please set DEEPSEEK_API_KEY environment variable")
	}

	contentFilter := services.NewContentFilterService()

	fmt.Println("🧪 Testing Content Filter with DeepSeek AI")
	fmt.Println("==========================================")

	for i, test := range testCases {
		fmt.Printf("\n📝 Test %d: %s\n", i+1, test.title)
		fmt.Printf("Content: %s\n", test.content)

		isClean, err := contentFilter.CheckContent(test.title, test.content)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		result := "INAPPROPRIATE"
		if isClean {
			result = "CLEAN"
		}

		fmt.Printf("🤖 AI Result: %s\n", result)

		if result == test.expected {
			fmt.Printf("✅ PASS - Expected %s, got %s\n", test.expected, result)
		} else {
			fmt.Printf("❌ FAIL - Expected %s, got %s\n", test.expected, result)
		}
	}

	fmt.Println("\n🔄 Testing completed!")
}