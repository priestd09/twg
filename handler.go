package main

import (
	"github.com/nsf/termbox-go"
	"github.com/waka/twg/twitter"
	"github.com/waka/twg/views"
)

// EventHandler
type Handler struct {
	args      []string
	apiClient *twitter.Client
	quit      bool
	container *views.Container
}

func NewHandler(args []string, apiClient *twitter.Client) *Handler {
	return &Handler{args: args, apiClient: apiClient, quit: false}
}

func (handler *Handler) MainLoop() error {
	if err := handler.setupContainer(); err != nil {
		return err
	}
	defer handler.finish()

	handler.reset()

	eventCh := make(chan termbox.Event)
	defer close(eventCh)

	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case event := <-eventCh:
			if event.Type == termbox.EventResize {
				handler.reset()
			} else {
				handler.handleEvent(event)
			}
		}
		if handler.quit {
			break
		}
	}

	return nil
}

func (handler *Handler) setupContainer() error {
	handler.container = views.NewContainer()
	if err := handler.container.Setup(); err != nil {
		return err
	}
	return nil
}

func (handler *Handler) reset() {
	handler.container.Render()
}

func (handler *Handler) finish() {
	//handler.apiClient.Close()
	handler.container.Dispose()
}

func (handler *Handler) handleEvent(event termbox.Event) {
	switch handler.getKeyEvent(event) {
	case ACTION_RELOAD:
		// refresh data and scroll top
	case ACTION_QUIT:
		// quit loop
		handler.quit = true
	case ACTION_UP:
		// select next tweet
	case ACTION_DOWN:
		// select prev tweet
	case ENTER_NORMAL_MODE:
		// disable command view
		handler.container.ChangeCommandMode(false)
	case ENTER_COMMAND_MODE:
		// enable command view
		handler.container.ChangeCommandMode(true)
		handler.container.SetRuneInCommand(':')
	default:
		if handler.container.IsCommandMode() && event.Ch != 0 {
			handler.container.SetRuneInCommand(event.Ch)
		}
	}
}

func (handler *Handler) getKeyEvent(event termbox.Event) Action {
	for _, keybind := range KeybindList {
		if event.Mod == keybind.Mod && event.Key == keybind.Key && event.Ch == keybind.Ch {
			return keybind.Action
		}
	}
	return NO_ACTION
}

func (handler *Handler) processCommand(event termbox.Event) {
}
