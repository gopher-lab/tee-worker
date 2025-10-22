package types_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/masa-finance/tee-worker/api/types"
)

var _ = Describe("RedditResponse", func() {
	Describe("Unmarshalling", func() {
		It("should unmarshal a user response", func() {
			jsonData := `{"type": "user", "id": "user123", "username": "testuser"}`
			var resp types.RedditResponse
			err := json.Unmarshal([]byte(jsonData), &resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.User).ToNot(BeNil())
			Expect(resp.Post).To(BeNil())
			Expect(resp.User.ID).To(Equal("user123"))
			Expect(resp.User.Username).To(Equal("testuser"))
		})

		It("should unmarshal a post response", func() {
			jsonData := `{"type": "post", "id": "post123", "title": "Test Post"}`
			var resp types.RedditResponse
			err := json.Unmarshal([]byte(jsonData), &resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Post).ToNot(BeNil())
			Expect(resp.User).To(BeNil())
			Expect(resp.Post.ID).To(Equal("post123"))
			Expect(resp.Post.Title).To(Equal("Test Post"))
		})

		It("should return an error for an unknown type", func() {
			jsonData := `{"type": "unknown", "id": "123"}`
			var resp types.RedditResponse
			err := json.Unmarshal([]byte(jsonData), &resp)
			Expect(err).To(MatchError("unknown Reddit response type: unknown"))
		})
	})

	Describe("Marshalling", func() {
		It("should marshal a user response", func() {
			now := time.Now()
			resp := types.RedditResponse{
				Type: types.RedditUserItem,
				User: &types.RedditUser{
					ID:           "user123",
					Username:     "testuser",
					CreatedAt:    now,
					CommentKarma: 10,
				},
			}

			expectedJSON, err := json.Marshal(resp.User)
			Expect(err).ToNot(HaveOccurred())

			actualJSON, err := json.Marshal(&resp)
			Expect(err).ToNot(HaveOccurred())

			Expect(actualJSON).To(MatchJSON(expectedJSON))
		})

		It("should marshal a post response", func() {
			resp := types.RedditResponse{
				Type: types.RedditPostItem,
				Post: &types.RedditPost{
					ID:    "post123",
					Title: "Test Post",
				},
			}

			expectedJSON, err := json.Marshal(resp.Post)
			Expect(err).ToNot(HaveOccurred())

			actualJSON, err := json.Marshal(&resp)
			Expect(err).ToNot(HaveOccurred())

			Expect(actualJSON).To(MatchJSON(expectedJSON))
		})

		It("should return an error for an unknown type", func() {
			resp := types.RedditResponse{
				Type: "unknown",
			}
			_, err := json.Marshal(&resp)
			Expect(err).To(HaveOccurred())
		})
	})
})
