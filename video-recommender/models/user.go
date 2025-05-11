package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	MongoID             primitive.ObjectID  `bson:"_id,omitempty" json:"-"`
	ID                  string              `bson:"id" json:"id"`
	Interests           []string            `bson:"interests" json:"interests"`
	WatchHistory        []string            `bson:"watchhistory" json:"watchhistory"`
	LikedVideos         []string            `bson:"likedvideos" json:"likedvideos"`
	PreferredCategories []string            `bson:"preferredcategories" json:"preferredcategories"`
	Age                 int                 `bson:"age" json:"age"`
	Gender              string              `bson:"gender" json:"gender"`                         // "male", "female", etc.
	Location            string              `bson:"location" json:"location"`                     // e.g., "New York, USA"
	Device              string              `bson:"device" json:"device"`                         // e.g., "mobile", "desktop"
	SubscriptionStatus  string              `bson:"subscriptionstatus" json:"subscriptionstatus"` // e.g., "free", "premium"
	NotificationPrefs   map[string]bool     `bson:"notificationprefs" json:"notificationprefs"`   // e.g., {"email": true, "sms": false}
	LanguagePreference  []string            `bson:"languagepreference" json:"languagepreference"` // e.g., "English", "Spanish"
	WatchTime           int                 `bson:"watchtime" json:"watchtime"`                   // total watch time in minutes
	DeviceHistory       []string            `bson:"devicehistory" json:"devicehistory"`           // e.g., ["mobile", "desktop"]
	Feedback            map[string]string   `bson:"feedback" json:"feedback"`                     // e.g., {"video1": "liked", "video2": "disliked"}
	SearchHistory       []string            `bson:"searchhistory" json:"searchhistory"`           // e.g., ["golang", "python"]
	Recommendations     []string            `bson:"recommendations" json:"recommendations"`       // e.g., ["video1", "video2"]
	BlockedUsers        []string            `bson:"blockedusers" json:"blockedusers"`             // e.g., ["user1", "user2"]
	CustomPlaylists     map[string][]string `bson:"customplaylists" json:"customplaylists"`       // e.g., {"Favorites": ["video1", "video2"]}
	ParentalControl     bool                `bson:"parentalcontrol" json:"parentalcontrol"`       // e.g., true or false
	WatchLater          []string            `bson:"watchlater" json:"watchlater"`                 // e.g., ["video1", "video2"]
	OfflineDownloads    []string            `bson:"offlinedownloads" json:"offlinedownloads"`     // e.g., ["video1", "video2"]
	DeviceSettings      map[string]string   `bson:"devicesettings" json:"devicesettings"`         // e.g., {"mobile": "dark", "desktop": "light"}
}
