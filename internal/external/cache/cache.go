package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	DefaultExpiration              = 5 * time.Minute
	DefaultExpirationNotifications = 90 * time.Second // 1.5 minutes
	DefaultExpirationTopics        = 1 * time.Hour

	// Default keys
	DefaultFeedCacheKey         = "feed."
	DefaultPostDetailCacheKey   = "post_detail."
	DefaultProfileCacheKey      = "profile."
	DefaultNotificationCacheKey = "notifications."
	DefaultBookmarkCacheKey     = "bookmarks."
	DefaultTopicCacheKey        = "topics."
	DefaultTopicFeedCacheKey    = "topics_feed."
	DefaultNotesCacheKey        = "notes."
)

type ICache interface {
	Get(key string) (any, bool)
	Set(key string, element any, duration time.Duration)
	Items() map[string]cache.Item
}

func NewCache(cacheMap map[string]cache.Item) ICache {
	return cache.NewFrom(DefaultExpiration, 2*DefaultExpiration, cacheMap)
}
