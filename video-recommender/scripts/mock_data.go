// Run manually to insert test data
package main

import (
	"context"
	"video-recommender/config"
	"video-recommender/models"
)

func main() {
	config.ConnectDB()

	// Insert diverse users
	users := []models.User{
		{
			ID:                  "user1",
			Interests:           []string{"tech", "gaming"},
			WatchHistory:        []string{"video1", "video3", "video5"},
			LikedVideos:         []string{"video1", "video2"},
			PreferredCategories: []string{"tech", "education"},
			Age:                 25,
			Gender:              "male",
			Location:            "New York, USA",
			Device:              "mobile",
			SubscriptionStatus:  "premium",
			NotificationPrefs:   map[string]bool{"email": true, "sms": false},
			LanguagePreference:  []string{"English"},
			WatchTime:           1200,
			DeviceHistory:       []string{"mobile", "desktop"},
			Feedback:            map[string]string{"video1": "liked", "video3": "disliked"},
			SearchHistory:       []string{"AI tutorials", "gaming reviews"},
			Recommendations:     []string{"video6", "video7"},
			BlockedUsers:        []string{"user5"},
			CustomPlaylists:     map[string][]string{"Favorites": {"video1", "video2"}},
			ParentalControl:     false,
			WatchLater:          []string{"video8"},
			OfflineDownloads:    []string{"video3"},
			DeviceSettings:      map[string]string{"mobile": "dark", "desktop": "light"},
		},
		{
			ID:                  "user2",
			Interests:           []string{"music", "sports"},
			WatchHistory:        []string{"video10", "video12"},
			LikedVideos:         []string{"video10"},
			PreferredCategories: []string{"entertainment", "sports"},
			Age:                 30,
			Gender:              "female",
			Location:            "London, UK",
			Device:              "desktop",
			SubscriptionStatus:  "free",
			NotificationPrefs:   map[string]bool{"email": false, "sms": true},
			LanguagePreference:  []string{"English", "Spanish"},
			WatchTime:           800,
			DeviceHistory:       []string{"desktop"},
			Feedback:            map[string]string{"video10": "liked"},
			SearchHistory:       []string{"concert highlights", "football matches"},
			Recommendations:     []string{"video15", "video16"},
			BlockedUsers:        []string{},
			CustomPlaylists:     map[string][]string{"Workout": {"video10", "video12"}},
			ParentalControl:     false,
			WatchLater:          []string{"video20"},
			OfflineDownloads:    []string{},
			DeviceSettings:      map[string]string{"desktop": "light"},
		},
	}

	for _, user := range users {
		config.DB.Collection("users").InsertOne(context.Background(), user)
	}

	// Insert diverse videos
	videos := []models.Video{
		{
			ID:               "video1",
			Title:            "Introduction to AI",
			Tags:             []string{"tech", "AI"},
			Categories:       []string{"education", "tech"},
			AgeGroup:         "adult",
			Gender:           "all",
			Views:            5000,
			Likes:            400,
			Dislikes:         20,
			Duration:         300,
			UploadDate:       "2023-01-15",
			Language:         "English",
			Description:      "Learn the basics of Artificial Intelligence.",
			Thumbnail:        "https://example.com/thumbnails/ai.jpg",
			ChannelID:        "channel1",
			ChannelName:      "Tech Guru",
			Rating:           4.5,
			Comments:         []string{"Great video!", "Very informative."},
			TagsScore:        90,
			CategoryScore:    85,
			PopularityScore:  88,
			EngagementScore:  92,
			WatchTime:        15000,
			Region:           "US",
			TagsBoost:        map[string]int{"tech": 10, "AI": 8},
			CategoriesBoost:  map[string]int{"education": 9, "tech": 10},
			AgeGroupBoost:    map[string]int{"adult": 10},
			LanguageBoost:    map[string]int{"English": 10},
			GenderBoost:      map[string]int{"all": 10},
			IsFamilyFriendly: true,
		},
		{
			ID:               "video10",
			Title:            "Top 10 Football Goals",
			Tags:             []string{"sports", "football"},
			Categories:       []string{"sports", "entertainment"},
			AgeGroup:         "teen",
			Gender:           "all",
			Views:            20000,
			Likes:            1500,
			Dislikes:         50,
			Duration:         600,
			UploadDate:       "2023-02-10",
			Language:         "English",
			Description:      "A compilation of the top 10 football goals.",
			Thumbnail:        "https://example.com/thumbnails/football.jpg",
			ChannelID:        "channel2",
			ChannelName:      "Sports Central",
			Rating:           4.8,
			Comments:         []string{"Amazing goals!", "Loved it."},
			TagsScore:        95,
			CategoryScore:    90,
			PopularityScore:  93,
			EngagementScore:  96,
			WatchTime:        30000,
			Region:           "EU",
			TagsBoost:        map[string]int{"sports": 10, "football": 9},
			CategoriesBoost:  map[string]int{"sports": 10, "entertainment": 8},
			AgeGroupBoost:    map[string]int{"teen": 10},
			LanguageBoost:    map[string]int{"English": 10},
			GenderBoost:      map[string]int{"all": 10},
			IsFamilyFriendly: true,
		},
	}

	for _, video := range videos {
		config.DB.Collection("videos").InsertOne(context.Background(), video)
	}
}
