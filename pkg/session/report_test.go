package session_test

import (
	"errors"

	"github.com/beatkind/watchtower/pkg/session"
	wt "github.com/beatkind/watchtower/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// fakeContainer embeds the full container interface and only implements the
// methods the session package reads during a run
type fakeContainer struct {
	wt.Container

	id        wt.ContainerID
	name      string
	imageName string
	imageID   wt.ImageID
	noNotify  bool
}

func (f fakeContainer) ID() wt.ContainerID      { return f.id }
func (f fakeContainer) Name() string            { return f.name }
func (f fakeContainer) ImageName() string       { return f.imageName }
func (f fakeContainer) SafeImageID() wt.ImageID { return f.imageID }
func (f fakeContainer) NoNotify() bool          { return f.noNotify }

var _ = Describe("the session report", func() {
	newProgress := func() session.Progress {
		return session.Progress{}
	}

	scanned := func(m session.Progress, name string, noNotify bool) wt.ContainerID {
		c := fakeContainer{
			id:        wt.ContainerID("01d1" + name),
			name:      name,
			imageName: "mock/" + name + ":latest",
			imageID:   wt.ImageID("01d1" + name),
			noNotify:  noNotify,
		}
		m.AddScanned(c, wt.ImageID("d0a1"+name))
		m.MarkForUpdate(c.id)
		return c.id
	}

	Describe("AllUpdatedNoNotify", func() {
		It("is true when every updated container carries the no-notify label", func() {
			m := newProgress()
			scanned(m, "updt1", true)
			scanned(m, "updt2", true)

			Expect(m.Report().AllUpdatedNoNotify()).To(BeTrue())
		})

		It("is false when any updated container lacks the no-notify label", func() {
			m := newProgress()
			scanned(m, "updt1", true)
			scanned(m, "updt2", false)

			Expect(m.Report().AllUpdatedNoNotify()).To(BeFalse())
		})

		It("is false when an update failed", func() {
			m := newProgress()
			scanned(m, "updt1", true)
			m.UpdateFailed(map[wt.ContainerID]error{wt.ContainerID("01d1updt1"): errors.New("update failed")})

			Expect(m.Report().AllUpdatedNoNotify()).To(BeFalse())
		})

		It("is false when a container was skipped", func() {
			m := newProgress()
			scanned(m, "updt1", true)
			c := fakeContainer{
				id:        wt.ContainerID("01d1skip1"),
				name:      "skip1",
				imageName: "mock/skip1:latest",
				imageID:   wt.ImageID("01d1skip1"),
			}
			m.AddSkipped(c, errors.New("unpossible"))

			Expect(m.Report().AllUpdatedNoNotify()).To(BeFalse())
		})

		It("is false when only fresh containers were seen", func() {
			m := newProgress()
			c := fakeContainer{
				id:        wt.ContainerID("01d1frsh1"),
				name:      "frsh1",
				imageName: "mock/frsh1:latest",
				imageID:   wt.ImageID("d0d1frsh1"),
			}
			m.AddScanned(c, c.imageID)

			Expect(m.Report().AllUpdatedNoNotify()).To(BeFalse())
		})

		It("is false when no containers were updated", func() {
			Expect(newProgress().Report().AllUpdatedNoNotify()).To(BeFalse())
		})
	})
})
