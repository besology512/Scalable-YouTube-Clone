package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Video struct {
	MongoID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ID               string             `bson:"id" json:"id"`
	Title            string             `bson:"title" json:"title"`
	Tags             []string           `bson:"tags" json:"tags"`
	Categories       []string           `bson:"categories" json:"categories"`
	AgeGroup         string             `bson:"ageGroup" json:"ageGroup"` // "teen", "adult", "senior"
	Gender           string             `bson:"gender" json:"gender"`     // target gender, "all" or "male"/"female"
	Views            int                `bson:"views" json:"views"`
	Likes            int                `bson:"likes" json:"likes"`
	Dislikes         int                `bson:"dislikes" json:"dislikes"`
	Duration         int                `bson:"duration" json:"duration"`     // in seconds
	UploadDate       string             `bson:"uploadDate" json:"uploadDate"` // e.g., "2023-10-01"
	Language         string             `bson:"language" json:"language"`     // e.g., "English", "Spanish"
	Description      string             `bson:"description" json:"description"`
	Thumbnail        string             `bson:"thumbnail" json:"thumbnail"`               // URL to thumbnail image
	ChannelID        string             `bson:"channelId" json:"channelId"`               // ID of the channel that uploaded the video
	ChannelName      string             `bson:"channelName" json:"channelName"`           // name of the channel
	Rating           float64            `bson:"rating" json:"rating"`                     // average rating from users
	Comments         []string           `bson:"comments" json:"comments"`                 // list of comments
	TagsScore        int                `bson:"tagsScore" json:"tagsScore"`               // score based on tags relevance
	CategoryScore    int                `bson:"categoryScore" json:"categoryScore"`       // score based on category relevance
	PopularityScore  int                `bson:"popularityScore" json:"popularityScore"`   // score based on views and likes
	EngagementScore  int                `bson:"engagementScore" json:"engagementScore"`   // score based on likes, dislikes, and comments
	WatchTime        int                `bson:"watchTime" json:"watchTime"`               // total watch time in seconds
	Region           string             `bson:"region" json:"region"`                     // target region, e.g., "US", "EU"
	TagsBoost        map[string]int     `bson:"tagsBoost" json:"tagsBoost"`               // boost score for specific tags
	CategoriesBoost  map[string]int     `bson:"categoriesBoost" json:"categoriesBoost"`   // boost score for specific categories
	AgeGroupBoost    map[string]int     `bson:"ageGroupBoost" json:"ageGroupBoost"`       // boost score for specific age groups
	LanguageBoost    map[string]int     `bson:"languageBoost" json:"languageBoost"`       // boost score for specific languages
	GenderBoost      map[string]int     `bson:"genderBoost" json:"genderBoost"`           // boost
	IsFamilyFriendly bool               `bson:"isFamilyFriendly" json:"isFamilyFriendly"` // true if the video is family-friendly
}
