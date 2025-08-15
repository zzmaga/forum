package service

import (
	"errors"
	"fmt"
	"log"

	"forum/architecture/models"
	"forum/architecture/service/post_vote"
)

func (s *Service) FillPost(post *models.Post, sesUserId int64) error {
	var err error

	post.WCategories, err = s.Category.GetByPostID(post.Id)
	switch {
	case err != nil:
		log.Printf("FillPost: PostCategory.GetByPostID(postId: %v): %v", post.Id, err)
	}

	post.WUser, err = s.User.GetByID(post.UserId)
	switch {
	case err != nil:
		log.Printf("FillPost: User.GetByID(userId: %v): %v", post.UserId, err)
	}

	vUp, vDown, err := s.PostVote.GetByPostID(post.Id)
	switch {
	case err != nil:
		log.Printf("FillPost: PostVote.GetByPostID(id: %v): %v", post.Id, err)
	}
	post.WVoteUp = vUp
	post.WVoteDown = vDown

	if sesUserId == 0 {
		return nil
	}

	vUser, err := s.PostVote.GetPostUserVote(sesUserId, post.Id)
	switch {
	case err == nil:
		post.WUserVote = vUser.Vote
	case errors.Is(err, post_vote.ErrNotFound):
	case err != nil:
		log.Printf("FillPost: PostVote.GetPostUserVote(userId: %v, postId: %v): %v", sesUserId, post.Id, err)
	}
	return nil
}

func (s *Service) FillPosts(posts []*models.Post, sesUserId int64) error {
	for _, post := range posts {
		err := s.FillPost(post, sesUserId)
		if err != nil {
			return fmt.Errorf("FillPosts: %w", err)
		}
	}
	return nil
}
