package bookmarkrouter

import (
	"github.com/awahids/bn-server/internal/delivery/handlers/bookmarkhandler"

	"github.com/gin-gonic/gin"
)

func RegisterBookmarkRoutes(group *gin.RouterGroup, handler *bookmarkhandler.BookmarkHandler, authMiddleware gin.HandlerFunc) {
	protected := group.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/bookmarks", handler.GetBookmarks)
		protected.POST("/bookmarks", handler.PostBookmark)
		protected.DELETE("/bookmarks/:id", handler.DeleteBookmark)
	}
}
