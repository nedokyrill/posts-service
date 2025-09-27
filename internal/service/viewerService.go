package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/pkg/logger"
)

type Viewer struct {
	ch chan *models.Comment
	id int
}

type ViewerServiceImpl struct {
	viewers map[uuid.UUID][]Viewer
	cnt     int
	mu      sync.Mutex
}

func NewViewerService() *ViewerServiceImpl {
	return &ViewerServiceImpl{
		viewers: make(map[uuid.UUID][]Viewer),
		cnt:     0,
		mu:      sync.Mutex{},
	}
}

// добавляем подписчика
func (s *ViewerServiceImpl) CreateViewer(postId uuid.UUID) (int, chan *models.Comment, error) {
	s.mu.Lock()

	comm := make(chan *models.Comment)
	s.viewers[postId] = append(s.viewers[postId], Viewer{ch: comm, id: s.cnt})
	s.cnt++

	s.mu.Unlock()

	logger.Logger.Infof("create viewer for post id %s", postId.String())
	return s.cnt, comm, nil
}

// удаляем подписчика из пула при закрытии подписки
func (s *ViewerServiceImpl) DeleteViewer(postId uuid.UUID, chanId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	viewers, ok := s.viewers[postId]
	if !ok {
		logger.Logger.Error(fmt.Sprintf("no post with postId: %s", postId.String()))
		return errors.New(fmt.Sprintf("no post with postId: %s", postId.String()))
	}
	for i, viewer := range viewers {
		if viewer.id == chanId {
			s.viewers[postId] = append(viewers[:i], viewers[i+1:]...)
		}
	}

	logger.Logger.Infof("delete viewer for post id %s", postId.String())
	return nil
}

// отправляем уведомление в виде комментария всем подписчикам
func (s *ViewerServiceImpl) NotifyViewers(postId uuid.UUID, comm models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	viewers, ok := s.viewers[postId]
	if !ok {
		logger.Logger.Error(fmt.Sprintf("no post with postId: %s", postId.String()))
		return errors.New(fmt.Sprintf("no post with postId: %s", postId.String()))
	}

	for _, viewer := range viewers {
		viewer.ch <- &comm
	}

	logger.Logger.Info(fmt.Sprintf("notify viewer for postId %s successfully", postId.String()))
	return nil
}
