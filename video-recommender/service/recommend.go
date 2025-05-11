package service

import (
	"fmt"
	"video-recommender/models"
	"video-recommender/repository"
)

func unique(input []string) []string {
	m := map[string]bool{}
	var result []string
	for _, v := range input {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

func RecommendVideos(userID string) ([]models.Video, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Collect tags and categories from user preferences and history
	tags := append([]string{}, user.Interests...)
	categories := append([]string{}, user.PreferredCategories...)

	// Include tags and categories from watched videos
	watched, err := repository.GetVideosByIDs(user.WatchHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to get watched videos: %w", err)
	}
	for _, video := range watched {
		tags = append(tags, video.Tags...)
		categories = append(categories, video.Categories...)
	}

	// Include tags and categories from liked videos
	likedVideos, err := repository.GetVideosByIDs(user.LikedVideos)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked videos: %w", err)
	}
	for _, video := range likedVideos {
		tags = append(tags, video.Tags...)
		categories = append(categories, video.Categories...)
	}

	// Remove duplicates
	tags = unique(tags)
	categories = unique(categories)

	// Fetch recommendations based on various criteria
	recommendedVideos, err := repository.FindVideosByTagsAndCategories(tags, categories, user.WatchHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to find videos by tags and categories: %w", err)
	}
	tagBasedVideos, err := repository.GetVideosByTags(tags)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by tags: %w", err)
	}
	categoryBasedVideos, err := repository.GetVideosByCategories(categories)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by categories: %w", err)
	}
	highRatedVideos, err := repository.GetVideosByRating(4.5)
	if err != nil {
		return nil, fmt.Errorf("failed to get high-rated videos: %w", err)
	}
	recentVideos, err := repository.GetVideosByUploadDate("2023-01-01")
	if err != nil {
		return nil, fmt.Errorf("failed to get recent videos: %w", err)
	}
	watchTimeBasedVideos, err := repository.GetVideosByWatchTime(user.WatchTime / 2)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by watch time: %w", err)
	}
	deviceBasedVideos, err := repository.GetVideosByChannelName(user.Device)
	if err != nil {
		return nil, fmt.Errorf("failed to get device-based videos: %w", err)
	}
	searchHistoryVideos, err := repository.GetVideosByTags(user.SearchHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to get search history videos: %w", err)
	}

	// Include boosted recommendations based on tags and categories
	for _, tag := range tags {
		boostedTagVideos, err := repository.GetVideosByBoostedTags(tag, 10)
		if err != nil {
			return nil, fmt.Errorf("failed to get boosted tag videos: %w", err)
		}
		recommendedVideos = append(recommendedVideos, boostedTagVideos...)
	}
	for _, category := range categories {
		boostedCategoryVideos, err := repository.GetVideosByBoostedCategories(category, 10)
		if err != nil {
			return nil, fmt.Errorf("failed to get boosted category videos: %w", err)
		}
		recommendedVideos = append(recommendedVideos, boostedCategoryVideos...)
	}

	// Combine all recommendations
	recommendedVideos = append(recommendedVideos, tagBasedVideos...)
	recommendedVideos = append(recommendedVideos, categoryBasedVideos...)
	recommendedVideos = append(recommendedVideos, highRatedVideos...)
	recommendedVideos = append(recommendedVideos, recentVideos...)
	recommendedVideos = append(recommendedVideos, watchTimeBasedVideos...)
	recommendedVideos = append(recommendedVideos, deviceBasedVideos...)
	recommendedVideos = append(recommendedVideos, searchHistoryVideos...)

	// Remove duplicates
	recommendedVideos = uniqueVideos(recommendedVideos)

	// Apply parental control filter if enabled
	if user.ParentalControl {
		recommendedVideos = filterParentalControl(recommendedVideos)
	}

	return recommendedVideos, nil
}

// filterParentalControl filters out videos that are not suitable for parental control.
func filterParentalControl(videos []models.Video) []models.Video {
	var filtered []models.Video
	for _, video := range videos {
		if video.IsFamilyFriendly { // Assuming `IsFamilyFriendly` is a field in `models.Video`
			filtered = append(filtered, video)
		}
	}
	return filtered
}

// uniqueVideos removes duplicate videos from the slice.
func uniqueVideos(input []models.Video) []models.Video {
	m := map[string]bool{}
	var result []models.Video
	for _, video := range input {
		if !m[video.ID] {
			m[video.ID] = true
			result = append(result, video)
		}
	}
	return result
}
