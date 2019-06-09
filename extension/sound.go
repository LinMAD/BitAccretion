package extension

import "github.com/LinMAD/BitAccretion/model"

// ISound player general interface
type ISound interface {
	// PlayAlert sound for given system name
	PlayAlert(name model.VertexName)
}
