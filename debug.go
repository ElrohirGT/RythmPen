package rythmpen

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var debugImage = func() *ebiten.Image {
	a := ebiten.NewImage(10, 10)
	a.Fill(color.RGBA{R: 255, G: 0, B: 255, A: 255})
	return a
}()

type DebugImage struct {
	Image    *ebiten.Image
	Position Vec2
}

func NewDebugImage(position Vec2) *DebugImage {
	return &DebugImage{
		Image:    debugImage,
		Position: position,
	}
}

func (db *DebugImage) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	db.Position.Vec2TranslateGeom(opt)
	parent.DrawImage(db.Image, opt)
	opt.GeoM.Reset()
}

type DebugImageManager struct {
	enabled       bool
	activationKey ebiten.Key
	images        []*DebugImage
}

func NewDebugImageManager(activationKey ebiten.Key) *DebugImageManager {
	return &DebugImageManager{
		images:        make([]*DebugImage, 0, 10),
		activationKey: activationKey,
	}
}

func (manager *DebugImageManager) Update() {
	if inpututil.IsKeyJustPressed(manager.activationKey) {
		manager.enabled = !manager.enabled
	}
}

func (manager *DebugImageManager) Draw(parent *ebiten.Image, opt *ebiten.DrawImageOptions) {
	if manager.enabled {
		for _, i := range manager.images {
			i.Draw(parent, opt)
		}
	}
}

func (manager *DebugImageManager) Add(img *DebugImage) {
	manager.images = append(manager.images, img)
}
