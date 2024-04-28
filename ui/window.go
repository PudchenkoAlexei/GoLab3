package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz     size.Event
	center image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.center.X = 800 / 2
	pw.center.Y = 800 / 2
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button != mouse.ButtonLeft { // Змінено на ліву кнопку
			break
		}
		if e.Direction != mouse.DirPress {
			break
		}
		if t == nil {
			pw.center = image.Point{
				int(e.X),
				int(e.Y),
			}
			pw.w.Send(paint.Event{})
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src) // Змінено на білий фон.

	x, y := pw.center.X, pw.center.Y
	clr := color.RGBA{0, 0, 255, 255} // Змінено на синій колір.

	// Перший прямокутник малюється так, що він перетинається по центру горизонтально
	pw.w.Fill(image.Rect(x-150, y-50, x+150, y+50), clr, draw.Src)

	// Другий прямокутник малюється так, що він перетинається по центру вертикально
	pw.w.Fill(image.Rect(x-50, y-100, x+50, y+100), clr, draw.Src)

	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
