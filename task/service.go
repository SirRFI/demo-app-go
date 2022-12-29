package task

type taskRepository interface {
	List() ([]Task, error)
	GetByID(id ID) (Task, error)
	Add(addTask AddTaskCommand) (Task, error)
	Save(task Task) error
	Delete(id ID) error
}

//type Service struct {
//	repository taskRepository
//}
//
//func (s *Service) ListTasks() ([]Task, error) {
//
//}
//
//func (s *Service) GetTask(id ID) (Task, error) {
//
//}
//
//func (s *Service) AddTask(command AddTaskCommand) (Task, error) {
//
//}

//func (s *Service)

//func (s *Service)
