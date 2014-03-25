package main

import (
	"fmt"
	"gopkg.in/v0/qml"
	"image/color"
	"math/rand"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	qml.Init(nil)
	engine := qml.NewEngine()
	container := &ModelContainer{models: make([]*Colors, 0), CurrentModel: &Colors{}}
	container.models = append(container.models, &Colors{})
	container.CurrentModel.CopyFrom(container.models[0])
	engine.Context().SetVar("container", container)
	component, err := engine.LoadFile("delegate.qml")
	if err != nil {
		return err
	}
	window := component.CreateWindow(nil)
	window.Show()
	go func() {
		n := func() uint8 { return uint8(rand.Intn(256)) }
		for i := 0; i < 100; i++ {
			container.CurrentModel.Add(color.RGBA{n(), n(), n(), 0xff})
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		//replace with newly created model after 5 seconds
		time.Sleep(5 * time.Second)
		//create second model
		container.models = append(container.models, &Colors{})
		//backup first model
		container.models[0].CopyFrom(container.CurrentModel)
		//swap models
		container.CurrentModel.CopyFrom(container.models[1])

		//switch back to old model 3 seconds later
		time.Sleep(3 * time.Second)
		//backup second model
		container.models[1].CopyFrom(container.CurrentModel)
		//swap models
		container.CurrentModel.CopyFrom(container.models[0])
	}()

	window.Wait()
	return nil
}

type Colors struct {
	list []color.RGBA
	Len  int
}

func (colors *Colors) Add(c color.RGBA) {
	colors.list = append(colors.list, c)
	colors.Len = len(colors.list)
	qml.Changed(colors, &colors.Len)
}

func (colors *Colors) Color(index int) color.RGBA {
	return colors.list[index]
}

type ModelContainer struct {
	CurrentModel *Colors
	models       []*Colors
}

func (dst *Colors) CopyFrom(src *Colors) {
	dst.list = src.list
	dst.Len = src.Len
	qml.Changed(dst, &dst.Len)
}
