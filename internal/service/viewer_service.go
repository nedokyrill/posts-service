package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/pkg/consts"
	"github.com/nedokyrill/posts-service/pkg/logger"
	"github.com/nedokyrill/posts-service/pkg/utils"
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
func (s *ViewerServiceImpl) CreateViewer(_ context.Context, postId uuid.UUID) (int, chan *models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comm := make(chan *models.Comment)
	s.viewers[postId] = append(s.viewers[postId], Viewer{ch: comm, id: s.cnt})
	s.cnt++

	logger.Logger.Infof("create viewer for post id %s", postId.String())
	return s.cnt, comm, nil
}

// удаляем подписчика из пула при закрытии подписки
func (s *ViewerServiceImpl) DeleteViewer(_ context.Context, postId uuid.UUID, id int) error {
	s.mu.Lock()

	viewers, ok := s.viewers[postId]
	if !ok {
		s.mu.Unlock()
		logger.Logger.Error(fmt.Sprintf("no post with postId: %s", postId.String()))
		return utils.GqlError{Msg: fmt.Sprintf("no post with postId: %s", postId.String()),
			Type: consts.BadRequestType}
	}
	for i, viewer := range viewers {
		if viewer.id == id {
			s.viewers[postId] = append(viewers[:i], viewers[i+1:]...)
			close(viewer.ch)
			break
		}
	}

	s.mu.Unlock()

	logger.Logger.Infof("delete viewer for post id %s", postId.String())
	return nil
}

// отправляем уведомление в виде комментария всем подписчикам
func (s *ViewerServiceImpl) NotifyViewers(_ context.Context, postId uuid.UUID, comm models.Comment) error {
	s.mu.Lock()

	viewers, ok := s.viewers[postId]
	if !ok { // подписчиков на данный пост нет
		s.mu.Unlock()
		return nil
	}

	snap := make([]Viewer, len(viewers))
	copy(snap, viewers)
	s.mu.Unlock() // сделали копию каналов и разлочили мьютекс, чтобы далее не блокировать данные

	for _, v := range snap {
		v.ch <- &comm
	}

	logger.Logger.Info(fmt.Sprintf("notify viewers for postId %s successfully", postId.String()))
	return nil
}
