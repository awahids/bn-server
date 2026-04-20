package bookmarkreq

type CreateBookmarkRequest struct {
	Type      string  `json:"type"`
	ContentID string  `json:"contentId"`
	Note      *string `json:"note"`
}
