package schema

type CommentWithSubs struct {
	*Comment
	Replys []*CommentWithSubs `json:",omitempty"`
}

type ArticleWithComments struct {
	*Article
	Comments []*CommentWithSubs
}
