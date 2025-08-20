package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genai"
)

// Initialize the Gemini client once to reuse across calls
var geminiClient *genai.Client

// InitializeAIClient sets up the Gemini API client
func InitializeAIClient() error {
	ctx := context.Background()
	var err error
	geminiClient, err = genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to initialize Gemini client: %v", err)
	}
	return nil
}

// Example function showing basic client usage
func aiTellMeHowAIWorks() {
	// Ensure client is initialized
	if geminiClient == nil {
		if err := InitializeAIClient(); err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()
	result, err := geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("Explain how AI works in a few words"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
}

// DetailToText converts a detail value to its textual representation from details.json
func DetailToText(detailName string, value int) string {
	// This should be populated from details.json, but for simplicity we handle a few common ones
	switch detailName {
	case "dogs":
		if value >= 0 && value <= 10 {
			descriptions := []string{
				"I tell people I'm allergic, but just to be polite. Dogs are chaos gremlins.",
				"They smell weird and stare at me like I owe them money.",
				"If one jumps on me, I will file a formal complaint—with HR or God.",
				"They're fine. From a distance. Like mountains. Or relatives.",
				"I pet them. I move on. That's it.",
				"As long as he doesn't shit on the floor or bark at ghosts, we're cool.",
				"I might sneak one a piece of cheese. Maybe.",
				"I wave at dogs from my car. That's where I'm at in life.",
				"I have dog biscuits in my pocket right now. Just in case.",
				"I show people photos of my dog before my family.",
				"I love dogs more than people. I'd run into traffic for a stranger's golden retriever.",
			}
			return descriptions[value]
		}
	}
	// For other details, just return the numeric value
	return strconv.Itoa(value)
}

// formatDetailsForPrompt creates a readable summary of a person's details
func formatDetailsForPrompt(person Person) string {
	var details strings.Builder

	details.WriteString("Details:\n")

	// Add dogs preference if available
	if person.Details.Dogs != 0 {
		details.WriteString(fmt.Sprintf("- Dogs: %d/10 - \"%s\"\n",
			person.Details.Dogs, DetailToText("dogs", person.Details.Dogs)))
	}

	// Add cats preference if available
	if person.Details.Cats != 0 {
		details.WriteString(fmt.Sprintf("- Cats: %d/10\n", person.Details.Cats))
	}

	// Add kids preference if available
	if person.Details.Kids != 0 {
		details.WriteString(fmt.Sprintf("- Kids: %d/10\n", person.Details.Kids))
	}

	// Add smoking preference if available
	if person.Details.Smoking != 0 {
		details.WriteString(fmt.Sprintf("- Smoking: %d/10\n", person.Details.Smoking))
	}

	// Add drinking preference if available
	if person.Details.Drinking != 0 {
		details.WriteString(fmt.Sprintf("- Drinking: %d/10\n", person.Details.Drinking))
	}

	// Add religion preference if available
	if person.Details.Religion != 0 {
		details.WriteString(fmt.Sprintf("- Religion: %d/10\n", person.Details.Religion))
	}

	// Add food preference if available
	if person.Details.Food != 0 {
		details.WriteString(fmt.Sprintf("- Food: %d/10\n", person.Details.Food))
	}

	// Add energy level if available
	if person.Details.EnergyLevel != 0 {
		details.WriteString(fmt.Sprintf("- Energy Level: %d/10\n", person.Details.EnergyLevel))
	}

	// Add outdoorsiness if available
	if person.Details.Outdoorsyness != 0 {
		details.WriteString(fmt.Sprintf("- Outdoorsyness: %d/10\n", person.Details.Outdoorsyness))
	}

	// Add travel preference if available
	if person.Details.Travel != 0 {
		details.WriteString(fmt.Sprintf("- Travel: %d/10\n", person.Details.Travel))
	}

	// Add bougieness if available
	if person.Details.Bouginess != 0 {
		details.WriteString(fmt.Sprintf("- Bougieness: %d/10\n", person.Details.Bouginess))
	}

	// Add political importance if available
	if person.Details.ImportanceOfPolitics != 0 {
		details.WriteString(fmt.Sprintf("- Importance of Politics: %d/10\n", person.Details.ImportanceOfPolitics))
	}

	return details.String()
}

// FormatPersonForPrompt formats person data for the AI prompt
func FormatPersonForPrompt(person Person) string {
	var personStr strings.Builder

	personStr.WriteString(fmt.Sprintf("Name: %s\n", person.Name))
	personStr.WriteString(fmt.Sprintf("Age: %d\n", person.Age))
	personStr.WriteString(fmt.Sprintf("Motto: %s\n", person.Motto))
	personStr.WriteString(fmt.Sprintf("Location: %.6f, %.6f\n", person.LatLocation, person.LongLocation))

	// Format description - if it's JSON (Delta format), just mention it's a rich text description
	if strings.HasPrefix(person.Description, "{") && strings.Contains(person.Description, "ops") {
		personStr.WriteString("Description: [Rich text description available]\n")
	} else if person.Description != "" {
		personStr.WriteString(fmt.Sprintf("Description: %s\n", person.Description))
	}

	personStr.WriteString(formatDetailsForPrompt(person))

	return personStr.String()
}

// GenerateChatIntroduction creates an introduction message for a new match
func GenerateChatIntroduction(match Match) (string, error) {
	// Ensure client is initialized
	if geminiClient == nil {
		if err := InitializeAIClient(); err != nil {
			return "", err
		}
	}

	// Format both profiles for the prompt
	person1 := FormatPersonForPrompt(match.OfferedProfile)
	person2 := FormatPersonForPrompt(match.AcceptedProfile)

	// Build the prompt
	prompt := fmt.Sprintf(`You are an AI matchmaker for a dating app. You need to create a friendly, fun introduction message for two people who just matched.

PERSON 1:
%s

PERSON 2:
%s

Based on their profiles, create a message that:
1. Introduces them to each other in a lighthearted, friendly way
2. Highlights what they might have in common or complementary traits
3. Suggests a fun ice breaker question related to their profiles 
4. Keeps the tone casual, positive, and fun

Your response should be around 5-10 sentences, direct, and conversation-starting. Do not mention their exact scores, but please do reference why they might be a good match for each other.

IMPORTANT: Format your response as if you're the dating app sending the first message in their chat. Don't include any meta commentary or notes to me.`, person1, person2)

	// Call Gemini API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate introduction: %v", err)
	}

	return result.Text(), nil
}

// CreateInitialChatMessage creates the first message in a match's chat history
func CreateInitialChatMessage(match Match) (*ChatMessage, error) {
	// Generate introduction text
	introText, err := GenerateChatIntroduction(match)
	if err != nil {
		return nil, err
	}

	// Create a system message (ID 0 indicates system message)
	message := &ChatMessage{
		MatchID: int(match.ID),
		Time:    time.Now(),
		Who:     0, // 0 indicates system/AI message
		Message: introText,
	}

	return message, nil
}

// Called when a chat is started for the first time
func GenerateMatchIntroduction(matchID uint) (*ChatMessage, error) {
	// Find the match
	var match Match
	result := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find match: %v", result.Error)
	}

	// Generate the initial chat message
	message, err := CreateInitialChatMessage(match)
	if err != nil {
		return nil, err
	}

	// Save the message to the database
	if err := db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to save introduction message: %v", err)
	}

	return message, nil
}

// GenerateVibeChat creates a new conversation starter based on recent messages
func GenerateVibeChat(matchID uint, recentMessages []ChatMessage) (*ChatMessage, error) {
	// Find the match
	var match Match
	dbResult := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("failed to find match: %v", dbResult.Error)
	}

	// Format profiles for the prompt
	person1 := FormatPersonForPrompt(match.OfferedProfile)
	person2 := FormatPersonForPrompt(match.AcceptedProfile)

	// Format recent messages for the prompt
	var messagesForPrompt strings.Builder
	messagesForPrompt.WriteString("Recent messages (from oldest to newest):\n")
	
	// Determine the number of messages to include
	messageCount := len(recentMessages)
	if messageCount > 15 {
		// Only include the 15 most recent messages
		recentMessages = recentMessages[messageCount-15:]
	}

	// Format each message
	for i, msg := range recentMessages {
		var sender string
		if msg.Who == match.Offered {
			sender = match.OfferedProfile.Name
		} else if msg.Who == match.Accepted {
			sender = match.AcceptedProfile.Name
		} else if msg.Who == 0 {
			sender = "AI Host"
		} else {
			sender = fmt.Sprintf("Unknown (%d)", msg.Who)
		}
		
		messagesForPrompt.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, sender, msg.Message))
	}

	// Build the prompt for Gemini
	prompt := fmt.Sprintf(`You are an AI conversation assistant for a dating app. You need to create a fresh, interesting conversation starter for two people who are chatting.

PERSON 1:
%s

PERSON 2:
%s

%s

Based on their profiles and recent conversation, create a message that:
1. Is friendly, thoughtful, and engaging
2. References their interests and compatibility points
3. Introduces a new topic if the conversation has gone stale or steered into sensitive territory
4. Asks a specific, interesting question that's easy to respond to
5. Keeps the tone casual, positive, and fun

Your response should be 2-4 sentences maximum, direct, and conversation-starting. Be genuine, warm, and helpful.

IMPORTANT: Format your response as if you're the dating app's AI host sending a message to both users. Don't include any meta commentary or notes. Don't mention that you are an AI unless relevant to the conversation.`, person1, person2, messagesForPrompt.String())

	// Ensure client is initialized
	if geminiClient == nil {
		if err := InitializeAIClient(); err != nil {
			return nil, fmt.Errorf("failed to initialize AI client: %v", err)
		}
	}

	// Call Gemini API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get response from Gemini
	aiResponse, err := geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate vibe chat message: %v", err)
	}
	
	// Get the generated text from the response
	generatedText := aiResponse.Text()
	
	// Create the chat message
	message := &ChatMessage{
		MatchID: int(matchID),
		Time:    time.Now(),
		Who:     0, // 0 indicates system/AI message
		Message: generatedText,
	}

	return message, nil
}

// GenerateDateSuggestion creates a personalized date recommendation based on profiles and location
func GenerateDateSuggestion(matchID uint, recentMessages []ChatMessage) (*ChatMessage, error) {
	// Find the match
	var match Match
	dbResult := db.Preload("OfferedProfile").Preload("AcceptedProfile").First(&match, matchID)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("failed to find match: %v", dbResult.Error)
	}

	// Format profiles for the prompt
	person1 := FormatPersonForPrompt(match.OfferedProfile)
	person2 := FormatPersonForPrompt(match.AcceptedProfile)
	
	// Format locations for the prompt
	location1 := fmt.Sprintf("%.6f,%.6f", match.OfferedProfile.LatLocation, match.OfferedProfile.LongLocation)
	location2 := fmt.Sprintf("%.6f,%.6f", match.AcceptedProfile.LatLocation, match.AcceptedProfile.LongLocation)
	
	// Calculate midpoint between users (approximation)
	midLat := (match.OfferedProfile.LatLocation + match.AcceptedProfile.LatLocation) / 2
	midLong := (match.OfferedProfile.LongLocation + match.AcceptedProfile.LongLocation) / 2
	midpoint := fmt.Sprintf("%.6f,%.6f", midLat, midLong)
	
	// Format recent messages for the prompt
	var messagesForPrompt strings.Builder
	messagesForPrompt.WriteString("Recent messages (from oldest to newest):\n")
	
	// Determine the number of messages to include
	messageCount := len(recentMessages)
	if messageCount > 15 {
		// Only include the 15 most recent messages
		recentMessages = recentMessages[messageCount-15:]
	}

	// Format each message
	for i, msg := range recentMessages {
		var sender string
		if msg.Who == match.Offered {
			sender = match.OfferedProfile.Name
		} else if msg.Who == match.Accepted {
			sender = match.AcceptedProfile.Name
		} else if msg.Who == 0 {
			sender = "AI Host"
		} else {
			sender = fmt.Sprintf("Unknown (%d)", msg.Who)
		}
		
		messagesForPrompt.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, sender, msg.Message))
	}

	// Build the prompt for Gemini
	prompt := fmt.Sprintf(`You are an AI dating assistant for a dating app. You need to suggest a perfect date for two people who have matched.

PERSON 1 PROFILE:
%s
PERSON 1 LOCATION: %s

PERSON 2 PROFILE:
%s
PERSON 2 LOCATION: %s

APPROXIMATE MIDPOINT BETWEEN USERS: %s

RECENT CONVERSATION:
%s

Create a personalized date suggestion that:
1. Recommends a specific activity based on their shared interests and preferences
2. Suggests a specific location for their date that's convenient for both (use the midpoint as reference)
3. Considers their energy levels, outdoorsiness, and preferences mentioned in their profiles
4. References something from their conversation if relevant
5. Is detailed enough to be actionable (specific venue names, activities)
6. If the users have discussed days, dates, or times, please incorporate that in your recommendation
7. Is thoughtful and considers both users' preferences equally
8. If there is something specific about the venue or activity that the users might find appealing, use your search functionality to confirm the details prior to your recommendation
9. If you recommend a location rather than an activity, use your search functionality to check whether there are any special events that coincide with the users timeline that they might enjoy.

IMPORTANT: Do NOT make up places or events. If you recommend a specific venue or location, be clear that the users should verify it exists and check opening hours/availability. Use your web search tooling to check that it appears to exist, what the hours are today (or for a recommended date), and any other details.

If the users are in different cities or far apart, acknowledge this and suggest options for a virtual date or meeting halfway.

Your response should be 3-15 sentences, warm and helpful. Format as a message from the dating app's AI host.`, 
		person1, location1, person2, location2, midpoint, messagesForPrompt.String())

	// Ensure client is initialized
	if geminiClient == nil {
		if err := InitializeAIClient(); err != nil {
			return nil, fmt.Errorf("failed to initialize AI client: %v", err)
		}
	}

	// Call Gemini API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get response from Gemini
	aiResponse, err := geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate date suggestion: %v", err)
	}
	
	// Get the generated text from the response
	generatedText := aiResponse.Text()
	
	// Create the chat message
	message := &ChatMessage{
		MatchID: int(matchID),
		Time:    time.Now(),
		Who:     0, // 0 indicates system/AI message
		Message: generatedText,
	}

	return message, nil
}