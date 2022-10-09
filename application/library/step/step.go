package step

import "github.com/webx-top/echo"

// NewGroup 创建新的步骤集
func NewGroup() *Steps {
	return &Steps{Detail: map[string]*Step{}}
}

// New 创建新的步骤
func New(stepName string, titleAndDescription ...string) *Step {
	var (
		title       string
		description string
	)
	switch len(titleAndDescription) {
	case 2:
		description = titleAndDescription[1]
	case 1:
		title = titleAndDescription[0]
	}
	return &Step{Name: stepName, Title: title, Description: description, Config: echo.H{}}
}

// Processor 步骤处理逻辑
type Processor func(echo.Context) error

// Step 单个步骤信息
type Step struct {
	Name        string // 英文标识名
	Title       string // 步骤标题
	Description string // 步骤说明
	Active      bool   // 是否已完成步骤
	Index       int
	Config      echo.H
	processors  []Processor
}

func (s *Step) SetTitle(title string) *Step {
	s.Title = title
	return s
}

func (s *Step) SetDescription(description string) *Step {
	s.Description = description
	return s
}

func (s *Step) SetActive(active bool) *Step {
	s.Active = active
	return s
}

func (s *Step) SetActiveByCurrentIndex(currentIndex int) *Step {
	s.Active = currentIndex >= s.Index
	return s
}

func (s *Step) SetConfig(config echo.H) *Step {
	s.Config = config
	return s
}

func (s *Step) SetConfigKV(key string, value interface{}) *Step {
	s.Config.Set(key, value)
	return s
}

// Next 下一步
func (s *Step) Next(ss *Steps) *Step {
	return ss.Next(s.Index)
}

// Prev 上一步
func (s *Step) Prev(ss *Steps) *Step {
	return ss.Prev(s.Index)
}

// Process 执行处理逻辑
func (s *Step) Process(ctx echo.Context) error {
	for _, processor := range s.processors {
		if err := processor(ctx); err != nil {
			return err
		}
	}
	return nil
}

// AddProcessor 添加处理逻辑
func (s *Step) AddProcessor(processors ...Processor) *Step {
	s.processors = append(s.processors, processors...)
	return s
}

// Steps 步骤集合
type Steps struct {
	Idents []string
	Detail map[string]*Step
}

// Add 添加步骤
func (s *Steps) Add(stepName string, processors ...Processor) *Steps {
	if step, found := s.Detail[stepName]; !found {
		s.Detail[stepName] = &Step{
			Name:       stepName,
			Index:      len(s.Idents),
			Config:     echo.H{},
			processors: processors,
		}
		s.Idents = append(s.Idents, stepName)
	} else {
		step.processors = append(step.processors, processors...)
	}
	return s
}

// AddStep 添加步骤
func (s *Steps) AddStep(step *Step) *Steps {
	_, found := s.Detail[step.Name]
	if found {
		panic(`The step already exists: ` + step.Name)
	}
	step.Index = len(s.Idents)
	s.Detail[step.Name] = step
	s.Idents = append(s.Idents, step.Name)
	return s
}

// Set 设置步骤
func (s *Steps) Set(stepName string, processors ...Processor) *Steps {
	if step, found := s.Detail[stepName]; !found {
		s.Detail[stepName] = &Step{
			Name:       stepName,
			Index:      len(s.Idents),
			Config:     echo.H{},
			processors: processors,
		}
		s.Idents = append(s.Idents, stepName)
	} else {
		step.processors = processors
	}
	return s
}

// SetStep 设置步骤
func (s *Steps) SetStep(step *Step) *Steps {
	if old, found := s.Detail[step.Name]; !found {
		step.Index = len(s.Idents)
		s.Detail[step.Name] = step
		s.Idents = append(s.Idents, step.Name)
	} else {
		step.Index = old.Index
		s.Detail[step.Name] = step
	}
	return s
}

// Get 根据名称获取
func (s *Steps) Get(stepName string) *Step {
	if step, found := s.Detail[stepName]; found {
		return step
	}
	return nil
}

// GetAny 根据名称获取
func (s *Steps) GetAny(stepNames ...string) *Step {
	for _, stepName := range stepNames {
		if step, found := s.Detail[stepName]; found {
			return step
		}
	}
	return nil
}

// First 第一步
func (s *Steps) First() *Step {
	if len(s.Idents) == 0 {
		return nil
	}
	return s.Get(s.Idents[0])
}

// Next 下一步
func (s *Steps) Next(index int) *Step {
	next := index + 1
	if len(s.Idents) <= next {
		return nil
	}
	return s.Get(s.Idents[next])
}

// NextByName 下一步
func (s *Steps) NextByName(stepName string) *Step {
	step := s.Get(stepName)
	if step == nil {
		return nil
	}
	return s.Next(step.Index)
}

// Prev 上一步
func (s *Steps) Prev(index int) *Step {
	prev := index - 1
	if prev < 1 || len(s.Idents) <= prev {
		return nil
	}
	return s.Get(s.Idents[prev])
}

// PrevByName 上一步
func (s *Steps) PrevByName(stepName string) *Step {
	step := s.Get(stepName)
	if step == nil {
		return nil
	}
	return s.Prev(step.Index)
}

// Last 最后一步
func (s *Steps) Last() *Step {
	if len(s.Idents) == 0 {
		return nil
	}
	return s.Get(s.Idents[len(s.Idents)-1])
}

// Size 步骤总数
func (s *Steps) Size() int {
	return len(s.Idents)
}

// Slice 按顺序返回所有步骤
func (s *Steps) Slice() []*Step {
	r := make([]*Step, s.Size())
	for index, name := range s.Idents {
		step := s.Get(name)
		if step == nil {
			r[index] = nil
		} else {
			stepCopy := *step
			r[index] = &stepCopy
		}
	}
	return r
}

func (s *Steps) GetTitle(stepName string) string {
	step := s.Get(stepName)
	if step == nil {
		return ``
	}
	return step.Title
}

func (s *Steps) GetDescription(stepName string) string {
	step := s.Get(stepName)
	if step == nil {
		return ``
	}
	return step.Description
}

func (s *Steps) GetIndex(stepName string) int {
	step := s.Get(stepName)
	if step == nil {
		return -1
	}
	return step.Index
}

func (s *Steps) Max() int {
	return s.Size() - 1
}
