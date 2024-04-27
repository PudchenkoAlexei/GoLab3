package lang

import (
	"image"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type State struct {
	backgroundColor     painter.Operation
	backgroundRectangle *painter.BgRect
	figuresArray        []*painter.Figure
	moveOperations      []painter.Operation
	updateOperation     painter.Operation
}

func (s *State) Reset() {
	s.backgroundColor = nil
	s.backgroundRectangle = nil
	s.figuresArray = nil
	s.moveOperations = nil
	s.updateOperation = nil
}

func (s *State) GetOperations() []painter.Operation {
	var ops []painter.Operation

	if s.backgroundColor != nil {
		ops = append(ops, s.backgroundColor)
	}
	if s.backgroundRectangle != nil {
		ops = append(ops, s.backgroundRectangle)
	}
	if len(s.moveOperations) != 0 {
		ops = append(ops, s.moveOperations...)
		s.moveOperations = nil
	}
	if len(s.figuresArray) != 0 {
		for _, figure := range s.figuresArray {
			ops = append(ops, figure)
		}
	}
	if s.updateOperation != nil {
		ops = append(ops, s.updateOperation)
	}

	return ops
}

func (s *State) ResetOperations() {
	if s.backgroundColor == nil {
		s.backgroundColor = painter.OperationFunc(painter.Reset)
	}
	if s.updateOperation != nil {
		s.updateOperation = nil
	}
}

func (s *State) GreenBackground() {
	s.backgroundColor = painter.OperationFunc(painter.GreenFill)
}

func (s *State) WhiteBackground() {
	s.backgroundColor = painter.OperationFunc(painter.WhiteFill)
}

func (s *State) BackgroundRectangle(firstPoint image.Point, secondPoint image.Point) {
	s.backgroundRectangle = &painter.BgRect{
		FirstPoint:  firstPoint,
		SecondPoint: secondPoint,
	}
}

func (s *State) AddFigure(centralPoint image.Point) {
	figure := painter.Figure{
		CentralPoint: centralPoint,
	}
	s.figuresArray = append(s.figuresArray, &figure)
}

func (s *State) AddMoveOperation(x int, y int) {
	moveOp := painter.Move{X: x, Y: y, FiguresArray: s.figuresArray}
	s.moveOperations = append(s.moveOperations, &moveOp)
}

func (s *State) ResetStateAndBackground() {
	s.Reset()
	s.backgroundColor = painter.OperationFunc(painter.Reset)
}

func (s *State) SetUpdateOperation() {
	s.updateOperation = painter.UpdateOp
}
