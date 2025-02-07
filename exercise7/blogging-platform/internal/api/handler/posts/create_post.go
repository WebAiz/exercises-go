package posts

import (
	"fmt"
	"net/http"

	"github.com/webaiz/exercise7/blogging-platform/internal/auth"
	"github.com/webaiz/exercise7/blogging-platform/internal/db/post"
	"github.com/webaiz/exercise7/blogging-platform/pkg/httputils/request"
	"github.com/webaiz/exercise7/blogging-platform/pkg/httputils/response"
)

type CreatePostRequest struct {
	Data *post.ModelPost `json:"data"`
}

type CreatePostResponse struct {
	Data *post.ModelPost `json:"data"`
}

func (h *Posts) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.logger.With("method", "CreatePost")

	user, ok := ctx.Value("user").(*auth.UserData)
	if !ok {
		log.ErrorContext(
			ctx,
			"failed to type cast user data",
		)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("user: %+v\n", *user)

	// request parse
	requestBody := &CreatePostRequest{}

	if err := request.JSON(w, r, requestBody); err != nil {
		log.ErrorContext(
			ctx,
			"failed to parse request body",
			"error", err,
		)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	// db request
	dbResp, err := h.db.CreatePost(ctx, requestBody.Data)

	if err != nil {
		log.ErrorContext(
			ctx,
			"failed to query from db",
			"error", err,
		)
		http.Error(w, "failed to query from db", http.StatusInternalServerError)
		return
	}

	if dbResp == nil {
		log.ErrorContext(
			ctx,
			"row is empty",
		)
		http.Error(w, "row is empty", http.StatusInternalServerError)
		return
	}

	// response
	resp := CreatePostResponse{
		Data: dbResp,
	}

	if err := response.JSON(
		w,
		http.StatusOK,
		resp,
	); err != nil {
		log.ErrorContext(
			ctx,
			"fail json",
			"error", err,
		)
		return
	}

	log.InfoContext(
		ctx,
		"success insert post",
		"post id", resp.Data.ID,
	)
	return
}
